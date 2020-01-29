package main

import (
	"testing"
)

func TestExtractStruct(t *testing.T) {
	s := `{"name":"John","age":30,"cars":{"car1":"Ford","car2":"BMW","car3":"Fiat"}}`

	if s == "" {
		t.Errorf("Expected a json string but got %s", s)
	}

	st, err := extractStruct(s)
	if err != nil {
		t.Errorf("extrated structs error. Got %s", st)
		t.Error(err)
	}
}

func TestCapFirst(t *testing.T) {
	s := `hello`
	res := capFirst(s)

	if res != `Hello` {
		t.Errorf("Expected `Hello`. Got %s", res)
	}
}
