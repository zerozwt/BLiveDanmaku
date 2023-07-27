package BLiveDanmaku

const (
	OP_HANDSHAKE       = 0
	OP_HANDSHAKE_REPLY = 1

	OP_HEARTBEAT       = 2
	OP_HEARTBEAT_REPLY = 3

	OP_SEND_MSG       = 4
	OP_SEND_MSG_REPLY = 5

	OP_AUTH       = 7
	OP_AUTH_REPLY = 8

	VER_NORMAL  = 1
	VER_DEFLATE = 2
	VER_BROTLI  = 3

	HEADER_LENGTH = 16

	ROOM_INFO_API    string = `https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom`
	DANMAKU_INFO_API string = `https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo`
	SEND_MSG_API     string = `https://api.live.bilibili.com/msg/send`
	SEND_DM_API      string = `https://api.vc.bilibili.com/web_im/v1/web_im/send_msg`
	UPLOAD_PIC_API   string = `https://api.bilibili.com/x/dynamic/feed/draw/upload_bfs`

	USER_AGENT string = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36`
)

var HEARTBEAT_MSG []byte = []byte(`[object Object]`)
