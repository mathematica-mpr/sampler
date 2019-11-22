package compare

import (
	"sampler/simulate"
	"testing"
)

func TestCompare(t *testing.T) {
	// Now, let's test that it gives correct results
	counts1 := []simulate.Coord{{X: 0.1, Y: 1, C: 1}, {X: 0.2, Y: 2, C: 3}}

	counts2 := []simulate.Coord{{X: 0.5, Y: 4, C: 1}, {X: 0.7, Y: 3, C: 4}}

	comp2 := Diff{
		Dista: counts1,
		Distb: counts2}

	tmp := Compare(comp2)

	if tmp.Diff-(8.8/21) > 0.001 {
		t.FailNow()
	}
	if tmp.Less-(8.8/21) > 0.001 {
		t.FailNow()
	}
	if tmp.More-0 > 0.001 {
		t.FailNow()
	}
}
