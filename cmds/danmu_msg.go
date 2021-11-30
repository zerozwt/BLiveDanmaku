package cmds

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
)

var ErrDanmakuRawDataIsNotArray error = errors.New("danmaku raw data is not array")

type DanmakuMsg struct {
	Style   DanmakuMsgStyle
	Content string
	Sender  DanmakuMsgSender
}

type DanmakuMsgStyle struct {
	DisplayMode int // 弹幕显示模式（滚动、顶部、底部）
	FontSize    int
	Color       int64
	Timestamp   int64
	Rnd         int64
	CRC32       string
	MsgType     int // 是否礼物弹幕（节奏风暴）
	Bubble      int
	IsEmotion   bool // 是否是表情图片弹幕
	Emotion     struct {
		BulgeDisplay  int    `json:"bulge_display"`
		EmotionUnique string `json:"emoticon_unique"` // 表情ID
		Height        int    `json:"height"`          // 单位：像素
		InPlayerArea  int    `json:"in_player_area"`
		IsDynamic     int    `json:"is_dynamic"`
		Url           string `json:"url"`   // 表情图片URL
		Width         int    `json:"width"` // 单位：像素
	}
	Extra struct {
		Mode           int    `json:"mode"`
		ShowPlayerType int    `json:"show_player_type"`
		Extra          string `json:"extra"` // json string
	}
}

type DanmakuMsgSender struct {
	UID                 int64
	UserName            string
	IsAdmin             int // 是否是房管
	IsVIP               int
	IsSVIP              int
	UserRank            int
	MobilePhoneVerified int
	UserNameColor       string
	Medal               MedalInfo
	UserLevel           int
	UserColor           interface{}
	UserLevelRank       interface{}
	OldTitle            string
	Title               string
	GuardLevel          int
}

func (i *DanmakuMsg) Decode(data []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, data)
	if iter.WhatIsNext() != jsoniter.ArrayValue {
		return ErrDanmakuRawDataIsNotArray
	}
	idx := 0
	var err error
	iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
		defer func() { idx += 1 }()
		switch idx {
		case 0:
			err = i.decodeDanmakuStyle(iter.SkipAndReturnBytes())
		case 1:
			i.Content = iter.ReadString()
		case 2:
			err = i.decodeSenderBasic(iter.SkipAndReturnBytes())
		case 3:
			err = i.decodeSenderMedal(iter.SkipAndReturnBytes())
		case 4:
			err = i.decodeSenderLevel(iter.SkipAndReturnBytes())
		case 5:
			err = i.decodeSenderTitle(iter.SkipAndReturnBytes())
		case 7:
			i.Sender.GuardLevel = iter.ReadInt()
		default:
			iter.Skip()
		}
		return err == nil
	})
	if iter.Error != nil {
		return iter.Error
	}
	return err
}

func (i *DanmakuMsg) decodeDanmakuStyle(data []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, data)
	if iter.WhatIsNext() != jsoniter.ArrayValue {
		return ErrDanmakuRawDataIsNotArray
	}
	idx := 0
	var err error
	iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
		defer func() { idx += 1 }()
		switch idx {
		case 1:
			i.Style.DisplayMode = iter.ReadInt()
		case 2:
			i.Style.FontSize = iter.ReadInt()
		case 3:
			i.Style.Color = iter.ReadInt64()
		case 4:
			i.Style.Timestamp = iter.ReadInt64()
		case 5:
			i.Style.Rnd = iter.ReadInt64()
		case 7:
			i.Style.CRC32 = iter.ReadString()
		case 9:
			i.Style.MsgType = iter.ReadInt()
		case 10:
			i.Style.Bubble = iter.ReadInt()
		case 13:
			i.Style.IsEmotion = iter.WhatIsNext() == jsoniter.ObjectValue
			if i.Style.IsEmotion {
				err = jsonDecode(iter.SkipAndReturnBytes(), &(i.Style.Emotion))
			} else {
				iter.Skip()
			}
		case 15:
			err = jsonDecode(iter.SkipAndReturnBytes(), &(i.Style.Extra))
		default:
			iter.Skip()
		}
		return err == nil
	})
	if iter.Error != nil {
		return iter.Error
	}
	return err
}

func (i *DanmakuMsg) decodeSenderBasic(data []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, data)
	if iter.WhatIsNext() != jsoniter.ArrayValue {
		return ErrDanmakuRawDataIsNotArray
	}
	idx := 0
	var err error
	iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
		defer func() { idx += 1 }()
		switch idx {
		case 0:
			i.Sender.UID = iter.ReadInt64()
		case 1:
			i.Sender.UserName = iter.ReadString()
		case 2:
			i.Sender.IsAdmin = iter.ReadInt()
		case 3:
			i.Sender.IsVIP = iter.ReadInt()
		case 4:
			i.Sender.IsSVIP = iter.ReadInt()
		case 5:
			i.Sender.UserRank = iter.ReadInt()
		case 6:
			i.Sender.MobilePhoneVerified = iter.ReadInt()
		case 7:
			i.Sender.UserNameColor = iter.ReadString()
		default:
			iter.Skip()
		}
		return err == nil
	})
	if iter.Error != nil {
		return iter.Error
	}
	return err
}

func (i *DanmakuMsg) decodeSenderMedal(data []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, data)
	if iter.WhatIsNext() != jsoniter.ArrayValue {
		return ErrDanmakuRawDataIsNotArray
	}
	idx := 0
	var err error
	iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
		defer func() { idx += 1 }()
		switch idx {
		case 0:
			i.Sender.Medal.MedalLevel = iter.ReadInt()
		case 1:
			i.Sender.Medal.MedalName = iter.ReadString()
		case 2:
			i.Sender.Medal.AnchorName = iter.ReadString()
		case 3:
			i.Sender.Medal.AnchorRoomID = iter.ReadInt()
		case 4:
			i.Sender.Medal.MedalColor = iter.ReadInt64()
		case 5:
			i.Sender.Medal.Special = iter.ReadString()
		case 6:
			i.Sender.Medal.IconID = iter.ReadInt64()
		case 7:
			i.Sender.Medal.MedalColorBorder = iter.ReadInt64()
		case 8:
			i.Sender.Medal.MedalColorStart = iter.ReadInt64()
		case 9:
			i.Sender.Medal.MedalColorEnd = iter.ReadInt64()
		case 10:
			i.Sender.Medal.GuardLevel = iter.ReadInt()
		case 11:
			i.Sender.Medal.IsLighted = iter.ReadInt()
		case 12:
			i.Sender.Medal.TargetID = iter.ReadInt64()
		default:
			iter.Skip()
		}
		return err == nil
	})
	if iter.Error != nil {
		return iter.Error
	}
	return err
}

func (i *DanmakuMsg) decodeSenderLevel(data []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, data)
	if iter.WhatIsNext() != jsoniter.ArrayValue {
		return ErrDanmakuRawDataIsNotArray
	}
	idx := 0
	var err error
	iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
		defer func() { idx += 1 }()
		switch idx {
		case 0:
			i.Sender.UserLevel = iter.ReadInt()
		case 2:
			i.Sender.UserColor = iter.Read()
		case 3:
			i.Sender.UserLevelRank = iter.Read()
		default:
			iter.Skip()
		}
		return err == nil
	})
	if iter.Error != nil {
		return iter.Error
	}
	return err
}

func (i *DanmakuMsg) decodeSenderTitle(data []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigCompatibleWithStandardLibrary, data)
	if iter.WhatIsNext() != jsoniter.ArrayValue {
		return ErrDanmakuRawDataIsNotArray
	}
	idx := 0
	var err error
	iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
		defer func() { idx += 1 }()
		switch idx {
		case 0:
			i.Sender.OldTitle = iter.ReadString()
		case 1:
			i.Sender.Title = iter.ReadString()
		default:
			iter.Skip()
		}
		return err == nil
	})
	if iter.Error != nil {
		return iter.Error
	}
	return err
}
