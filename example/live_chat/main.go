package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	dm "github.com/zerozwt/BLiveDanmaku"
	"github.com/zerozwt/BLiveDanmaku/cmds"
)

var room_id int

func waitSignal(done chan bool) {
	ch_sig := make(chan os.Signal, 1)
	signal.Notify(ch_sig, syscall.SIGINT, syscall.SIGTERM)
	<-ch_sig
	fmt.Println("exit signal recieved")
	close(done)
}

func writeGuard(level int, builder *strings.Builder) {
	switch level {
	case cmds.GUARD_LEVEL_KANCHOU:
		builder.WriteByte('[')
		builder.WriteString("舰长")
		builder.WriteByte(']')
	case cmds.GUARD_LEVEL_TEITOKU:
		builder.WriteByte('[')
		builder.WriteString("提督")
		builder.WriteByte(']')
	case cmds.GUARD_LEVEL_SOUTOKU:
		builder.WriteByte('[')
		builder.WriteString("总督")
		builder.WriteByte(']')
	}
}

func writeMedal(medal *cmds.MedalInfo, builder *strings.Builder) {
	if medal.MedalLevel <= 0 {
		return
	}
	builder.WriteByte('[')
	builder.WriteString(strconv.Itoa(medal.MedalLevel))
	builder.WriteByte('|')
	builder.WriteString(medal.MedalName)
	builder.WriteByte(']')
}

func onChat(client *dm.Client, cmd string, data []byte) bool {
	obj := cmds.DanmakuMsg{}
	if err := obj.Decode(data); err != nil {
		fmt.Printf("decode %s failed: %v\n", cmd, err)
		return false
	}

	builder := strings.Builder{}
	builder.WriteString("[CHAT]")

	writeGuard(obj.Sender.GuardLevel, &builder)
	writeMedal(&(obj.Sender.Medal), &builder)

	builder.WriteString(obj.Sender.UserName)
	builder.WriteString(": ")
	builder.WriteString(obj.Content)

	if obj.Style.IsEmotion {
		builder.WriteString("[EMOTION: ")
		builder.WriteString(obj.Style.Emotion.EmotionUnique)
		builder.WriteByte(']')
	}

	fmt.Println(builder.String())
	return false
}

func onSuperChat(client *dm.Client, cmd string, data []byte) bool {
	obj := cmds.SuperChatMessage{}
	if err := obj.Decode(data); err != nil {
		fmt.Printf("decode %s failed: %v\n", cmd, err)
		return false
	}

	builder := strings.Builder{}
	builder.WriteString("[SUPERCHAT]")

	builder.WriteByte('[')
	builder.WriteString("CNY ")
	builder.WriteString(strconv.Itoa(int(obj.Price)))
	builder.WriteByte(']')

	writeMedal(&(obj.Medal), &builder)

	builder.WriteString(obj.User.UserName)
	builder.WriteString(": ")
	builder.WriteString(obj.Message)

	if len(obj.MessageTrans) > 0 {
		builder.WriteString(" TRANS: ")
		builder.WriteString(obj.MessageTrans)
	}

	fmt.Println(builder.String())
	return false
}

func onNewGuard(client *dm.Client, cmd string, data []byte) bool {
	obj := cmds.GuardBuy{}
	if err := obj.Decode(data); err != nil {
		fmt.Printf("decode %s failed: %v\n", cmd, err)
		return false
	}

	builder := strings.Builder{}
	builder.WriteString("[NEW_GUARD]")

	s := fmt.Sprintf("UID:%d NAME:%s NUM:%d PRICE:%d GiftName:%s GuardLevel:%d",
		obj.UID, obj.UserName, obj.Num, obj.Price, obj.GiftName, obj.GuardLevel)
	builder.WriteString(s)

	fmt.Println(builder.String())
	return false
}

func onVIPEntry(client *dm.Client, cmd string, data []byte) bool {
	obj := cmds.EntryEffect{}
	if err := obj.Decode(data); err != nil {
		fmt.Printf("decode %s failed: %v\n", cmd, err)
		return false
	}

	builder := strings.Builder{}
	builder.WriteString("[VIP_ENTRY]")

	builder.WriteString(obj.CopyWritting)

	fmt.Println(builder.String())
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
	flag.IntVar(&room_id, "room_id", 0, "room id of bilibili live room")
	flag.Parse()

	if room_id == 0 {
		fmt.Println("--room_id=BLIVE_ROOM_ID")
		return
	}

	conf = &dm.ClientConf{
		OnNetError: onNetError,
	}
	conf.AddCmdHandler(dm.CMD_DANMU_MSG, onChat).AddCmdHandler(dm.CMD_SUPER_CHAT_MESSAGE, onSuperChat)
	conf.AddCmdHandler(dm.CMD_GUARD_BUY, onNewGuard)
	conf.AddCmdHandler(dm.CMD_ENTRY_EFFECT, onVIPEntry)
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
