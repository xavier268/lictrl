package client

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// Client is the main object to manage the license-related permissions
type Client struct {
	start     time.Time     // time when Client was first constructed
	locked    bool          // if set, application is definitively locked
	rnd       *rand.Rand    // random generator for client
	ssid      int           // unique session id
	lastError error         // last error recorded
	lastCheck time.Time     // the last successful check
	offLimit  time.Duration // max duration we can stay offline without a successful check
	minLimit  time.Duration // minimal time between 2 actual server requests, internally set to 1/3 of offLimit
	license   string        // license identification string
	surl      string        // server  entry point
	ticker    *time.Ticker  // Ticker to trigger automatic checks. minLimit does not apply here.
	done      chan bool     // signal end of process to all concurrent goroutines by closing this channel. Never send data to it !
}

// New creates a new Client, using the provided configuration.
func New(conf Configuration) *Client {

	c := new(Client)
	c.start = time.Now()
	c.lastCheck = c.start
	c.offLimit = conf.OfflineLimit
	c.minLimit = c.offLimit / 3

	c.locked = false
	c.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	c.ssid = c.rnd.Int()

	c.license = url.PathEscape(conf.License)
	c.surl = conf.ServerURL + "/" + c.license

	c.done = make(chan bool)

	if conf.AutoRepeat != 0 {
		c.ticker = time.NewTicker(conf.AutoRepeat)
		go c.repeatChecks()
	}
	return c
}

// Close release all Client ressources and locks client.
// All subsequent checks will fail.
func (c *Client) Close() error {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	close(c.done)
	c.locked = true
	return nil
}

// repeatCkecks perform repeated checks based on ticker.
// This is called asynchroneously by constructor if auto mode was defined.
// Do not call directly.
func (c *Client) repeatChecks() {
	for {
		select {
		case <-c.done: // When the channel is closed, this will receive the zero-value, false.
			fmt.Println("Close requested !")
			return
		case tt := <-c.ticker.C:
			fmt.Println("Auto-checking ..", tt)
			c.checkServer()
		}

	}
}

// Locked verify if Client was locked.
// No actual additionnal check is done.
func (c *Client) Locked() bool {
	return c.locked
}

// Check the validity of the license, asynchroneously.
// If already locked, return immediately.
// What happens exactly depends on the configuration that was passed upon creation.
func (c *Client) Check() {
	if c.locked {
		return
	}
	if c.lastError != nil {
		c.locked = true
		return
	}

	// do the actual asynchroneous check, only if enough time has elapsed
	if time.Now().After(c.lastCheck.Add(c.minLimit)) {
		go c.checkServer()
	}

}

// checkServer sends the Get query to the server and verify response is valid.
// it should never be called synchroneously to avoid freezing the application.
func (c *Client) checkServer() {
	resp, err := http.Get(c.surl)
	if err != nil && resp != nil {
		fmt.Printf("Server responded : %s\n", resp.Status)
	} else {
		fmt.Printf("Server error : %v\n", err)
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		if c.isOfflineOk() {
			// ignore network failure when within acceptable timout
			return
		}
		// error beyond acceptable limit !
		// Client was locked already by offlimit check, but let's be safe ...
		c.locked = true
		return
	}
	// additionnal checks here ...

	// everything is fine, update time of valid check
	c.lastCheck = time.Now()
	return
}

// isOfflineOk tells if it is still ok to be offline.
// If not, Client is also locked.
func (c *Client) isOfflineOk() bool {

	if time.Now().After(c.lastCheck.Add(c.offLimit)) {
		c.locked = true
		c.lastError = fmt.Errorf("Exceeded timeout without being able to connect to authetication server")
		return false
	}
	return true
}
