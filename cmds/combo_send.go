package cmds

type ComboSend struct {
	Action           string      `json:"action"`
	BatchComboID     string      `json:"batch_combo_id"`
	BatchComboNum    int         `json:"batch_combo_num"`
	ComboID          string      `json:"combo_id"`
	ComboNum         int         `json:"combo_num"`
	ComboTotalCoin   int64       `json:"combo_total_coin"`
	DMScore          int64       `json:"dmscore"`
	GiftID           int64       `json:"gift_id"`
	GiftName         string      `json:"gift_name"`
	GiftNum          int         `json:"gift_num"`
	IsShow           int         `json:"is_show"`
	Medal            MedalInfo   `json:"medal_info"`
	NameColor        string      `json:"name_color"`
	RecieverUserName string      `json:"r_uname"`
	RecieverUID      int64       `json:"ruid"`
	SendMaster       interface{} `json:"send_master"`
	TotalNum         int         `json:"total_num"`
	UID              int64       `json:"uid"`
	UserName         string      `json:"uname"`
}

type MedalInfo struct {
	AnchorRoomID     int         `json:"anchor_roomid"`
	AnchorName       string      `json:"anchor_uname"`
	GuardLevel       int         `json:"guard_level"` // 1-总督 2-提督 3-舰长
	IconID           int64       `json:"icon_id"`
	IsLighted        int         `json:"is_lighted"`
	MedalColor       interface{} `json:"medal_color"` // maybe float64(integer) or string(in SUPER_CHAT_MESSAGE) <- Fxxk bulibuli
	MedalColorBorder int64       `json:"medal_color_border"`
	MedalColorEnd    int64       `json:"medal_color_end"`
	MedalColorStart  int64       `json:"medal_color_start"`
	MedalLevel       int         `json:"medal_level"`
	MedalName        string      `json:"medal_name"`
	Special          string      `json:"special"`
	TargetID         int64       `json:"target_id"` // 主播UID
	Score            int         `json:"score"`
}

func (i *ComboSend) Decode(data []byte) error {
	return jsonDecode(data, i)
}
