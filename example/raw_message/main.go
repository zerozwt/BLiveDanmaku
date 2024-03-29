package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dm "github.com/zerozwt/BLiveDanmaku"
)

var room_id int64
var log string

func waitSignal(done chan bool) {
	ch_sig := make(chan os.Signal, 1)
	signal.Notify(ch_sig, syscall.SIGINT, syscall.SIGTERM)
	<-ch_sig
	fmt.Println("exit signal recieved")
	close(done)
}

func onMsg(_ *dm.Client, msg *dm.RawMessage) bool {
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

var client *dm.Client
var conf *dm.ClientConf

func onNetError(_ *dm.Client, err error) {
	fmt.Printf("Connection loss detected [%v], reconnect...\n", err)
	for err != nil {
		client, err = dm.Dial(room_id, conf)
		fmt.Printf("Reconnect failed: %v, try again ...", err)
	}
}

func main() {
	flag.Int64Var(&room_id, "room_id", 0, "room id of bilibili live room")
	flag.StringVar(&log, "log", "", "[optional] log file")
	flag.Parse()

	if room_id == 0 {
		fmt.Println("--room_id=BLIVE_ROOM_ID")
		return
	}

	conf = &dm.ClientConf{
		OpHandlerMap: map[uint32][]dm.OpHandler{
			dm.OP_SEND_MSG_REPLY: {onMsg},
		},
		OnNetError: onNetError,
	}
	var err error
	client, err = dm.Dial(room_id, conf)
	if err != nil {
		fmt.Println("Connect to live room failed:", err)
		return
	}
	defer client.Close()

	ch_done := make(chan bool)
	go waitSignal(ch_done)
	<-ch_done
}
