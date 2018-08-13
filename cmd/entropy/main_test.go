package main

import (
	"testing"
	"math/rand"
	"time"
)

func TestTRNG_ChiSquared(t *testing.T) {
	// Heuristic limit on chisquared value
	// Corresponds to ~95% confidence
	heuristic := 3.841
	sampleSize := 100000 // number of bytes to sample

	var g = NewTRNG(uint8(time.Now().UnixNano()))

	expectedFrequencyPerByte := 4.0
	expected := []float64{float64(sampleSize) * expectedFrequencyPerByte, float64(sampleSize) * expectedFrequencyPerByte}
	observed := []float64{0.0, 0.0}

	for i := 0; i < sampleSize; i++ {
		sample, gen := g.Next()
		g = gen

		z, o := CountBits(sample)
		observed[0] += float64(z)
		observed[1] += float64(o)
	}

	t.Logf("Frequencies: %v", observed)

	chisquared := ChiSquared(observed, expected)

	if chisquared > heuristic {
		t.Fatalf("Failed Hypothesis; Chi-Squared %f not less than %f", chisquared, heuristic)
	} else {
		t.Logf("Hypothesis Holds: Chi-Squared %f less than %f", chisquared, heuristic)
	}
}

// Test Go's RNG to sanity check test method
func TestTRNG_GoRand(t *testing.T) {
	// Heuristic limit on chisquared value
	// Corresponds to ~95% confidence
	heuristic := 3.841
	sampleSize := 10000 // number of bytes to sample

	expectedFrequencyPerByte := 4.0
	expected := []float64{float64(sampleSize) * expectedFrequencyPerByte, float64(sampleSize) * expectedFrequencyPerByte}
	observed := []float64{0.0, 0.0}

	for i := 0; i < sampleSize; i++ {
		sample := uint8(rand.Intn(255))

		z, o := CountBits(sample)
		observed[0] += float64(z)
		observed[1] += float64(o)
	}

	t.Logf("Frequencies: %v", observed)

	chisquared := ChiSquared(observed, expected)

	if chisquared > heuristic {
		t.Fatalf("Failed Hypothesis; Chi-Squared %f not less than %f", chisquared, heuristic)
	} else {
		t.Logf("Hypothesis Holds: Chi-Squared %f less than %f", chisquared, heuristic)
	}
}

func TestTRNG_Next(t *testing.T) {
	t.Run("Consecutive values should be equivalent", func(t *testing.T) {
		var g = NewTRNG(uint8(time.Now().UnixNano()))

		first, _ := g.Next()
		second, _ := g.Next()

		if first != second {
			t.Fatalf("Expected equal values, got {%d, %d}", first, second)
		}
	})
}

func TestCountZeros(t *testing.T) {
	t.Run("255", func(t *testing.T) {
		z, o := CountBits(255)
		if o != 8 {
			t.Fail()
		}
		if z != 0 {
			t.Fail()
		}
	})

	t.Run("0", func(t *testing.T) {
		z, o := CountBits(0)
		if o != 0 {
			t.Fail()
		}
		if z != 8 {
			t.Fail()
		}
	})

	t.Run("1", func(t *testing.T) {
		z, o := CountBits(1)
		if z != 7 {
			t.Fail()
		}
		if o != 1 {
			t.Fail()
		}
	})

	t.Run("7", func(t *testing.T) {
		z, o := CountBits(7)
		if z != 5 {
			t.Fail()
		}
		if o != 3 {
			t.Fail()
		}
	})

	t.Run("42", func(t *testing.T) {
		z, o := CountBits(42)
		if z != 5 {
			t.Fail()
		}
		if o != 3 {
			t.Fail()
		}
	})
}

/* Test Helpers */

// Calculate Chi Squared test statistic
func ChiSquared(exp []float64, sample []float64) float64 {
	var result = 0.0

	sampleSize := len(sample)
	for i := 0; i < sampleSize; i++ {
		if sample[i] == 0 && exp[i] == 0 {
			continue
		}
		result += (exp[i] - sample[i]) * (exp[i] - sample[i]) / exp[i]
	}
	return result
}

// return number of zeros and ones in the bit representation of a uint8
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
