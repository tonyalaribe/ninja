package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/tonyalaribe/ninja/datalayer"
)

type NewCollectionVM struct {
	Name   string                 `json:"name"`
	Meta   map[string]interface{} `json:"meta"`
	Schema map[string]interface{} `json:"schema"`
}

func (server *Server) CreateCollection(w http.ResponseWriter, r *http.Request) (responseData interface{}, statusCode int, err error) {
	resource := NewCollectionVM{}
	err = json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "REST: CreateCollection failed")
	}

	err = server.core.CreateCollection(r.Context(), resource.Name, resource.Schema, resource.Meta)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "REST: CreateCollection failed")
	}

	return "Collection created successfully", http.StatusOK, nil
}

func (server *Server) GetCollections(w http.ResponseWriter, r *http.Request) (responseData interface{}, statusCode int, err error) {
	collections, err := server.core.GetCollections(r.Context())
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "REST: GetCollections failed")
	}
	return collections, http.StatusOK, nil
}

func (server *Server) GetSchema(w http.ResponseWriter, r *http.Request) (responseData interface{}, statusCode int, err error) {
	collectionName := chi.URLParam(r, "collectionName")

	schema, err := server.core.GetSchema(r.Context(), collectionName)
	if err != nil {
		return nil, http.StatusNotFound, errors.Wrap(err, "REST: GetSchema failed")
	}
	return schema, http.StatusOK, nil
}

func (server *Server) SaveItem(w http.ResponseWriter, r *http.Request) (responseData interface{}, statusCode int, err error) {
	collectionName := chi.URLParam(r, "collectionName")

	resource := map[string]interface{}{}
	err = json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "REST: SaveItem failed")
	}

	err = server.core.SaveItem(r.Context(), collectionName, resource)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "REST: SaveItem failed")
	}
	return "Saved Item Successfully", http.StatusOK, nil
}

func (server *Server) GetItem(w http.ResponseWriter, r *http.Request) (responseData interface{}, statusCode int, err error) {
	collectionName := chi.URLParam(r, "collectionName")
	itemID := chi.URLParam(r, "itemID")

	item, err := server.core.GetItem(r.Context(), collectionName, itemID)
	if err != nil {
		return item, http.StatusInternalServerError, errors.Wrap(err, "REST: GetItem failed")
	}

	return item, http.StatusOK, nil
}

type ItemsResponse struct {
	Items []map[string]interface{}
	Meta  datalayer.ItemsResponseInfo
}

func (server *Server) GetItems(w http.ResponseWriter, r *http.Request) (responseData interface{}, statusCode int, err error) {
	collectionName := chi.URLParam(r, "collectionName")

	query := datalayer.QueryMeta{}
	items, respInfo, err := server.core.GetItems(r.Context(), collectionName, query)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "REST: GetItem failed")
	}

	return ItemsResponse{
		Items: items,
		Meta:  respInfo,
	}, http.StatusOK, nil
}

// TODO: Get single Item
// TODO: Get paginated list of items
// Item list should include meta. Get collections should be adjusted to return meta as well. eg
/*
{
	Items  []Items,
	Count int
	PerPage int
	ItemsSkipped int
	PagesCount int
	TotalCount int
}
*/
