package BLiveDanmaku

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var g_logger atomic.Value

func init() {
	SetLogWriter(NullWriter{})
}

func GetRoomInfo(room_id int) (*RoomInfo, error) {
	room_info := &RoomInfo{}
	info_rsp := struct {
		Code    int       `json:"code"`
		Message string    `json:"message"`
		Data    *RoomInfo `json:"data"`
	}{
		Data: room_info,
	}
	err := httpGet(ROOM_INFO_API, map[string]string{"room_id": strconv.Itoa(room_id)}, &info_rsp)
	if err != nil {
		return nil, err
	}
	if info_rsp.Code != 0 {
		return nil, errors.New("Get room info failed: [" + strconv.Itoa(info_rsp.Code) + "] " + info_rsp.Message)
	}
	return room_info, nil
}

func GetDanmakuInfo(room_id int) (*DanmakuInfo, error) {
	dm_rsp := struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    DanmakuInfo `json:"data"`
	}{}
	err := httpGet(DANMAKU_INFO_API, map[string]string{"id": strconv.Itoa(room_id), "type": "0"}, &dm_rsp)
	if err != nil {
		return nil, err
	}
	if dm_rsp.Code != 0 {
		return nil, errors.New("Get danmaku info failed: [" + strconv.Itoa(dm_rsp.Code) + "] " + dm_rsp.Message)
	}

	if len(dm_rsp.Data.HostList) == 0 {
		dm_rsp.Data.HostList = []DanmakuHost{{
			Host:    `broadcastlv.chat.bilibili.com`,
			Port:    2243,
			WssPort: 443,
			WsPort:  2244,
		}}
	}
	return &dm_rsp.Data, nil
}

func SendMsg(msg string, room *RoomInfo, sess_data, jct string) error {
	body := url.Values{}
	body.Set("bubble", "0")
	body.Set("msg", msg)
	body.Set("color", strconv.Itoa(0xFFFFFF))
	body.Set("mode", "1")
	body.Set("fontsize", "25")
	body.Set("rnd", fmt.Sprint(time.Now().Unix()))
	body.Set("roomid", strconv.Itoa(room.Base.RoomID))
	body.Set("csrf", jct)
	body.Set("csrf_token", jct)
	req, _ := http.NewRequest("POST", SEND_MSG_API, strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", fmt.Sprintf("SESSDATA=%s; bili_jct=%s", sess_data, jct))

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return fmt.Errorf("http request failed: %d", rsp.StatusCode)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logger().Printf("http read body failed: url: %s err: %v", SEND_MSG_API, err)
		return err
	}

	tmp := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	if tmp.Code != 0 {
		return fmt.Errorf("send msg failed: [%d] %s", tmp.Code, tmp.Message)
	}
	return nil
}

func httpGet(base_url string, params map[string]string, rsp interface{}) error {
	tmp := strings.Builder{}
	tmp.WriteString(base_url)
	first := true
	for k, v := range params {
		if first {
			tmp.WriteByte('?')
			first = false
		} else {
			tmp.WriteByte('&')
		}
		tmp.WriteString(k)
		tmp.WriteByte('=')
		tmp.WriteString(url.QueryEscape(v))
	}
	real_url := tmp.String()

	http_rsp, err := http.Get(real_url)
	if err != nil {
		logger().Printf("http get %s failed: %v", real_url, err)
		return err
	}
	defer http_rsp.Body.Close()

	if http_rsp.StatusCode != 200 {
		return fmt.Errorf("http request failed: %d", http_rsp.StatusCode)
	}

	data, err := ioutil.ReadAll(http_rsp.Body)
	if err != nil {
		logger().Printf("http read body failed: url: %s err: %v", real_url, err)
		return err
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, rsp)
}

func logger() *log.Logger {
	return g_logger.Load().(*log.Logger)
}

func SetLogWriter(out io.Writer) {
	g_logger.Store(log.New(out, "BLiveDanmaku", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile))
}

type NullWriter struct{}

func (w NullWriter) Write(data []byte) (int, error) {
	return len(data), nil
}
