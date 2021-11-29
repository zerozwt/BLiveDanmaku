package cmds

type HotRankChanged struct {
	Rank        int    `json:"rank"`
	Trend       int    `json:"trend"`
	CountDown   int    `json:"countdown"`
	Timestamp   int64  `json:"timestamp"`
	WebUrl      string `json:"web_url"`
	LiveUrl     string `json:"live_url"`
	BlinkUrl    string `json:"blink_url"`
	LiveLinkUrl string `json:"live_link_url"`
	PcLinkUrl   string `json:"pc_link_url"`
	Icon        string `json:"icon"`
	AreaName    string `json:"area_name"`
	RankDesc    string `json:"rank_desc"`
}

func (i *HotRankChanged) Decode(data []byte) error {
	return jsonDecode(data, i)
}
