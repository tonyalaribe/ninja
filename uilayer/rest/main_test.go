package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tonyalaribe/ninja/core"
	"github.com/tonyalaribe/ninja/datalayer"
	"github.com/tonyalaribe/ninja/datalayer/mock"
)

var dataStore datalayer.DataStore

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

// AssertNotEqual checks if values are not equal
func AssertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected to not equal %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(ErrorWrapper(PingPong))
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

func TestCreateCollection(t *testing.T) {
	req := `
	{
		"name": "testcollection", 
		"schema": {
			"title": "A registration form",
			"description": "A simple form example.",
			"type": "object",
			"required": [
				"firstName",
				"lastName"
			],
			"properties": {
				"firstName": {
					"type": "string",
					"title": "First name"
				},
				"lastName": {
					"type": "string",
					"title": "Last name"
				},
				"age": {
					"type": "integer",
					"title": "Age"
				},
				"bio": {
					"type": "string",
					"title": "Bio"
				},
				"password": {
					"type": "string",
					"title": "Password",
					"minLength": 3
				},
				"telephone": {
					"type": "string",
					"title": "Telephone",
					"minLength": 10
				}
			}
		}, 
		"meta":{}
	}
	`
	var collectionData NewCollectionVM
	err := json.Unmarshal([]byte(req), &collectionData)
	AssertEqual(t, err, nil)
	coreManager, mockDataStore, mockCtrler, err := GetCoreManager(t)
	if mockCtrler != nil {
		mockDataStore.EXPECT().CreateCollection(collectionData.Name, collectionData.Schema, collectionData.Meta)
		defer mockCtrler.Finish()
	}

	s := &Server{
		core: coreManager,
	}
	server := httptest.NewServer(ErrorWrapper(s.CreateCollection))

	// Close the server when test finishes
	defer server.Close()
	body := bytes.NewReader([]byte(req))
	resp, err := server.Client().Post(server.URL, "application/json", body)
	AssertEqual(t, err, nil)

	RespIsNotError(t, resp.Body)
}

func RespIsNotError(t *testing.T, resp io.Reader) {
	respData := map[string]interface{}{}
	err := json.NewDecoder(resp).Decode(&respData)
	AssertEqual(t, err, nil)

	_, ok := respData["error"]
	if ok {
		t.Errorf("got an error response: %v", respData)
	}

}
