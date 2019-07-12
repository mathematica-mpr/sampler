package simulate

import (
	sampler "sampler/simulate"
	"testing"
)

// TestSimulate : testing how long it takes to run Simulate()
func TestSimulate(b *testing.B) {

	for i := 0; i < b.N; i++ {
		sampler.Simulate(1000, 121356654, 39, 14580, 41, 61, 100)
	}
}
