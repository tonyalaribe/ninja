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

const (
	TestSchema1 = `
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
	`
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
	server := httptest.NewServer(ResponseWrapper(s.CreateCollection))
	defer server.Close()

	body := bytes.NewReader([]byte(req))
	resp, err := server.Client().Post(server.URL, "application/json", body)
	AssertEqual(t, err, nil)

	RespIsNotError(t, resp.Body)
}

func TestCreateCollection(t *testing.T) {
	req := fmt.Sprintf(TestSchema1, uuid.Must(uuid.NewV4()).String())
	createCollection(t, req)
}

func TestGetCollections(t *testing.T) {
	req := fmt.Sprintf(TestSchema1, uuid.Must(uuid.NewV4()).String())

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
	server := httptest.NewServer(ResponseWrapper(s.GetCollections))
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	AssertEqual(t, err, nil)

	var buf bytes.Buffer
	body := io.TeeReader(resp.Body, &buf)
	RespIsNotError(t, body)

	var responseData ResponseResource
	err = json.NewDecoder(&buf).Decode(&responseData)
	AssertEqual(t, err, nil)

	collections, ok := responseData.Data.([]datalayer.CollectionVM)
	AssertEqual(t, ok, true)
	if len(collections) < 1 {
		t.Errorf("empty collections list returned")
	}
}

func TestGetSchema(t *testing.T) {
	req := fmt.Sprintf(TestSchema1, uuid.Must(uuid.NewV4()).String())

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
	r.Get("/{collectionName}/schema", ResponseWrapper(s.GetSchema))
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/" + reqData.Name + "/schema")
	AssertEqual(t, err, nil)

	RespIsNotError(t, resp.Body)
}

func TestSaveItem(t *testing.T) {
	req := fmt.Sprintf(TestSchema1, uuid.Must(uuid.NewV4()).String())

	reqData := NewCollectionVM{}
	err := json.Unmarshal([]byte(req), &reqData)
	AssertEqual(t, err, nil)

	createCollection(t, req)

	coreManager, mockDataStore, mockCtrler, err := GetCoreManager(t)
	AssertEqual(t, err, nil)
	s := &Server{
		core: coreManager,
	}
	r := chi.NewMux()
	r.Post("/{collectionName}", ResponseWrapper(s.SaveItem))
	server := httptest.NewServer(r)
	defer server.Close()

	itemBody := fmt.Sprintf(`
		{
			"_id":"%s",
			"firstName":"Anthony Alaribe"
		}
	`, uuid.Must(uuid.NewV4()).String())

	var itemData map[string]interface{}
	json.Unmarshal([]byte(itemBody), &itemData)

	if mockCtrler != nil {
		mockDataStore.EXPECT().GetSchema(reqData.Name).Return(reqData.Schema, nil)
		mockDataStore.EXPECT().SaveItem(reqData.Name, itemData["_id"].(string), itemData).Return(nil)
		defer mockCtrler.Finish()
	}

	resp, err := server.Client().Post(server.URL+"/"+reqData.Name, "application/json", bytes.NewBufferString(itemBody))
	AssertEqual(t, err, nil)

	RespIsNotError(t, resp.Body)
}
