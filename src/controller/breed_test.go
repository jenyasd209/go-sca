package controller

import (
	"testing"
)

func TestBreedValidator(t *testing.T) {
	bv := NewBreedValidator()
	err := bv.Init()
	if err != nil {
		t.Logf("expected no errors, but got: %s\n", err)
		t.Fail()
	}
}
