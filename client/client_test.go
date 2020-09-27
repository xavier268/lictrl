package client

import (
	"fmt"
	"testing"
	"time"
)

func (c *Client) String() string {
	return fmt.Sprintf("Client dump :\n============\nstarted:\t%v\nlast check:\t%v\noffline max:\t%v\n"+
		"locked:\t\t%v\nssid:\t\t%v\nerror:\t\t%v\nlicense:\t%s\nserver url:\t%v\n",
		c.start, c.lastCheck, c.offLimit,
		c.locked, c.ssid, c.lastError, c.license, c.surl)
}

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
	if !c.Locked() {
		t.Fatalf("Should be locked ?\n%v\n", c)
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
	if c.Locked() {
		t.Fatalf("Should NOT be locked ?\n%v\n", c)
		t.FailNow()
	}
}

func TestReapeatChecks(t *testing.T) {
	conf := Configuration{
		ServerURL:    "https://github.com", // no trailing slash - it will be added ...
		OfflineLimit: 4 * time.Second,
		AutoRepeat:   1 * time.Second}

	c := New(conf)
	if c.Locked() {
		t.FailNow()
	}
	time.Sleep(10 * time.Second)
	if c.Locked() {
		t.FailNow()
	}
	c.Close() // closing will lock the Client immediately.
	if !c.Locked() {
		t.FailNow()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Closed successful")
	if !c.Locked() {
		t.FailNow()
	}
}
