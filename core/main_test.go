package main

import (
	"fmt"
	"testing"

	"github.com/apxxxxxxe/Bouyomi/data"
)

func TestList(t *testing.T) {
	voices, err := data.ListVoices(true)
	if err != nil {
		t.Error(err)
	}
	for _, v := range voices {
		fmt.Printf("%v,%v\n", v.BouyomiChanNumber, v.Name)
	}
}
