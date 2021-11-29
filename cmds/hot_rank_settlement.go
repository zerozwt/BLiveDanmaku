package cmds

type HotRankSettlement struct {
	AreaName  string `json:"area_name"`
	CacheKey  string `json:"cache_key"`
	DMMsg     string `json:"dm_msg"`
	DMScore   int64  `json:"dmscore"`
	Face      string `json:"face"`
	Icon      string `json:"icon"`
	Rank      int    `json:"rank"`
	Timestamp int64  `json:"timestamp"`
	UserName  string `json:"uname"`
	Url       string `json:"url"`
}

func (i *HotRankSettlement) Decode(data []byte) error {
	return jsonDecode(data, i)
}
