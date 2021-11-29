package cmds

type UserToastMsg struct {
	AnchorShow       bool   `json:"anchor_show"`
	Color            string `json:"color"`
	DMScore          int64  `json:"dmscore"`
	StartTime        int64  `json:"start_time"`
	EndTime          int64  `json:"end_time"`
	GuardLevel       int    `json:"guard_level"`
	IsShow           int    `json:"is_show"`
	Num              int    `json:"num"`
	OpType           int    `json:"op_type"`
	PayFlowID        string `json:"payflow_id"`
	Price            int64  `json:"price"` // 单位：金瓜子
	RoleName         string `json:"role_name"`
	SvgaBlock        int    `json:"svga_block"`
	TargetGuardCount int    `json:"target_guard_count"`
	ToastMsg         string `json:"toast_msg"`
	UID              int64  `json:"uid"`
	Unit             string `json:"unit"`
	UserShow         bool   `json:"user_show"`
	UserName         string `json:"username"`
}

func (i *UserToastMsg) Decode(data []byte) error {
	return jsonDecode(data, i)
}
