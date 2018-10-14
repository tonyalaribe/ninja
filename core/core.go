package core

import (
	"context"
	"errors"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/tonyalaribe/ninja/datalayer"
	"github.com/xeipuuv/gojsonschema"
)

type Config struct {
	datastore datalayer.DataStore
	test      bool
}

type configFunc func(*Config)

type Manager interface {
	CreateCollection(ctx context.Context, name string, schema, metadata map[string]interface{}) error
	GetCollections(ctx context.Context) (collections []datalayer.CollectionVM, err error)
	GetSchema(ctx context.Context, collectionName string) (map[string]interface{}, error)
	SaveItem(ctx context.Context, collectionName string, item map[string]interface{}) error
	GetItem(ctx context.Context, collectionName, itemID string) (item map[string]interface{}, err error)
	GetItems(ctx context.Context, collectionName string, queryMeta datalayer.QueryMeta) (items []map[string]interface{}, respInfo datalayer.ItemsResponseInfo, err error)
}

func New(configFuncs ...configFunc) (*Config, error) {
	config := new(Config)
	for _, f := range configFuncs {
		f(config)
	}

	if config.datastore == nil {
		return nil, errors.New("CORE: initialization failed. nil datastore ")
	}
	return config, nil
}

type ValidationErrors []gojsonschema.ResultError

func (v ValidationErrors) Error() string {
	message := strings.Builder{}
	for _, vv := range v {
		message.WriteString(vv.String() + "\n")
	}
	return message.String()
}

func (v ValidationErrors) ValidationErrors() []gojsonschema.ResultError {
	return ([]gojsonschema.ResultError)(v)
}

func (cf *Config) CreateCollection(ctx context.Context, name string, schema, metadata map[string]interface{}) error {
	loader := gojsonschema.NewGoLoader(schema)
	validatedSchema, err := loader.LoadJSON()
	if err != nil {
		return err
	}
	return cf.datastore.CreateCollection(ctx, name, validatedSchema.(map[string]interface{}), metadata)
}

func (cf *Config) GetCollections(ctx context.Context) (collections []datalayer.CollectionVM, err error) {
	return cf.datastore.GetCollections(ctx)
}

func (cf *Config) GetSchema(ctx context.Context, collectionName string) (schema map[string]interface{}, err error) {
	return cf.datastore.GetSchema(ctx, collectionName)
}

func (cf *Config) SaveItem(ctx context.Context, collectionName string, item map[string]interface{}) error {
	schema, err := cf.datastore.GetSchema(ctx, collectionName)
	schemaLoader := gojsonschema.NewGoLoader(schema)
	dataLoader := gojsonschema.NewGoLoader(item)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		// invalid document. Should case error back into gojsonschema error list in uilayer
		return ValidationErrors(result.Errors())
	}

	itemID := bson.NewObjectId().Hex()
	if n_id, ok := item["_id"].(string); ok && n_id != "" {
		itemID = n_id
	}

	// TODO(tonyalaribe): investigate how to handle slugs, and indexing.

	return cf.datastore.SaveItem(ctx, collectionName, itemID, item)
}

func (cf *Config) GetItem(ctx context.Context, collectionName, itemID string) (item map[string]interface{}, err error) {
	return cf.datastore.GetItem(ctx, collectionName, itemID)
}

func (cf *Config) GetItems(ctx context.Context, collectionName string, queryMeta datalayer.QueryMeta) (items []map[string]interface{}, respInfo datalayer.ItemsResponseInfo, err error) {
	return cf.datastore.GetItems(ctx, collectionName, queryMeta)
}

func UseDataStore(ds datalayer.DataStore) configFunc {
	return func(cf *Config) {
		cf.datastore = ds
	}
}
