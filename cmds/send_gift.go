package cmds

const (
	COIN_TYPE_SILVER = "silver"
	COIN_TYPE_GOLD   = "gold"
)

type SendGift struct {
	Action         string    `json:"action"`
	BatchComboID   string    `json:"batch_combo_id"`
	BatchComboSend *struct { // maybe nil
		Action        string      `json:"action"`
		BatchComboID  string      `json:"batch_combo_id"`
		BatchComboNum int         `json:"batch_combo_num"`
		BlindGift     interface{} `json:"blind_gift"`
		GiftID        int64       `json:"gift_id"`
		GiftName      string      `json:"gift_name"`
		GiftNum       int         `json:"gift_num"`
		SendMaster    interface{} `json:"send_master"`
		UID           int64       `json:"uid"`
		UserName      string      `json:"uname"`
	} `json:"batch_combo_send"`
	BeatID          string      `json:"beatId"`
	BizSource       string      `json:"biz_source"`
	BlindGift       interface{} `json:"blind_gift"`
	BroadcastID     int64       `json:"broadcast_id"`
	CoinType        string      `json:"coin_type"` // silver or gold
	ComboResourceID int64       `json:"combo_resources_id"`
	ComboSend       *struct {   // maybe nil
		Action     string      `json:"action"`
		ComboID    string      `json:"combo_id"`
		ComboNum   int         `json:"combo_num"`
		GiftID     int64       `json:"gift_id"`
		GiftName   string      `json:"gift_name"`
		GiftNum    int         `json:"gift_num"`
		SendMaster interface{} `json:"send_master"`
		UID        int64       `json:"uid"`
		UserName   string      `json:"uname"`
	} `json:"combo_send"`
	ComboStayTime     int         `json:"combo_stay_time"`
	ComboTotalCoin    int64       `json:"combo_total_coin"`
	CritProb          int         `json:"crit_prob"`
	Demarcation       int         `json:"demarcation"`
	DiscountPrice     int64       `json:"discount_price"`
	DMScore           int64       `json:"dmscore"`
	Draw              int         `json:"draw"`
	Effect            int         `json:"effect"`
	EffectBlock       int         `json:"effect_block"`
	Face              string      `json:"face"`
	FloatSCResourceID int64       `json:"float_sc_resource_id"`
	GiftID            int64       `json:"giftId"`
	GiftName          string      `json:"giftName"`
	GiftType          int         `json:"giftType"`
	Gold              int64       `json:"gold"`
	GuardLevel        int         `json:"guard_level"`
	IsFirst           bool        `json:"is_first"`
	IsSpecialbatch    int         `json:"is_special_batch"`
	Magnification     int         `json:"magnification"`
	Medal             MedalInfo   `json:"medal_info"`
	NameColor         string      `json:"name_color"`
	Num               int         `json:"num"`
	OriginalGiftName  string      `json:"original_gift_name"`
	Price             int64       `json:"price"`
	RCost             int64       `json:"rcost"`
	Remain            int         `json:"remain"`
	Rnd               string      `json:"rnd"`
	SendMaster        interface{} `json:"send_master"`
	Silver            int         `json:"silver"`
	Super             int         `json:"super"`
	SuperBatchGiftNum int         `json:"super_batch_gift_num"`
	SuperGiftNum      int         `json:"super_gift_num"`
	SvgaBlock         int         `json:"svga_block"`
	TagImage          string      `json:"tag_image"`
	Tid               string      `json:"tid"`
	Timestamp         int64       `json:"timestamp"`
	TopList           interface{} `json:"top_list"`
	TotalCoin         int64       `json:"total_coin"`
	UID               int64       `json:"uid"`
	UserName          string      `json:"uname"`
}

func (i *SendGift) Decode(data []byte) error {
	return jsonDecode(data, i)
}
