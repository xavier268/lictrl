package client

import (
	"fmt"
	"testing"
	"time"
)

func TestConstructClient(t *testing.T) {

	conf := Configuration{
		ServerURL: "https://testerver.com/something", // no trailing slash - it will be added ...
		License:   "LI/../../../LI"}

	c := New(conf)

	if c.surl.String() != "https://testerver.com/something/LI%252F..%252F..%252F..%252FLI" {
		t.Fatalf("Unexpected client server url\n%v\n", c)
	}

	fmt.Printf("%v\n", c)

}

func TestInvalidServer0sec(t *testing.T) {
	conf := Configuration{
		ServerURL: "https://invalid.com/snonexistent", // no trailing slash - it will be added ...
		License:   "testing"}

	c := New(conf)
	c.checkServer()
	if c.Check() == nil {
		t.Fatalf("Check should have failed ?\n%v\n", c)
		t.FailNow()
	}
}

func TestInvalidServer1sec(t *testing.T) {
	conf := Configuration{
		ServerURL:    "https://invalid.com/snonexistent", // no trailing slash - it will be added ...
		License:      "testing",
		OfflineLimit: 1 * time.Second}

	c := New(conf)
	c.checkServer()
	if c.Check() != nil {
		t.Fatalf("Check should NOT have failed ?\n%v\n", c)
		t.FailNow()
	}
}

func TestReapeatChecks(t *testing.T) {
	conf := Configuration{
		ServerURL:    "http://www.google.com", // no trailing slash - it will be added ...
		OfflineLimit: 5 * time.Second,
		AutoRepeat:   2 * time.Second}

	c := New(conf)
	time.Sleep(20 * time.Second)
	c.Close()
	time.Sleep(5 * time.Second)
	fmt.Println("Closed successful")
}
