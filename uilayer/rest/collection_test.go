package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/tonyalaribe/ninja/datalayer"
)

func createCollection(t *testing.T, req string) {
	var collectionData NewCollectionVM
	err := json.Unmarshal([]byte(req), &collectionData)
	AssertEqual(t, err, nil)

	coreManager, mockDataStore, mockCtrler, err := GetCoreManager(t)
	if mockCtrler != nil {
		mockDataStore.EXPECT().CreateCollection(collectionData.Name, collectionData.Schema, collectionData.Meta).Return(nil).MinTimes(1)
		defer mockCtrler.Finish()
	}

	s := &Server{
		core: coreManager,
	}
	server := httptest.NewServer(ErrorWrapper(s.CreateCollection))
	defer server.Close()

	body := bytes.NewReader([]byte(req))
	resp, err := server.Client().Post(server.URL, "application/json", body)
	AssertEqual(t, err, nil)

	RespIsNotError(t, resp.Body)
}

func TestCreateCollection(t *testing.T) {
	req := fmt.Sprintf(`
	{
		"name": "%s", 
		"schema": {
			"title": "A registration form",
			"description": "A simple form example.",
			"type": "object",
			"required": [
				"firstName"
			],
			"properties": {
				"firstName": {
					"type": "string",
					"title": "First name"
				}
			}
		}, 
		"meta":{}
	}
	`, uuid.Must(uuid.NewV4()).String())
	createCollection(t, req)
}

func TestGetCollections(t *testing.T) {
	req := fmt.Sprintf(`
	{
		"name": "%s", 
		"schema": {
			"title": "A registration form",
			"description": "A simple form example.",
			"type": "object",
			"required": [
				"firstName"
			],
			"properties": {
				"firstName": {
					"type": "string",
					"title": "First name"
				}
			}
		}, 
		"meta":{}
	}
	`, uuid.Must(uuid.NewV4()).String())

	reqData := datalayer.CollectionVM{}
	err := json.Unmarshal([]byte(req), &reqData)
	AssertEqual(t, err, nil)

	createCollection(t, req)

	coreManager, mockDataStore, mockCtrler, err := GetCoreManager(t)
	if mockCtrler != nil {
		mockDataStore.EXPECT().GetCollections().Return([]datalayer.CollectionVM{reqData}, nil)
		defer mockCtrler.Finish()
	}

	s := &Server{
		core: coreManager,
	}
	server := httptest.NewServer(ErrorWrapper(s.GetCollections))
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	AssertEqual(t, err, nil)

	var buf bytes.Buffer
	body := io.TeeReader(resp.Body, &buf)
	RespIsNotError(t, body)

	var collections []datalayer.CollectionVM
	json.NewDecoder(&buf).Decode(&collections)
	if len(collections) < 1 {
		t.Errorf("empty collections list returned")
	}
}

func TestGetSchema(t *testing.T) {
	req := fmt.Sprintf(`
	{
		"name": "%s", 
		"schema": {
			"title": "A registration form",
			"description": "A simple form example.",
			"type": "object",
			"required": [
				"firstName"
			],
			"properties": {
				"firstName": {
					"type": "string",
					"title": "First name"
				}
			}
		}, 
		"meta":{}
	}
	`, uuid.Must(uuid.NewV4()).String())

	reqData := NewCollectionVM{}
	err := json.Unmarshal([]byte(req), &reqData)
	AssertEqual(t, err, nil)

	createCollection(t, req)

	coreManager, mockDataStore, mockCtrler, err := GetCoreManager(t)
	if mockCtrler != nil {
		mockDataStore.EXPECT().GetSchema(reqData.Name).Return(reqData.Schema, nil)
		defer mockCtrler.Finish()
	}

	s := &Server{
		core: coreManager,
	}
	r := chi.NewMux()
	r.Get("/{collectionName}", ErrorWrapper(s.GetSchema))
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/" + reqData.Name)
	AssertEqual(t, err, nil)

	RespIsNotError(t, resp.Body)
}
