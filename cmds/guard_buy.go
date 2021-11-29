package cmds

const (
	GUARD_LEVEL_SOUTOKU = 1
	GUARD_LEVEL_TEITOKU = 2
	GUARD_LEVEL_KANCHOU = 3
)

type GuardBuy struct {
	UID        int64  `json:"uid"`
	UserName   string `json:"username"`
	GuardLevel int    `json:"guard_level"`
	Num        int    `json:"num"`
	Price      int64  `json:"price"` // 单位：金瓜子
	GiftID     int64  `json:"gift_id"`
	GiftName   string `json:"gift_name"`  // 舰长、提督、总督
	StartTime  int64  `json:"start_time"` // unix timestamp
	EndTime    int64  `json:"end_time"`
}

func (i *GuardBuy) Decode(data []byte) error {
	return jsonDecode(data, i)
}
