package main

import (
	"crypto/rand"
	"log"
	"bytes"
	"math/big"
)

func main() {
	c := 2
	bigint := big.NewInt(int64(c))
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("error:", err)
		return
	}
	log.Println(int(b))
	bigint.SetBytes(b)
	log.Println(bigint)
	// The slice should now contain random bytes instead of only zeroes.
	log.Println(bytes.Equal(b, make([]byte, c)))
}
