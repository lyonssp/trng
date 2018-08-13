package main

import (
	"time"
	"encoding/binary"
	"os"
	"runtime"
)

var (
	seed = MakeEntropy()
	trng = NewTRNG(seed)
)

func main() {
	bytesWritten := 0
	fd, err := os.Create("./pool.txt")
	if err != nil {
		panic(err)
	}

	var next uint8
	for bytesWritten < 1024 {
		next, trng = trng.Next()
		if err := binary.Write(fd, binary.BigEndian, next); err != nil {
			panic(err)
		}
		bytesWritten += 1
	}
}

type TRNG struct {
	seed uint8
}

func (g *TRNG) Next() (uint8, TRNG) {
	next := g.seed
	return next, NewTRNG(MakeEntropy())
}

func NewTRNG(s uint8) TRNG {
	return TRNG{seed: s}
}

func MakeEntropy() uint8 {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	var result uint8 = 0x00
	for i := uint(0); i < 4; i++ {
		t := time.Now().UnixNano()
		timeBits := uint8(t & 0x3)
		result = result ^ (timeBits << (i << 1))
	}

	return result
}
