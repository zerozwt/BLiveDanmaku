package cmds

type SuperChatMessageJpn struct {
	ID                    string    `json:"id"`
	UID                   string    `json:"uid"`
	Price                 int       `json:"price"` // 单位：人民币
	Rate                  int       `json:"rate"`
	Message               string    `json:"message"`
	MessageJpn            string    `json:"message_jpn"`
	IsRanked              int       `json:"is_ranked"`
	BackgroundImage       string    `json:"background_image"`
	BackgroundColor       string    `json:"background_color"`
	BackgroundIcon        string    `json:"background_icon"`
	BackgroundPriceColor  string    `json:"background_price_color"`
	BackgroundBottomColor string    `json:"background_bottom_color"`
	Timestamp             int64     `json:"ts"`
	Token                 string    `json:"token"`
	Medal                 MedalInfo `json:"medal_info"`
	User                  UserInfo  `json:"user_info"`
	Time                  int       `json:"time"` // 持续秒数
	StartTime             int64     `json:"start_time"`
	EndTime               int64     `json:"end_time"`
	Gift                  struct {
		GiftID   int64  `json:"gift_id"`
		GiftName string `json:"gift_name"`
		Num      int    `json:"num"`
	} `json:"gift"`
}

type UserInfo struct {
	UserName   string `json:"uname"`
	Face       string `json:"face"`
	FaceFrame  string `json:"face_frame"`
	GuardLevel int    `json:"guard_level"`
	UserLevel  int    `json:"user_level"`
	LevelColor string `json:"level_color"`
	IsVIP      int    `json:"is_vip"`
	IsSVIP     int    `json:"is_svip"`
	IsMainVIP  int    `json:"is_main_vip"`
	Title      string `json:"title"`
	Manager    int    `json:"manager"`
	NameColor  string `json:"name_color"`
}

func (i *SuperChatMessageJpn) Decode(data []byte) error {
	return jsonDecode(data, i)
}
