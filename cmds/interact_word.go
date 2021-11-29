package cmds

type InteractWord struct {
	Contribution struct {
		Grade int `json:"grade"`
	} `json:"contribution"`
	DMScore       int64     `json:"dmscore"`
	Medal         MedalInfo `json:"fans_medal"`
	Indentities   []int     `json:"identities"`
	IsSpread      int       `json:"is_spread"`
	MsgType       int       `json:"msg_type"`
	RoomID        int       `json:"roomid"`
	Score         int64     `json:"score"`
	SpreadDesc    string    `json:"spread_desc"`
	SpreadInfo    string    `json:"spread_info"`
	TallIcon      int       `json:"tail_icon"`
	Timestamp     int64     `json:"timestamp"`
	TriggerTime   int64     `json:"trigger_time"` // 单位：纳秒
	UID           int64     `json:"uid"`
	UserName      string    `json:"uname"`
	UserNameColor string    `json:"uname_color"`
}

func (i *InteractWord) Decode(data []byte) error {
	return jsonDecode(data, i)
}
