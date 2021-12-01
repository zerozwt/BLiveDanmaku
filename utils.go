package BLiveDanmaku

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

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
		return err
	}
	defer http_rsp.Body.Close()

	data, err := ioutil.ReadAll(http_rsp.Body)
	if err != nil {
		return err
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, rsp)
}
