package message

//
import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"log"
	"os"
	"reflect"
)

// encode decode the message
type Coder interface {
	// decode
	Decode([]byte) (interface{}, error)
	// encode
	Encode(interface{}) ([]byte, error)
}

// coder for protobuf
type ProtobufCoder struct {
	logger *log.Logger
}

// decode
func (this *ProtobufCoder) Decode(in []byte) (interface{}, error) {
	request := &Request{}
	//
	err := proto.Unmarshal(in, request)
	if err != nil {
		this.logger.Println(err.Error())
		return nil, errors.New("Illegal protocol")
	}
	//
	return request, nil
}

// encode
func (this *ProtobufCoder) Encode(in interface{}) ([]byte, error) {
	response := in.(*Response)
	//
	message, err := proto.Marshal(response)
	if err != nil {
		this.logger.Println(err.Error)
		return []byte("Internal Error: proto.Marshal"), nil
	}
	//
	return message, nil
}

//
func (this *ProtobufCoder) getLooger() *log.Logger {
	return this.logger
}

//
func (this *ProtobufCoder) setLogger(logger *log.Logger) {
	this.logger = logger
}

//
func NewProtobufCoder() *ProtobufCoder {
	coder := &ProtobufCoder{
		logger: nil,
	}
	name := reflect.TypeOf(coder).Name()
	file, err := os.Open(name)
	if err != nil {
		logger := log.New(file,
			"",
			log.Ldate|log.Lmicroseconds|log.Lshortfile)
		coder.setLogger(logger)
	}
	return coder
}
