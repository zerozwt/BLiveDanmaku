package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"syscall"

	dm "github.com/zerozwt/BLiveDanmaku"
)

var room_id int
var log string

func waitSignal(done chan bool) {
	ch_sig := make(chan os.Signal, 1)
	signal.Notify(ch_sig, syscall.SIGINT, syscall.SIGTERM)
	<-ch_sig
	fmt.Println("exit signal recieved")
	close(done)
}

func onMsg(msg *dm.Message) bool {
	if len(log) == 0 {
		fmt.Println(string(msg.Data))
		return false
	}

	f, err := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer f.Close()
	f.Write(msg.Data)
	f.Write([]byte("\n"))

	return false
}

func main() {
	flag.IntVar(&room_id, "room_id", 0, "room id of bilibili live room")
	flag.StringVar(&log, "log", "", "[optional] log file")
	flag.Parse()

	if room_id == 0 {
		fmt.Println("--room_id=BLIVE_ROOM_ID")
		return
	}

	conf := dm.ClientConf{
		OpHandlerMap: map[uint32][]dm.OpHandler{
			dm.OP_SEND_MSG_REPLY: []dm.OpHandler{onMsg},
		},
	}
	client, err := dm.Dial(room_id, &conf)
	if err != nil {
		fmt.Println("Connect to live room failed:", err)
		return
	}
	defer client.Close()

	ch_done := make(chan bool)
	go waitSignal(ch_done)
	<-ch_done
}