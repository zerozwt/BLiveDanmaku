package cmds

type RoomRealtimeMessageUpdate struct {
	RoomID    int `json:"roomid"`
	Fans      int `json:"fans"`
	RedNotice int `json:"red_notice"`
	FansClub  int `json:"fans_club"`
}

func (i *RoomRealtimeMessageUpdate) Decode(data []byte) error {
	return jsonDecode(data, i)
}
