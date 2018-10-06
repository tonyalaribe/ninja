package rest

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tonyalaribe/ninja/mocks"
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
	server := httptest.NewServer(ErrorWrapper(PingPong))
	// Close the server when test finishes
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	AssertEqual(t, err, nil)

	bb, err := ioutil.ReadAll(resp.Body)
	AssertEqual(t, err, nil)
	fmt.Println(string(bb))
}

func TestCreateCollection(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockManager := mocks.NewMockManager(mockCtrl)
	s := &Server{
		core: mockManager,
	}
	server := httptest.NewServer(ErrorWrapper(s.CreateCollection))

	// Close the server when test finishes
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	AssertEqual(t, err, nil)

	bb, err := ioutil.ReadAll(resp.Body)
	AssertEqual(t, err, nil)
	fmt.Println(string(bb))
}
