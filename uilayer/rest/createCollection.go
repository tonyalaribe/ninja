package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type NewCollectionVM struct {
	Name   string
	Meta   map[string]interface{}
	Schema map[string]interface{}
}

func (server *Server) CreateCollection(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	resource := NewCollectionVM{}
	err = json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: CreateCollection failed")
	}

	err = server.core.CreateCollection(resource.Name, resource.Schema, resource.Meta)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: CreateCollection failed")
	}

	return http.StatusOK, nil
}

func (server *Server) GetSchema(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	collectionName := chi.URLParam(r, "collectionName")

	schema, err := server.core.GetSchema(collectionName)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: GetSchema failed")
	}
	render.JSON(w, r, schema)
	return http.StatusOK, nil
}

/*
func (server *Server) SaveItem(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	collectionName := chi.Param("collectionName")

	resource := GetSchemaVM{}
	err = json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: GetSchema failed")
	}

	err = server.core.GetSchema(resource.Name, resource.Schema, resource.Meta)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: GetSchema failed")
	}

	return http.StatusOK, nil
}
*/

/*
	CreateCollection(name string, schema, metadata map[string]interface{}) error
	GetSchema(collectionName string) (map[string]interface{}, error)
	SaveItem(collectionName, itemID string, item map[string]interface{}) error
	GetItem(collectionName, itemID string) (item map[string]interface{}, err error)
	GetItems(collectionName string, queryMeta QueryMeta) (items []map[string]interface{}, respInfo ItemsResponseInfo, err error)
*/
