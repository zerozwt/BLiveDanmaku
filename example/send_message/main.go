package main

import (
	"flag"
	"fmt"

	dm "github.com/zerozwt/BLiveDanmaku"
)

func main() {
	msg := ""
	room_id := 0
	sess_data := ""
	jct := ""

	flag.StringVar(&msg, "msg", "", "danmaku msg to send")
	flag.IntVar(&room_id, "room_id", 0, "room id to send msg")
	flag.StringVar(&sess_data, "sess_data", "", "your SESS_DATA")
	flag.StringVar(&jct, "jct", "", "your JCT")
	flag.Parse()

	if room_id == 0 {
		fmt.Println("room_id not specified")
		return
	}
	if len(sess_data) == 0 {
		fmt.Println("sess_data is empty")
		return
	}
	if len(jct) == 0 {
		fmt.Println("jct is empty")
		return
	}
	if len(msg) == 0 {
		fmt.Println("msg is empty")
		return
	}

	room := &dm.RoomInfo{}
	room.Base.RoomID = room_id
	err := dm.SendMsg(msg, room, sess_data, jct)

	if err != nil {
		fmt.Println("send msg failed:", err)
		return
	}
}
