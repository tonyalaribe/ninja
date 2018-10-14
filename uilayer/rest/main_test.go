package rest

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tonyalaribe/ninja/core"
	"github.com/tonyalaribe/ninja/datalayer"
	"github.com/tonyalaribe/ninja/datalayer/mock"
)

var dataStore datalayer.DataStore

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(PingPong))
	defer server.Close()
	resp, err := server.Client().Get(server.URL)
	AssertEqual(t, err, nil)

	bb, err := ioutil.ReadAll(resp.Body)
	AssertEqual(t, err, nil)
	AssertEqual(t, string(bb), "pong")
}

func GetCoreManager(t *testing.T) (coreManager core.Manager, mockDataStore *mock.MockDataStore, mockCtrler *gomock.Controller, err error) {
	if dataStore != nil {
		coreManager, err := core.New(core.UseDataStore(dataStore))
		return coreManager, nil, nil, err
	}

	mockCtrler = gomock.NewController(t)
	// NOTE: defer mockCtrler.Finish() on caller

	mockDataStore = mock.NewMockDataStore(mockCtrler)
	coreManager, err = core.New(core.UseDataStore(mockDataStore))
	return coreManager, mockDataStore, mockCtrler, err
}

func RespIsNotError(t *testing.T, resp io.Reader) {
	var respData ResponseResource
	err := json.NewDecoder(resp).Decode(&respData)
	if err != nil {
		// most likely an array type, so definitely not the generic error return message
		return
	}

	if respData.Error != "" {
		t.Errorf("got an error response: %v", respData)
	}
}

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}

	t.Helper()
	t.Errorf("Received %v (type %v), expected %v (type %v) \n %s", a, reflect.TypeOf(a), b, reflect.TypeOf(b), string(debug.Stack()))
}

// AssertNotEqual checks if values are not equal
func AssertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected to not equal %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
