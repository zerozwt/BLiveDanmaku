package BLiveDanmaku

import (
	"bytes"
	"compress/flate"
	"errors"
	"io/ioutil"

	"github.com/andybalholm/brotli"
)

var ErrMsgIncomplete error = errors.New("msg not complete")
var ErrMsgUnknownProtocolVersion error = errors.New("msg has an unknown version")

type RawMessage struct {
	Ver  int
	Op   uint32
	Seq  uint32
	Data []byte
}

func (msg *RawMessage) Encode() []byte {
	total_len := uint32(len(msg.Data) + HEADER_LENGTH)
	ret := make([]byte, total_len)

	buf := writeUint32(ret, total_len)
	buf = writeUint16(buf, HEADER_LENGTH)
	buf = writeUint16(buf, VER_NORMAL)
	buf = writeUint32(buf, msg.Op)
	buf = writeUint32(buf, msg.Seq)

	copy(buf, msg.Data)

	return ret
}

func (msg *RawMessage) Decode(buf []byte) ([]byte, error) {
	if len(buf) < 4 {
		return buf, ErrMsgIncomplete
	}

	total_len, buf := readUint32(buf)
	if total_len > uint32(len(buf)+4) {
		return buf, ErrMsgIncomplete
	}

	header_len, buf := readUint16(buf)
	ver, buf := readUint16(buf)
	msg.Op, buf = readUint32(buf)
	msg.Seq, buf = readUint32(buf)

	buf = buf[header_len-HEADER_LENGTH:]
	ret := buf[total_len-uint32(header_len):]
	buf = buf[:total_len-uint32(header_len)]

	msg.Ver = int(ver)

	switch ver {
	case 0:
		msg.Data = make([]byte, len(buf))
		copy(msg.Data, buf)
	case VER_NORMAL:
		msg.Data = make([]byte, len(buf))
		copy(msg.Data, buf)
	case VER_DEFLATE:
		reader := flate.NewReader(bytes.NewReader(buf))
		defer reader.Close()
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return ret, err
		}
		msg.Data = data
	case VER_BROTLI:
		reader := brotli.NewReader(bytes.NewReader(buf))
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return ret, err
		}
		msg.Data = data
	default:
		return ret, ErrMsgUnknownProtocolVersion
	}

	return ret, nil
}

func writeUint32(buf []byte, value uint32) []byte {
	buf[0] = byte((value >> 24) & 0xFF)
	buf[1] = byte((value >> 16) & 0xFF)
	buf[2] = byte((value >> 8) & 0xFF)
	buf[3] = byte(value & 0xFF)
	return buf[4:]
}

func writeUint16(buf []byte, value uint16) []byte {
	buf[0] = byte((value >> 8) & 0xFF)
	buf[1] = byte(value & 0xFF)
	return buf[2:]
}

func readUint32(buf []byte) (uint32, []byte) {
	return (uint32(buf[0]) << 24) | (uint32(buf[1]) << 16) | (uint32(buf[2]) << 8) | uint32(buf[3]), buf[4:]
}

func readUint16(buf []byte) (uint16, []byte) {
	return (uint16(buf[0]) << 8) | uint16(buf[1]), buf[2:]
}
