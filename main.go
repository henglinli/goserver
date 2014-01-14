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
	loginManager := login.NewLoginManager()
	// 
	db := "user.db"
	err := loginManager.OpenDB(db)
	if err != nil {
		return
	}
	defer loginManager.CloseDB()
	//
	sessionManager := server.NewAddrSessionManager()
	//
	handler := server.NewDemuxerHandler(loginManager, sessionManager)
	s := server.NewServer(address)
	s.Serve(handler)
	utils.Wait()
	s.Stop()
	//
	log.Println("Good Bye!")
}
