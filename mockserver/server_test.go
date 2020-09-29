package mockserver

import (
	"testing"
	"time"
)

func TestMockServer(t *testing.T) {

	port := 16524
	srv := New(port)

	go srv.ListenAndServe()
	time.Sleep(300 * time.Millisecond)
	srv.Close()

}
