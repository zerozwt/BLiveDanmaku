package cmds

type SuperChatMessageDelete struct {
	Ids []int64 `json:"ids"`
}

func (i *SuperChatMessageDelete) Decode(data []byte) error {
	return jsonDecode(data, i)
}
