package chat

//
import (
	"log"
	"net"
)

//
type Handler interface {
	Read(message string) error
	Write(message string) error
}

//
type Channel struct {
	identity  int64
	head      *Client
	tail      *Client
	from_head chan string
	from_tail chan string
}

//
func (channel *Channel) Join(conn net.Conn) {
	switch {
	case nil == channel.head:
		channel.head = NewClient(conn)
		go func() {
			if channel.tail != nil {
				channel.listen()
			}
			for {
				// recv from head
				channel.from_head <- <-channel.head.incoming
				//log.Println("message from head");
			}
		}()
	case nil == channel.tail:
		channel.tail = NewClient(conn)
		go func() {
			if channel.tail != nil {
				channel.listen()
			}
			for {
				// recv from tail
				channel.from_tail <- <-channel.tail.incoming
				//log.Println("message from tail");
			}
		}()
	default:
		// not going here
		conn.Close()
	}
}

//
func (channel *Channel) listen() {
	go func() {
		for {
			select {
			case to_tail := <-channel.from_head:
				// forward to tail
				channel.tail.outgoing <- to_tail
			case to_head := <-channel.from_tail:
				// forward to head
				channel.head.outgoing <- to_head
			}
		}
	}()
}

//
func NewChannel(id int64) *Channel {
	channel := &Channel{
		identity:  id,
		head:      nil,
		tail:      nil,
		from_head: make(chan string),
		from_tail: make(chan string),
	}
	return channel
}

//
func TestChannelServer(address string) error {
	log.Println("testing Channel...")
	log.Println("open two terminal and run [telnet localhost 9999]")
	channels := make([]*Channel, 0)
	//
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	var counter int64
	counter = 0
	for {
		channel := NewChannel(counter)
		channels = append(channels, channel)
		conn1, err1 := listener.Accept()
		if err != nil {
			log.Println("Accept error: ", err1.Error())
		}
		channel.Join(conn1)
		//
		conn2, err2 := listener.Accept()
		if err2 != nil {
			log.Println("Accept error: ", err2.Error())
		}
		channel.Join(conn2)
		//
		counter = counter + 1
		//
		log.Println("channel connected")
	}
}
