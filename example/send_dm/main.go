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
	dev_id := ""

	flag.UintVar(&sender, "sender", 0, "sender UID")
	flag.UintVar(&reciever, "reciever", 0, "reciever UID")
	flag.StringVar(&msg, "msg", "", "direct msg to send")
	flag.StringVar(&dev_id, "dev", "", "direct msg send dev_id")
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
	if len(dev_id) == 0 {
		fmt.Println("dev_id not specified")
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

	rsp, err := dm.SendDirectMsg(int64(sender), int64(reciever), msg, dev_id, sess_data, jct)

	if err != nil {
		fmt.Println("send msg failed:", err)
		return
	}

	fmt.Printf("code: %d message: %s\n", rsp.Code, rsp.Message)
}
