package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	dm "github.com/zerozwt/BLiveDanmaku"
)

func main() {
	sess_data := ""
	jct := ""
	file_name := ""

	flag.StringVar(&file_name, "file", "", "picture to upload")
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
		fmt.Println("file is empty")
		return
	}

	// open file
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

	fmt.Printf("uploaded file: %+v\n", info)
}
