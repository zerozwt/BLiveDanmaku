package cmds

import (
	jsoniter "github.com/json-iterator/go"
)

type AnchorLotAward struct {
	ID         int64                `json:"id"`
	AwardImage string               `json:"award_image"`
	AwardName  string               `json:"award_name"`
	AwardNum   int                  `json:"award_num"`
	AwardUsers []AnchorLotAwardUser `json:"award_users"`
	LotStatus  int                  `json:"lot_status"` // 2 might be lottery end
	Url        string               `json:"url"`
	WebUrl     string               `json:"web_url"`
}

type AnchorLotAwardUser struct {
	UID      int64  `json:"uid"`
	Name     string `json:"uname"`
	FaceIcon string `json:"face"`
	Level    int    `json:"level"`
	Color    int64  `json:"color"`
}

func (i *AnchorLotAward) Decode(data []byte) error {
	return jsonDecode(data, i)
}

func jsonDecode(data []byte, ptr interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, ptr)
}
