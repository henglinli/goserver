package main

import (
	"fmt"
)

func main() {
	ch := make(chan []byte)
	buf := [100]byte{}
	go func() {
		buffer := <-ch
		fmt.Println(len(buffer))
	}()
	ch <- buf[:]

	//
	copy(buf[0:5], string("hello"))
	fmt.Println(buf)
}
