package cmds

type AnchorLotEnd struct {
	ID int64 `json:"id"`
}

func (i *AnchorLotEnd) Decode(data []byte) error {
	return jsonDecode(data, i)
}
