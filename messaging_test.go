package nativemessaging

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strings"
	"testing"
)

// Test write
func write(t *testing.T, endian binary.ByteOrder) {
	var host MessagingHost
	buf := new(bytes.Buffer)
	value := "native message host"
	if endian == nil {
		host = NativeHost(nil, buf)
	} else {
		host = New(nil, buf, endian)
	}
	i, err := host.Write(strings.NewReader(value))

	if err != nil {
		t.Fatal(err)
	}

	if i != len(value)+binary.Size(uint32(0)) {
		t.Fatal("Invalid write length")
	}

	result := buf.String()[4:]
	if result != value {
		t.Fatalf("Expected: %s Got: %s", value, result)
	}
}

func TestWriteNativeEndian(t *testing.T) {
	write(t, nil)
}

func TestWriteLittleEndian(t *testing.T) {
	write(t, binary.LittleEndian)
}

func TestWriteBigEndian(t *testing.T) {
	write(t, binary.BigEndian)
}

// Test send

func send(t *testing.T, endian binary.ByteOrder) {
	var host MessagingHost
	value := struct{ Text string }{Text: "native messaging host"}
	buf := new(bytes.Buffer)

	if endian == nil {
		host = NativeHost(nil, buf)
	} else {
		host = New(nil, buf, endian)
	}

	i, err := host.Send(value)
	if err != nil {
		t.Fatal(err)
	}
	if i != buf.Len() {
		t.Fatalf("Invalid write length: %d", 1)
	}

	var result struct{ Text string }

	err = json.Unmarshal(buf.Bytes()[binary.Size(uint32(0)):], &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Text != value.Text {
		t.Fatalf("Invalid result: %s", buf)
	}
}
func TestSendNativeEndian(t *testing.T) {
	send(t, nil)
}

func TestSendLittleEndian(t *testing.T) {
	send(t, binary.LittleEndian)
}

func TestSendBigEndian(t *testing.T) {
	send(t, binary.BigEndian)
}

// Test Read

func read(t *testing.T, endian binary.ByteOrder) {
	var host MessagingHost
	value := struct{ Text string }{Text: "native messaging host"}
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	header := make([]byte, binary.Size(uint32(0)))

	if endian == nil {
		endian = NativeEndian
		endian.PutUint32(header, uint32(len(data)))
		host = NativeHost(bytes.NewReader(append(header, data...)), nil)
	} else {
		endian.PutUint32(header, uint32(len(data)))
		host = New(bytes.NewReader(append(header, data...)), nil, endian)
	}

	result, err := host.Read()

	if err != nil {
		t.Fatal(err)
	}

	if string(result) != string(data) {
		t.Fatalf("Got: %s: Expected: %s", string(data), string(result))
	}
}

func TestReadNativeEndian(t *testing.T) {
	read(t, nil)
}

func TestReadLittleEndian(t *testing.T) {
	read(t, binary.LittleEndian)
}

func TestReadBigEndian(t *testing.T) {
	read(t, binary.BigEndian)
}

// Test Receive

func receive(t *testing.T, endian binary.ByteOrder) {
	var host MessagingHost
	value := struct{ Text string }{Text: "native messaging host"}
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	header := make([]byte, binary.Size(uint32(0)))
	if endian == nil {
		endian = NativeEndian
		endian.PutUint32(header, uint32(len(b)))
		host = NativeHost(bytes.NewReader(append(header, b...)), nil)
	} else {
		endian.PutUint32(header, uint32(len(b)))
		host = New(bytes.NewReader(append(header, b...)), nil, endian)
	}
	var result struct{ Text string }
	err = host.Receive(&result)

	if err != nil {
		t.Fatal(err)
	}

	if result.Text != value.Text {
		t.Fatalf("Invalid result: %#v", result)
	}
}

func TestReceiveNativeEndian(t *testing.T) {
	receive(t, nil)
}

func TestReceiveLittleEndian(t *testing.T) {
	receive(t, binary.LittleEndian)
}

func TestReceiveBigEndian(t *testing.T) {
	receive(t, binary.BigEndian)
}
