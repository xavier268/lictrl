package client

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the main object to manage the license-related permissions
type Client struct {
	start  time.Time  // time when Client was first constructed
	locked bool       // if set, application is definitively locked
	rnd    *rand.Rand // random generator for client
	ssid   int        // unique session id
	e      error      // last error recorded

	lastCheck time.Time     // the last successful check
	offLimit  time.Duration // max duration we can stay offline without a successful check

	license string   // license identification string
	surl    *url.URL // server url entry point

	ticker *time.Ticker // Ticker to trigger automatic checks.

	done chan bool // signal end of process to all concurrent goroutines
}

// New creates a new Client, using the provided configuration.
func New(conf Configuration) *Client {

	c := new(Client)
	c.start = time.Now()
	c.lastCheck = c.start
	c.offLimit = conf.OfflineLimit

	c.locked = false
	c.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	c.ssid = c.rnd.Int()

	c.license = url.PathEscape(conf.License)
	c.surl, c.e = url.Parse(conf.ServerURL)
	c.surl.Path = strings.Join([]string{c.surl.Path, c.license}, "/")

	c.done = make(chan bool)

	if conf.AutoRepeat != 0 {
		c.ticker = time.NewTicker(conf.AutoRepeat)
		go c.repeatChecks()
	}
	return c
}

// Close release all Client ressources.
// All subsequent checks will fail.
func (c *Client) Close() error {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.done <- true
	c.locked = true
	return nil
}

// repeatCkecks perform repeated checks based on ticker.
// This is called asynchroneously by constructor if auto mode was defined.
// Do not call directly.
func (c *Client) repeatChecks() {
	for {
		select {
		case <-c.done:
			fmt.Println("Close requested !")
			return
		case tt := <-c.ticker.C:
			fmt.Println("Auto-checking ..", tt)
			c.checkServer()
		}

	}
}

func (c *Client) String() string {
	return fmt.Sprintf("Client dump :\n============\nstarted:\t%v\nlast check:\t%v\noffline max:\t%v\n"+
		"locked:\t\t%v\nssid:\t\t%v\nerror:\t\t%v\nlicense:\t%s\nserver url:\t%v\n",
		c.start, c.lastCheck, c.offLimit,
		c.locked, c.ssid, c.e, c.license, c.surl)
}

// Check the validity of the license, returning error if invalid.
// What happens exactly depends on the configuration that was passed upon creation.
func (c *Client) Check() error {
	if c.locked {
		return fmt.Errorf("Application was locked, (%v)", c.e)
	}
	if c.e != nil {
		return c.e
	}

	// add more checks here ...

	return nil
}

// checkServer sends the Get query to the server and verify response is valid.
func (c *Client) checkServer() {
	resp, err := http.Get(c.surl.String())

	if err != nil || resp.StatusCode != http.StatusOK {
		if c.isOfflineOk() {
			// ignore network failure when within acceptable timout
			return
		}
		// error beyond acceptable limit !
		// Client was locked already by offlimit check.
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
		c.e = fmt.Errorf("Exceeded timeout without being able to connect to authetication server")
		return false
	}

	return true
}
