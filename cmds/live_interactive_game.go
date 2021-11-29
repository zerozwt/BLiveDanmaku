package cmds

type LiveInteractiveGame struct {
	Type           int         `json:"type"`
	UID            int64       `json:"uid"`
	UserName       string      `json:"uname"`
	UserFace       string      `json:"uface"`
	GiftID         int64       `json:"gift_id"`
	GiftName       string      `json:"gift_name"`
	GiftNum        int         `json:"gift_num"`
	Price          int64       `json:"price"` // 单位：金瓜子
	Paid           bool        `json:"paid"`
	Msg            string      `json:"msg"`
	FansMedalLevel int         `json:"fans_medal_level"`
	GuardLevel     int         `json:"guard_level"`
	Timestamp      int64       `json:"timestamp"`
	AnchorLottery  interface{} `json:"anchor_lottery"`
	PKInfo         interface{} `json:"pk_info"`
	AnchorInfo     struct {
		UID      int64  `json:"uid"`
		UserName string `json:"uname"`
		UserFace string `json:"uface"`
	} `json:"anchor_info"`
}

func (i *LiveInteractiveGame) Decode(data []byte) error {
	return jsonDecode(data, i)
}
