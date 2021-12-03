package BLiveDanmaku

import (
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
var ErrHeartbeatTimeout error = errors.New("heartbeat timeout")

type OpHandler func(*Client, *RawMessage) bool
type CmdHandler func(*Client, string, []byte) bool
type ErrHandler func(*Client, error)
type ServerDisconnectHandler func(*Client, error)

func Dial(bilibili_live_room_id int, conf *ClientConf) (*Client, error) {
	ret := &Client{}

	conf = conf.Clone()
	conf.AddOpHandler(OP_AUTH_REPLY, ret.onAuthReply).AddOpHandler(OP_SEND_MSG_REPLY, ret.onChatMsg)
	conf.AddOpHandler(OP_HEARTBEAT_REPLY, ret.onServerHeartbeat)
	conf.PrependCmdHandler(CMD_LIVE, ret.onLiveStateChange).PrependCmdHandler(CMD_PREPARING, ret.onLiveStateChange)

	if conf.HeartbeatInterval <= 0 {
		conf.HeartbeatInterval = time.Second * 30
	}

	if conf.HeartbeatTimeout < conf.HeartbeatInterval*2 {
		conf.HeartbeatTimeout = conf.HeartbeatInterval * 2
	}

	if conf.HandshakeTimeout <= 0 {
		conf.HandshakeTimeout = time.Second * 10
	}

	ret.conf = conf
	ret.last_heartbeat.Store(time.Now())
	err := ret.connect(bilibili_live_room_id)
	return ret, err
}

type ClientConf struct {
	OpHandlerMap  map[uint32][]OpHandler
	CmdHandlerMap map[string][]CmdHandler

	OnNetError         ErrHandler
	OnServerDisconnect ServerDisconnectHandler

	HeartbeatInterval time.Duration
	HeartbeatTimeout  time.Duration

	HandshakeTimeout time.Duration
}

func (c *ClientConf) Clone() *ClientConf {
	ret := &ClientConf{
		OpHandlerMap:  map[uint32][]OpHandler{},
		CmdHandlerMap: map[string][]CmdHandler{},
	}
	if c == nil {
		return ret
	}

	for key, handlers := range c.OpHandlerMap {
		ret.OpHandlerMap[key] = append([]OpHandler{}, handlers...)
	}

	for key, handlers := range c.CmdHandlerMap {
		ret.CmdHandlerMap[key] = append([]CmdHandler{}, handlers...)
	}

	ret.OnNetError = c.OnNetError
	ret.OnServerDisconnect = c.OnServerDisconnect

	ret.HeartbeatInterval = c.HeartbeatInterval
	ret.HeartbeatTimeout = c.HeartbeatTimeout
	ret.HandshakeTimeout = c.HandshakeTimeout
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

func (c *ClientConf) PrependCmdHandler(cmd string, cb CmdHandler) *ClientConf {
	if c.CmdHandlerMap == nil {
		c.CmdHandlerMap = make(map[string][]CmdHandler)
	}
	c.CmdHandlerMap[cmd] = append([]CmdHandler{cb}, c.CmdHandlerMap[cmd]...)
	return c
}

type Client struct {
	conn *websocket.Conn

	conf *ClientConf
	room atomic.Value

	closed         int32
	heartbeat_init int32
	last_heartbeat atomic.Value // time.Time
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
		dailer := websocket.Dialer{HandshakeTimeout: c.conf.HandshakeTimeout}
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
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	auth_data, _ := json.Marshal(auth)
	auth_msg := &RawMessage{
		Op:   OP_AUTH,
		Seq:  1,
		Data: auth_data,
	}
	if err := c.conn.WriteMessage(2, auth_msg.Encode()); err != nil {
		c.onError(err)
		return
	}

	// msg loop
	for {
		mt, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.onError(err)
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
	<-time.After(c.conf.HeartbeatInterval)
	if atomic.LoadInt32(&(c.closed)) != 0 {
		return
	}
	if time.Since(c.lastHeartbeat()) > c.conf.HeartbeatTimeout {
		c.Close()
		if c.conf.OnNetError != nil {
			c.conf.OnNetError(c, ErrHeartbeatTimeout)
		}
		return
	}
	c.SendMsg(&RawMessage{Op: OP_HEARTBEAT, Seq: 1, Data: HEARTBEAT_MSG})
	go c.heartBeat()
}

func (c *Client) lastHeartbeat() time.Time {
	if ret, ok := c.last_heartbeat.Load().(time.Time); ok {
		return ret
	}
	return time.Time{}
}

func (c *Client) onServerHeartbeat(_ *Client, msg *RawMessage) bool {
	c.last_heartbeat.Store(time.Now())
	return false
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
