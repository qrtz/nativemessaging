package nativemessaging

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

var (

	// ErrInvalidMessageSize unable to read message size
	ErrInvalidMessageSize = errors.New("Invalid message size")
	// ErrByteOrderNotSet byter order is set on the first read
	ErrByteOrderNotSet = errors.New("Byte order not set")
	messageSizeInBytes = binary.Size(uint32(0))
)

// MessagingHost interface represents the native messaging communication
type MessagingHost interface {
	Read() ([]byte, error)
	Write(io.Reader) (int, error)
	Send(interface{}) (int, error)
	Receive(v interface{}) error
}

type host struct {
	r  io.Reader
	w  io.Writer
	bo binary.ByteOrder
}

func (h *host) Read() ([]byte, error) {
	return Read(h.r, h.bo)
}

func (h *host) Write(message io.Reader) (int, error) {
	return Write(h.w, message, h.bo)
}

func (h *host) Send(v interface{}) (int, error) {
	return Send(h.w, v, h.bo)

}
func (h *host) Receive(v interface{}) error {
	return Receive(h.r, v, h.bo)
}

// NativeHost creates and returns an implementation of MessagingHost with native byte order
func NativeHost(stdin io.Reader, stdout io.Writer) MessagingHost {
	return &host{r: stdin, w: stdout, bo: NativeEndian}
}

// New creates and returns an implementation of MessagingHost
// This is convenient so you don't have to always pass around a reader and writer
func New(stdin io.Reader, stdout io.Writer, order binary.ByteOrder) MessagingHost {
	return &host{r: stdin, w: stdout, bo: order}
}

// Read reads a message from the reader
func Read(r io.Reader, order binary.ByteOrder) ([]byte, error) {
	b := make([]byte, messageSizeInBytes)
	i, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if i == 0 {
		return nil, ErrInvalidMessageSize
	}

	ln := order.Uint32(b)

	if ln == 0 {
		return nil, ErrInvalidMessageSize
	}
	m := make([]byte, ln)
	_, err = r.Read(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Receive parses the incoming JSON-encoded data and stores the result in the value pointed to by v.
func Receive(r io.Reader, v interface{}, order binary.ByteOrder) error {
	b, err := Read(r, order)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return errors.New(err.Error() + string(b))
	}
	return nil
}

// Write write a message to the writer
func Write(w io.Writer, message io.Reader, order binary.ByteOrder) (i int, err error) {
	data, err := ioutil.ReadAll(message)
	if err != nil {
		return 0, err
	}

	header := make([]byte, messageSizeInBytes)
	order.PutUint32(header, uint32(len(data)))
	return w.Write(append(header, data...))
}

// Send writes the json encoded value of v to the writer
func Send(w io.Writer, v interface{}, order binary.ByteOrder) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	return Write(w, bytes.NewReader(b), order)
}
