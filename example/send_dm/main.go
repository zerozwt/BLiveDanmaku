package main

import (
	"flag"
	"fmt"

	dm "github.com/zerozwt/BLiveDanmaku"
)

func main() {
	var sender uint
	var reciever uint
	msg := ""
	sess_data := ""
	jct := ""

	flag.UintVar(&sender, "sender", 0, "sender UID")
	flag.UintVar(&reciever, "reciever", 0, "reciever UID")
	flag.StringVar(&msg, "msg", "", "direct msg to send")
	flag.StringVar(&sess_data, "sess_data", "", "your SESS_DATA")
	flag.StringVar(&jct, "jct", "", "your JCT")
	flag.Parse()

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
	if sender == 0 {
		fmt.Println("sender not specified")
		return
	}
	if reciever == 0 {
		fmt.Println("reciever not specifed")
		return
	}

	dev_id, err := dm.GetDMDeviceID()
	if err != nil {
		fmt.Printf("Get DM device id failed: %v", err)
		return
	}
	fmt.Println("dev_id:", dev_id)

	rsp, err := dm.SendDirectMsg(int64(sender), int64(reciever), msg, dev_id, sess_data, jct)

	if err != nil {
		fmt.Println("send msg failed:", err)
		return
	}

	fmt.Printf("code: %d message: %s\n", rsp.Code, rsp.Message)
}
