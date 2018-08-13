package main

import (
	"time"
	"github.com/prometheus/common/log"
	"encoding/binary"
	"os"
	"runtime"
)

var (
	seed = GenSeed()
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
	return next, NewTRNG(GenSeed())
}

func NewTRNG(s uint8) TRNG {
	return TRNG{seed: s}
}

func GenSeed() uint8 {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	var result uint8 = 0x00
	for i := uint(0); i < 4; i++ {
		time.Sleep(20 * time.Nanosecond)
		t := time.Now().UnixNano()
		timeBits := uint8(t & 0x3)
		result = result ^ (timeBits << (i << 1))
	}

	return result
}

func CountBits(bits uint8) (int, int) {
	ones := 0
	zeroes := 0
	for i := uint(0); i < 8; i++ {
		next := (bits >> i) & 1
		if next == 0 {
			zeroes += 1
		}
		if next == 1 {
			ones += 1
		}
	}

	return zeroes, ones
}
