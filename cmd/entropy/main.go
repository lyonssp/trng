package main

import (
	"time"
	"github.com/prometheus/common/log"
	"encoding/binary"
	"os"
	"runtime"
)

var (
	seed = MakeEntropy()
	trng = NewTRNG(seed)
)

func main() {
	var next uint8
	for true {
		next, trng = trng.Next()
		log.Debugf("Generated %d", next)
		if err := binary.Write(os.Stderr, binary.BigEndian, next); err != nil {
			panic(err)
		}
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

