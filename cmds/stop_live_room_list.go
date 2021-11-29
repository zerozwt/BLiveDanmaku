package cmds

type StopLiveRoomList struct {
	RoomIDList []int `json:"room_id_list"`
}

func (i *StopLiveRoomList) Decode(data []byte) error {
	return jsonDecode(data, i)
}
