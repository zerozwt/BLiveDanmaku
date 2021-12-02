package BLiveDanmaku

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var ErrNotConnected error = errors.New("not connected")
var ErrClosed error = errors.New("underlyinh websocket closed")

type OpHandler func(*Client, *RawMessage) bool
type CmdHandler func(*Client, string, []byte) bool
type ErrHandler func(*Client, error)
type ServerDisconnectHandler func(*Client, error)

func Dial(bilibili_live_room_id int, conf *ClientConf) (*Client, error) {
	ret := &Client{}

	conf = conf.Clone()
	conf.AddOpHandler(OP_AUTH_REPLY, ret.onAuthReply).AddOpHandler(OP_SEND_MSG_REPLY, ret.onChatMsg)
	conf.AddCmdHandler(CMD_LIVE, ret.onLiveStateChange).AddCmdHandler(CMD_PREPARING, ret.onLiveStateChange)

	ret.conf = conf
	err := ret.connect(bilibili_live_room_id)
	return ret, err
}

type ClientConf struct {
	OpHandlerMap  map[uint32][]OpHandler
	CmdHandlerMap map[string][]CmdHandler

	OnNetError         ErrHandler
	OnServerDisconnect ServerDisconnectHandler
}

func (c *ClientConf) Clone() *ClientConf {
	ret := &ClientConf{
		OpHandlerMap:  map[uint32][]OpHandler{},
		CmdHandlerMap: map[string][]CmdHandler{},
	}
	if c == nil {
		return ret
	}

	ret.OnNetError = c.OnNetError
	ret.OnServerDisconnect = c.OnServerDisconnect

	for key, handlers := range c.OpHandlerMap {
		ret.OpHandlerMap[key] = append([]OpHandler{}, handlers...)
	}

	for key, handlers := range c.CmdHandlerMap {
		ret.CmdHandlerMap[key] = append([]CmdHandler{}, handlers...)
	}

	return ret
}

func (c *ClientConf) AddOpHandler(op uint32, cb OpHandler) *ClientConf {
	if c.OpHandlerMap == nil {
		c.OpHandlerMap = make(map[uint32][]OpHandler)
	}
	c.OpHandlerMap[op] = append(c.OpHandlerMap[op], cb)
	return c
}

func (c *ClientConf) AddCmdHandler(cmd string, cb CmdHandler) *ClientConf {
	if c.CmdHandlerMap == nil {
		c.CmdHandlerMap = make(map[string][]CmdHandler)
	}
	c.CmdHandlerMap[cmd] = append(c.CmdHandlerMap[cmd], cb)
	return c
}

type Client struct {
	conn *websocket.Conn

	conf *ClientConf
	room atomic.Value

	closed         int32
	heartbeat_init int32
}

func (c *Client) Close() {
	if c.conn != nil && atomic.CompareAndSwapInt32(&(c.closed), 0, 1) {
		c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{})
		c.conn.Close()
	}
}

func (c *Client) Room() *RoomInfo {
	if room, ok := c.room.Load().(*RoomInfo); ok {
		return room
	}
	return &RoomInfo{}
}

func (c *Client) SendMsg(msg *RawMessage) error {
	if atomic.LoadInt32(&(c.closed)) != 0 {
		return ErrClosed
	}
	if c.conn != nil {
		return c.conn.WriteMessage(2, msg.Encode())
	}

	return ErrNotConnected
}

func (c *Client) connect(bilibili_live_room_id int) error {
	// get room info
	room_info, err := GetRoomInfo(bilibili_live_room_id)
	if err != nil {
		return err
	}
	c.room.Store(room_info)

	// get danmaku info
	danmaku_info, err := GetDanmakuInfo(room_info.Base.RoomID)
	if err != nil {
		return err
	}

	// try to connect hosts one by one
	for _, host := range danmaku_info.HostList {
		dailer := websocket.Dialer{}
		ws_url := "wss://" + host.Host + ":" + strconv.Itoa(host.WssPort) + "/sub"
		c.conn, _, err = dailer.Dial(ws_url, http.Header{})
		if err == nil {
			go c.msgLoop(danmaku_info.Token)
			return nil
		}
		logger().Printf("connect danmaku websocket server %s:%d failed: %v, try next server ...", host.Host, host.WssPort, err)
	}

	return err
}

func (c *Client) msgLoop(token string) {
	if atomic.LoadInt32(&(c.closed)) != 0 {
		return
	}

	// auth
	auth := map[string]interface{}{
		"uid":      0,
		"roomid":   c.Room().Base.RoomID,
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
		if atomic.LoadInt32(&(c.closed)) == 0 {
			c.onError(err)
		}
		return
	}

	// msg loop
	for {
		mt, msg, err := c.conn.ReadMessage()
		if err != nil {
			if atomic.LoadInt32(&(c.closed)) == 0 {
				c.onError(err)
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
			if cb(c, msg) {
				break
			}
		}
	}
}

func (c *Client) onError(err error) {
	if atomic.LoadInt32(&(c.closed)) != 0 {
		return
	}
	c.Close()
	if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		if c.conf.OnNetError != nil {
			c.conf.OnNetError(c, err)
		}
	} else {
		if c.conf.OnServerDisconnect != nil {
			c.conf.OnServerDisconnect(c, err)
		}
	}
}

func (c *Client) onAuthReply(_ *Client, msg *RawMessage) bool {
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

func (c *Client) onChatMsg(_ *Client, msg *RawMessage) bool {
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
			if cb(c, cmd, data) {
				break
			}
		}
	}
	return true
}

func (c *Client) onLiveStateChange(*Client, string, []byte) bool {
	// get room info
	room_info, err := GetRoomInfo(c.Room().Base.RoomID)
	if err == nil {
		c.room.Store(room_info)
	}
	return false
}
