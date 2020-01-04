// Package protoio provides functions for reading and writing protobuf messages
// in a manner that is appropriate for streaming data.
package protoio

import (
	"encoding/binary"
	"io"

	proto "github.com/golang/protobuf/proto"
)

// Read reads a protobuf message from r into msg.
//
// It expects that the data is formated as first an integer in big endian
// followed by the binary data to be unmarshaled as a protobuf message.
func Read(r io.Reader, msg proto.Message) error {
	var length int32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	readBytes := make([]byte, length)
	_, err = io.ReadFull(r, readBytes)
	if err != nil {
		return err
	}

	return proto.Unmarshal(readBytes, msg)
}

// Write writes the msg to w and returns the total number of bytes written.
//
// This function will first write the length of the message in big endian
// followed by the binary data representing the protobuf message.
func Write(w io.Writer, msg proto.Message) (int64, error) {
	out, err := proto.Marshal(msg)
	if err != nil {
		return 0, err
	}

	err = binary.Write(w, binary.BigEndian, int32(len(out)))
	if err != nil {
		return 0, err
	}

	byteCount, err := w.Write(out)

	return int64(byteCount + 4), err
}
