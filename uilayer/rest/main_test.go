package rest

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"
)

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func TestPing(t *testing.T) {
	// s := &Server{
	// core: nil,
	// }
	// Start a local HTTP server
	server := httptest.NewServer(ErrorWrapper(PingPong))
	// Close the server when test finishes
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	AssertEqual(t, err, nil)

	bb, err := ioutil.ReadAll(resp.Body)
	AssertEqual(t, err, nil)
	fmt.Println(string(bb))
}
