package cmds

type AnchorLotStart struct {
	ID             int64  `json:"id"`
	AssetIcon      string `json:"asset_icon"`
	AwardImage     string `json:"award_image"`
	AwardName      string `json:"award_name"`
	AwardNum       int    `json:"award_num"`
	CurGiftNum     int    `json:"cur_gift_num"`
	CurrentTime    int64  `json:"current_time"` // unix timestamp
	Danmaku        string `json:"danmu"`
	GiftID         int64  `json:"gift_id"`
	GiftName       string `json:"gift_name"`
	GiftNum        int    `json:"gift_num"`
	GiftPrice      int64  `json:"gift_price"`  // 单位：金瓜子
	GoAwayTime     int    `json:"goaway_time"` // in seconds
	GoodsID        int64  `json:"goods_id"`
	IsBroadcast    int    `json:"is_broadcast"`
	JoinType       int    `json:"join_type"`
	LotStatus      int    `json:"lot_status"`
	MaxTime        int    `json:"max_time"` // in seconds
	RequireText    string `json:"require_text"`
	RequireType    int    `json:"require_type"` // 2 means anchor medal level reach $require_value
	RequireValue   int64  `json:"require_value"`
	RoomID         int    `json:"room_id"`
	SendGiftEnsure int    `json:"send_gift_ensure"`
	ShowPanel      int    `json:"show_panel"`
	Status         int    `json:"status"`
	RemainTime     int    `json:"time"` // in seconds
	Url            string `json:"url"`
	WebUrl         string `json:"web_url"`
}

func (i *AnchorLotStart) Decode(data []byte) error {
	return jsonDecode(data, i)
}
