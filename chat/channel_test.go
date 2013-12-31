package chat

import (
	"log"
	"net"
	"testing"
)

//
func TestChannel(t *testing.T) {
	log.Println("testing Channel...")
	log.Println("open two terminal and run [telnet localhost 9999]")
	address := ":9999"
	channels := make([]*Channel, 0)
	//
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
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
