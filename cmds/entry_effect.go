package cmds

type EntryEffect struct {
	ID               int64   `json:"id"`
	UID              int64   `json:"uid"`
	TargetID         int64   `json:"target_id"`
	MockEffect       int     `json:"mock_effect"`
	Face             string  `json:"face"`
	PrivilegeType    int     `json:"privilege_type"`
	CopyWritting     string  `json:"copy_writing"`
	CopyColor        string  `json:"copy_color"`
	HighlightColor   string  `json:"highlight_color"`
	Priority         int     `json:"priority"`
	BasemapUrl       string  `json:"basemap_url"`
	ShowAvatar       int     `json:"show_avatar"`
	EffectiveTime    int     `json:"effective_time"`
	WebBasemapUrl    string  `json:"web_basemap_url"`
	WebEffectiveTime int     `json:"web_effective_time"`
	WebEffectClose   int     `json:"web_effect_close"`
	WebCloseTime     int     `json:"web_close_time"`
	Business         int     `json:"business"`
	CopyWrittingV2   string  `json:"copy_writing_v2"`
	IconList         []int64 `json:"icon_list"`
	MaxDelayTime     int     `json:"max_delay_time"`
	TriggerTime      int64   `json:"trigger_time"` // 单位：纳秒
	Identities       int     `json:"identities"`
}

func (i *EntryEffect) Decode(data []byte) error {
	return jsonDecode(data, i)
}
