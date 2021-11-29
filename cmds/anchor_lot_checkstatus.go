package cmds

type AnchorLotCheckStatus struct {
	ID           int64  `json:"id"`
	Status       int    `json:"status"`
	AnchorUID    int64  `json:"uid"`
	RejectReason string `json:"reject_reason"`
}

func (i *AnchorLotCheckStatus) Decode(data []byte) error {
	return jsonDecode(data, i)
}
