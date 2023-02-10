package demo

import "testing"

func TestHTTPServer(t *testing.T) {
	StartWeb(":3002")
}
