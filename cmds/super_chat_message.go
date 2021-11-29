package cmds

type SuperChatMessage struct {
	ID                    int64   `json:"id"`
	BackgroundBottomColor string  `json:"background_bottom_color"`
	BackgroundColor       string  `json:"background_color"`
	BackgroundColorEnd    string  `json:"background_color_end"`
	BackgroundColorStart  string  `json:"background_color_start"`
	BackgroundIcon        string  `json:"background_icon"`
	BackgroundImage       string  `json:"background_image"`
	BackgroundPriceColor  string  `json:"background_price_color"`
	ColorPoint            float64 `json:"color_point"`
	DMScore               int64   `json:"dmscore"`
	Time                  int     `json:"time"` // 持续秒数
	StartTime             int64   `json:"start_time"`
	EndTime               int64   `json:"end_time"`
	Gift                  struct {
		GiftID   int64  `json:"gift_id"`
		GiftName string `json:"gift_name"`
		Num      int    `json:"num"`
	} `json:"gift"`
	IsRanked         int       `json:"is_ranked"`
	IsSendAudit      int       `json:"is_send_audit"`
	Medal            MedalInfo `json:"medal_info"`
	Message          string    `json:"message"`
	MessageFontColor string    `json:"message_font_color"`
	MessageTrans     string    `json:"message_trans"`
	Price            int       `json:"price"` // 单位：人民币
	Rate             int       `json:"rate"`
	Token            string    `json:"token"`
	TransMark        int       `json:"trans_mark"`
	Timestamp        int64     `json:"ts"`
	UID              int64     `json:"uid"`
	User             UserInfo  `json:"user_info"`
}

func (i *SuperChatMessage) Decode(data []byte) error {
	return jsonDecode(data, i)
}
