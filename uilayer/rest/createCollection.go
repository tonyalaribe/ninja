package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type NewCollectionVM struct {
	Name   string                 `json:"name"`
	Meta   map[string]interface{} `json:"meta"`
	Schema map[string]interface{} `json:"schema"`
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

	render.JSON(w, r, ResponseMessage(http.StatusOK, "Collection created successfully"))
	return http.StatusOK, nil
}

func (server *Server) GetCollections(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	collections, err := server.core.GetCollections()
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: GetCollections failed")
	}
	render.JSON(w, r, collections)
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

func (server *Server) SaveItem(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	collectionName := chi.URLParam(r, "collectionName")

	resource := map[string]interface{}{}
	err = json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: SaveItem failed")
	}

	err = server.core.SaveItem(collectionName, resource)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "REST: SaveItem failed")
	}

	render.JSON(w, r, ResponseMessage(http.StatusOK, "Saved Item Successfully"))
	return http.StatusOK, nil
}
