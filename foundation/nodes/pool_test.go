package nodes

import "testing"

func TestGetMax3BelowTarget(t *testing.T) {
	target := 3
	expected := 2
	p := NewPool([]string{"ip1", "ip2"})

	got := p.GetMaxTargetRandomIPs(target)
	if len(got) != expected {
		t.Fatalf("Got %d. Expected %d\n", len(got), expected)
	}
}

func TestGetMax3OnTarget(t *testing.T) {
	target := 3
	expected := 3
	p := NewPool([]string{"ip1", "ip2", "ip3"})

	got := p.GetMaxTargetRandomIPs(target)
	if len(got) != expected {
		t.Fatalf("Got %d. Expected %d\n", len(got), expected)
	}
}

func TestGetMax3AboveTarget(t *testing.T) {
	target := 3
	expected := 3
	p := NewPool([]string{"ip1", "ip2", "ip3", "ip4"})

	got := p.GetMaxTargetRandomIPs(target)
	if len(got) != expected {
		t.Fatalf("Got %d. Expected %d\n", len(got), expected)
	}
}

func TestGetMax3NoDuplicates(t *testing.T) {
	target := 3
	p := NewPool([]string{"ip1", "ip2", "ip3", "ip4"})

	got := p.GetMaxTargetRandomIPs(target)
	duplicateSet := make(map[string] struct{})

	for _, ip := range got {
		duplicateSet[ip] = struct{}{}
	}

	if len(got) != len(duplicateSet) {
		t.Fatalf("Duplicates found.")
	}
}
