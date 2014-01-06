package main

//
import (
	"./server"
	"./utils"
	"./login"
	"flag"
	"log"
)

//
//
type EchoMessage struct {
	// nil
}

func (handler *EchoMessage) Handle(message []byte) []byte {
	// nil
	return message
}

func main() {
	// flag
	ip := flag.String("ip", "0.0.0.0", "server ip")
	port := flag.String("port", "9999", "server port")
	flag.Parse()
	address := *ip + ":" + *port
	//
	log.Println("starting Server...")
	//
	manager := login.NewLoginManager()
	
	handler := server.NewDemuxerHandler(manager)
	s := server.NewServer(address)
	s.Serve(handler)
	utils.Wait()
	s.Stop()
	//
	log.Println("Good Bye!")
}
