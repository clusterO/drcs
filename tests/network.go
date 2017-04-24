package tests

import (
	network "DCRS/net"
	"testing"
)

func TestDial(t *testing.T) {
	expected := "localhost:8181"
	result := network.Dial("localhost", "8181")

	if result != expected {
		t.Errorf("Expected server address %s, but got %s", expected, result)
	}
}
