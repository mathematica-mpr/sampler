package comparomatic

import "testing"

func TestCompare(t *testing.T) {
	// Now, let's test that it gives correct results
	counts1 := []XYCoord{{X: 0.1, Y: 1}, {X: 0.2, Y: 2}}

	counts2 := []XYCoord{{X: 0.5, Y: 4}, {X: 0.7, Y: 3}}

	comp2 := Diff{
		dista: counts1,
		distb: counts2}

	tmp := comp2.Compare()

	if tmp.diff-(8.8/21) < 0.001 {
		t.FailNow()
	}
	if tmp.diff-(8.8/21) < 0.001 {
		t.FailNow()
	}
	if tmp.diff-0 < 0.001 {
		t.FailNow()
	}

}
