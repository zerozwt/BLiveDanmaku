package BLiveDanmaku

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var ErrNotConnected error = errors.New("not connected")

type OpHandler func(*RawMessage) bool
type CmdHandler func(string, []byte) bool
type ErrHandler func(error)

func Dial(bilibili_live_room_id int, conf *ClientConf) (*Client, error) {
	ret := &Client{}

	if conf == nil {
		conf = &ClientConf{}
	}

	if conf.OpHandlerMap == nil {
		conf.OpHandlerMap = make(map[uint32][]OpHandler)
	}

	if conf.CmdHandlerMap == nil {
		conf.CmdHandlerMap = make(map[string][]CmdHandler)
	}

	conf.OpHandlerMap[OP_AUTH_REPLY] = append(conf.OpHandlerMap[OP_AUTH_REPLY], ret.onAuthReply)
	conf.OpHandlerMap[OP_SEND_MSG_REPLY] = append(conf.OpHandlerMap[OP_SEND_MSG_REPLY], ret.onChatMsg)

	ret.conf = conf
	err := ret.connect(bilibili_live_room_id)
	return ret, err
}

type ClientConf struct {
	OpHandlerMap  map[uint32][]OpHandler
	CmdHandlerMap map[string][]CmdHandler

	OnDisconnect ErrHandler
}

type Client struct {
	conn *websocket.Conn
	lock sync.Mutex

	conf *ClientConf
	room *RoomInfo

	host_shift     uint32
	closed         int32
	heartbeat_init int32
}

func (c *Client) Close() {
	if c.conn != nil && atomic.CompareAndSwapInt32(&(c.closed), 0, 1) {
		c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{})
	}
}

func (c *Client) Room() *RoomInfo {
	return c.room
}

func (c *Client) SendMsg(msg *RawMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn != nil {
		return c.conn.WriteMessage(2, msg.Encode())
	}

	return ErrNotConnected
}

func (c *Client) connect(bilibili_live_room_id int) error {
	// get room info
	room_info := &RoomInfo{}
	info_rsp := struct {
		Code    int       `json:"code"`
		Message string    `json:"message"`
		Data    *RoomInfo `json:"data"`
	}{
		Data: room_info,
	}
	err := c.httpGet(ROOM_INFO_API, map[string]string{"room_id": strconv.Itoa(bilibili_live_room_id)}, &info_rsp)
	if err != nil {
		return err
	}
	if info_rsp.Code != 0 {
		return errors.New("Get room info failed: [" + strconv.Itoa(info_rsp.Code) + "] " + info_rsp.Message)
	}
	c.room = room_info

	return c.reconnect()
}

func (c *Client) reconnect() error {
	// get danmaku info
	dm_rsp := struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    DanmakuInfo `json:"data"`
	}{}
	err := c.httpGet(DANMAKU_INFO_API, map[string]string{"id": strconv.Itoa(c.room.Base.RoomID), "type": "0"}, &dm_rsp)
	if err != nil {
		return err
	}
	if dm_rsp.Code != 0 {
		return errors.New("Get danmaku info failed: [" + strconv.Itoa(dm_rsp.Code) + "] " + dm_rsp.Message)
	}

	if len(dm_rsp.Data.HostList) == 0 {
		dm_rsp.Data.HostList = []DanmakuHost{{
			Host:    `broadcastlv.chat.bilibili.com`,
			Port:    2243,
			WssPort: 443,
			WsPort:  2244,
		}}
	}

	hosts := dm_rsp.Data.HostList
	shift := atomic.LoadUint32(&(c.host_shift))
	for i := uint32(0); i < shift%uint32(len(hosts)); i++ {
		hosts = append(hosts[1:], hosts[0])
	}

	// try to connect hosts one by one
	for _, host := range hosts {
		dailer := websocket.Dialer{}
		ws_url := "wss://" + host.Host + ":" + strconv.Itoa(host.WssPort) + "/sub"
		c.conn, _, err = dailer.Dial(ws_url, http.Header{})
		if err == nil {
			go c.msgLoop(dm_rsp.Data.Token)
			return nil
		}
	}

	return err
}

func (c *Client) httpGet(base_url string, params map[string]string, rsp interface{}) error {
	tmp := strings.Builder{}
	tmp.WriteString(base_url)
	first := true
	for k, v := range params {
		if first {
			tmp.WriteByte('?')
			first = false
		} else {
			tmp.WriteByte('&')
		}
		tmp.WriteString(k)
		tmp.WriteByte('=')
		tmp.WriteString(url.QueryEscape(v))
	}
	real_url := tmp.String()

	http_rsp, err := http.Get(real_url)
	if err != nil {
		return err
	}
	defer http_rsp.Body.Close()

	data, err := ioutil.ReadAll(http_rsp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, rsp)
}

func (c *Client) msgLoop(token string) {
	if atomic.LoadInt32(&(c.closed)) != 0 {
		c.conn.Close()
		c.conn = nil
		return
	}

	// auth
	auth := map[string]interface{}{
		"uid":      0,
		"roomid":   c.room.Base.RoomID,
		"protover": 3,
		"platform": "web",
		"type":     2,
		"key":      token,
	}
	auth_data, _ := json.Marshal(auth)
	auth_msg := &RawMessage{
		Op:   OP_AUTH,
		Seq:  1,
		Data: auth_data,
	}
	if err := c.conn.WriteMessage(2, auth_msg.Encode()); err != nil {
		atomic.AddUint32(&(c.host_shift), 1)
		if atomic.LoadInt32(&(c.closed)) == 0 {
			c.tryReconnect(err)
		}
		return
	}

	// msg loop
	for {
		mt, msg, err := c.conn.ReadMessage()
		if err != nil {
			if atomic.LoadInt32(&(c.closed)) == 0 {
				c.tryReconnect(err)
			}
			return
		}
		if mt == 1 || mt == 2 {
			c.decodeMessage(msg)
		}
	}
}

func (c *Client) decodeMessage(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	msg := &RawMessage{}
	data, err := msg.Decode(data)

	if err != nil {
		return data, err
	}

	if msg.Ver <= VER_NORMAL {
		c.dispatchMessage(msg)
		return data, nil
	}

	bundle := msg.Data
	for len(bundle) > 0 {
		bundle, err = c.decodeMessage(bundle)

		if err != nil {
			return data, err
		}
	}

	return data, nil
}

func (c *Client) dispatchMessage(msg *RawMessage) {
	if handlers, ok := c.conf.OpHandlerMap[msg.Op]; ok {
		for _, cb := range handlers {
			if cb(msg) {
				break
			}
		}
	}
}

func (c *Client) tryReconnect(err error) {
	if atomic.LoadInt32(&(c.closed)) != 0 {
		return
	}

	c.conn.Close()
	c.conn = nil
	if !websocket.IsCloseError(err, websocket.CloseNormalClosure) && c.conf.OnDisconnect != nil {
		c.conf.OnDisconnect(err)
	}
	go c.reconnect()
}

func (c *Client) onAuthReply(msg *RawMessage) bool {
	c.SendMsg(&RawMessage{Op: OP_HEARTBEAT, Seq: 1, Data: HEARTBEAT_MSG})
	if atomic.CompareAndSwapInt32(&(c.heartbeat_init), 0, 1) {
		go c.heartBeat()
	}
	return true
}

func (c *Client) heartBeat() {
	<-time.After(time.Second * 30)
	if atomic.LoadInt32(&(c.closed)) != 0 {
		return
	}
	c.SendMsg(&RawMessage{Op: OP_HEARTBEAT, Seq: 1, Data: HEARTBEAT_MSG})
	go c.heartBeat()
}

func (c *Client) onChatMsg(msg *RawMessage) bool {
	// get cmd
	iter := jsoniter.NewIterator(jsoniter.ConfigCompatibleWithStandardLibrary)
	iter = iter.ResetBytes(msg.Data)
	cmd := ""
	var data []byte = nil
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, key string) bool {
		if key == "cmd" {
			cmd = iter.ReadString()
			return true
		}
		json_data := iter.SkipAndReturnBytes()
		if key == "data" || key == "info" {
			data = json_data
		}
		return true
	})

	if handlers, ok := c.conf.CmdHandlerMap[cmd]; ok {
		for _, cb := range handlers {
			if cb(cmd, data) {
				break
			}
		}
	}
	return true
}
