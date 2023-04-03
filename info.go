package BLiveDanmaku

type RoomInfo struct {
	Base struct {
		Uid           int64  `json:"uid"`
		RoomID        int    `json:"room_id"`
		Title         string `json:"title"`
		Cover         string `json:"cover"`
		LiveStatus    int    `json:"live_status"`
		LiveStartTime int64  `json:"live_start_time"`
	} `json:"room_info"`
	Liver struct {
		Base struct {
			Name   string `json:"uname"`
			Icon   string `json:"face"`
			Gender string `json:"gender"`
		} `json:"base_info"`
		Medal struct {
			Name    string `json:"medal_name"`
			Id      int    `json:"medal_id"`
			FanClub int    `json:"fansclub"`
		} `json:"medal_info"`
	} `json:"anchor_info"`
}

type DanmakuHost struct {
	Host    string `json:"host"` // Host only contains domain name
	Port    int    `json:"port"`
	WssPort int    `json:"wss_port"`
	WsPort  int    `json:"ws_port"`
}

type DanmakuInfo struct {
	Token    string        `json:"token"`
	HostList []DanmakuHost `json:"host_list"`
}

type UploadedPic struct {
	Height   int    `json:"image_height"`
	Width    int    `json:"image_width"`
	ImageURL string `json:"image_url"`
}
