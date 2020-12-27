package main

import (
	"testing"
)


func TestFindEncryptionKey(t *testing.T) {
	cardPubKey := 5764801
	doorPubKey := 17807724

	got := FindEncryptionKey(doorPubKey, cardPubKey)
	expected := 14897079

	if got != expected {
		t.Errorf("Got %d expected %d", got, expected)
	}
}
