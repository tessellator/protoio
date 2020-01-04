package protoio

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"

	proto "github.com/golang/protobuf/proto"
)

func TestRead(t *testing.T) {
	msg := testProtoMessage{content: "some data"}
	marshaledBytes, err := proto.Marshal(&msg)
	length := len(marshaledBytes)

	if err != nil {
		t.Fatalf("Marshal() returned err: %+v", err)
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	err = binary.Write(w, binary.BigEndian, int32(length))
	if err != nil {
		t.Fatalf("received error writing length: %+v", err)
	}

	_, err = w.Write(marshaledBytes)
	if err != nil {
		t.Fatalf("received error writing data: %+v", err)
	}

	w.Flush()

	r := bufio.NewReader(&buf)
	readMsg := testProtoMessage{}

	err = Read(r, &readMsg)

	if err != nil {
		t.Fatalf("Read() returned err: %+v", err)
	}

	want := "some data"
	got := string(msg.content)

	if want != got {
		t.Errorf("event.GetData(): want %s; got %s", want, got)
	}
}

func TestWrite(t *testing.T) {
	msg := testProtoMessage{content: "some data"}
	marshaledBytes, err := proto.Marshal(&msg)
	length := len(marshaledBytes)

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	err = Write(w, &msg)
	if err != nil {
		t.Fatalf("Write() returned err: %+v", err)
	}
	w.Flush()

	var readLength int32
	r := bufio.NewReader(&buf)
	binary.Read(r, binary.BigEndian, &readLength)

	if int32(length) != readLength {
		t.Fatalf("Did not read correct length; got %d, want %d", readLength, length)
	}

	readBytes := make([]byte, readLength)
	_, err = r.Read(readBytes)

	if err != nil {
		t.Fatalf("Error reading bytes from reader: %+v", err)
	}

	readMsg := testProtoMessage{}
	err = proto.Unmarshal(readBytes, &readMsg)
	if err != nil {
		t.Errorf("Error unmarshaling msg: %+v", err)
	}

	want := "some data"
	got := readMsg.content

	if want != got {
		t.Errorf("event.GetData(): want %s; got %s", want, got)
	}
}

// -----------------------------------------------------------------------------
// A stubbed proto.Message object

type testProtoMessage struct {
	content string
}

func (msg *testProtoMessage) Reset() {
	// nothing to do
}

func (msg *testProtoMessage) String() string {
	return msg.content
}

func (msg *testProtoMessage) ProtoMessage() {
	// nothing to do
}

func (msg *testProtoMessage) Marshal() ([]byte, error) {
	return []byte(msg.content), nil
}

func (msg *testProtoMessage) Unmarshal(bytes []byte) error {
	msg.content = string(bytes)
	return nil
}
