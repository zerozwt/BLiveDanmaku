package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	dm "github.com/zerozwt/BLiveDanmaku"
)

func main() {
	var sender uint
	var reciever uint
	file_name := ""
	sess_data := ""
	jct := ""

	flag.UintVar(&sender, "sender", 0, "sender UID")
	flag.UintVar(&reciever, "reciever", 0, "reciever UID")
	flag.StringVar(&file_name, "file_name", "", "picture to send")
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
	if len(file_name) == 0 {
		fmt.Println("file_name is empty")
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

	// get device ID
	dev_id, err := dm.GetDMDeviceID()
	if err != nil {
		fmt.Printf("Get DM device id failed: %v", err)
		return
	}
	fmt.Println("dev_id:", dev_id)

	// upload picture
	data, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Printf("open picture file %s failed: %v\n", file_name, err)
		return
	}

	info, err := dm.UploadPic(data, filepath.Base(file_name), sess_data, jct)
	if err != nil {
		fmt.Printf("upload file to bilibili failed: %v\n", err)
		return
	}

	// get ext
	ext := filepath.Ext(file_name)
	if len(ext) > 0 {
		ext = ext[1:]
	}

	rsp, err := dm.SendDirectMsgPicture(int64(sender), int64(reciever), info, ext, dev_id, sess_data, jct)

	if err != nil {
		fmt.Println("send msg failed:", err)
		return
	}

	fmt.Printf("code: %d message: %s\n", rsp.Code, rsp.Message)
}
