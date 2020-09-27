package client

import (
	"fmt"
	"testing"
)

func TestConstructClient(t *testing.T) {

	conf := Configuration{
		ServerURL: "https://testerver.com/something", // no trailing slash - it will be added ...
		License:   "LI/../../../LI"}

	c := New(conf)
	fmt.Printf("%v\n", c)

}

func TestInvalidServer(t *testing.T) {
	conf := Configuration{
		ServerURL: "https://invalid.com/snonexistent", // no trailing slash - it will be added ...
		License:   "testing"}

	c := New(conf)
	c.queryServer()
	if c.Check() == nil {
		t.Fatalf("Check should have failed ?\n%v\n", c)
		t.FailNow()
	}
}
