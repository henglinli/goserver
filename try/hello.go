package main

import (
	"fmt"
	"net"
)

func main() {
	counter := 0
	counter = counter + 1
	ch := make(chan string)
	go func() {
		ch <- string("hello")
	}()
	hello := <-ch
	fmt.Println(hello)
	//
	listener, _ := net.Listen("tcp", ":6666")
	for {
		one, _ := listener.Accept()
		//other, _:= listener.Accept()
		//a, b := net.Pipe()
		one.Close()
	}
}
