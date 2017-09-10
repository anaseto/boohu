package main

import (
	"crypto/rand"
	"log"
	"math/big"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func RandInt(n int) int {
	if n <= 0 {
		return 0
	}
	x, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		log.Println(err, ":", n)
	}
	return int(x.Int64())
}
