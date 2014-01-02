package main 
//
import (
	"flag"
	"net"
	"log"
	"os"
	"os/signal"
	"syscall"
	"../chat"
)
//
func main() {
	// flag
	ip := flag.String("ip", "0.0.0.0", "server ip")
	port := flag.String("port", "9999", "server port")
	flag.Parse()
	address := *ip + ":" + *port
	go func() {
		chat.TestChannelServer(address)
	}()
	// 
	exit_chan := make(chan int)
	signal_chan := make(chan os.Signal, 1)
	go func() {
		<-signal_chan
		log.Println("Caught signal, exiting...")
		exit_chan <- 1
	}()
	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM)
	<- exit_chan
	// listen
/*
	listener, err := net.Listen("tcp", address)
	if nil != err {
		log.Println("Error listen: " , err.Error())
		return
	} else {
		log.Println("serve at: ", address)
	}
	// accept
	for {
		connection, err := listener.Accept()
		if nil != err {
			log.Println("Error accept: ", err.Error())
			return
		}
		go Handler(connection)
	}
*/
}
const (
	kRecvBufLen = 2048
)
// Handler
func Handler(connection net.Conn) {
	buf := make([]byte, kRecvBufLen)
	// read
	_, err := connection.Read(buf)
	if nil != err {
		log.Println("Error read: ", err.Error())
		return
	} else {
		log.Println("Read: ", buf)
	}
	// send
	_, err = connection.Write(buf)
	if nil != err {
		log.Println("Error write: ", err.Error())		
	} else {
		log.Println("Relay")
	}
}
