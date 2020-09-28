package mockserver

import (
	"fmt"
	"html"
	"net/http"
	"time"
)

// New constructs a mock server.
// Start it by calling LitenAndServe() and then Close() it of send a /quit request.
func New(addr string) *http.Server {
	svr := &http.Server{}
	svr.Addr = ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Closing server ...")
		fmt.Println("Closing server ...")
		go func() {
			time.Sleep(300 * time.Millisecond)
			svr.Close()

		}()
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Url requested : %q", html.EscapeString(r.URL.Path))
		fmt.Printf("Url requested : %q\n", html.EscapeString(r.URL.Path))
	})

	svr.Handler = mux
	return svr
}
