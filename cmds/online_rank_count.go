package cmds

type OnlineRankCount struct {
	Count int `json:"count"`
}

func (i *OnlineRankCount) Decode(data []byte) error {
	return jsonDecode(data, i)
}
