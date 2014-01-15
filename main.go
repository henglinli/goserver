package main

//
import (
	"./client"
	"./server"
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
	log.Println("starting GoServer ...")
	//
	//handler := server.NewEchoMessageHandler()
	handler := client.NewProtobufClientManager()
	server.Go(address, handler)
	//
	/*

		//
		db := "user.db"
		err := clientManager.OpenDB(db)
		if err != nil {
			return
		}
		defer clientManager.CloseDB()
		//
		sessionManager := server.NewAddrSessionManager()
		//
		handler := server.NewDemuxerHandler(clientManager, sessionManager)
		s := server.NewGoServer(address)
		s.Serve(handler)
		utils.Wait()
		s.Stop()
	*/
	//
	log.Println("Good Bye!")
}
