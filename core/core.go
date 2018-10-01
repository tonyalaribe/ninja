package core

import (
	"github.com/globalsign/mgo/bson"
	"github.com/tonyalaribe/ninja/datalayer"
)

type Config struct {
	datastore datalayer.DataStore
	test      bool
}

type configFunc func(*Config)

func New(configFuncs ...configFunc) *Config {
	config := new(Config)
	for _, f := range configFuncs {
		f(config)
	}
	return config
}

func (cf *Config) CreateCollection(name string, schema, metadata map[string]interface{}) error {
	return cf.datastore.CreateCollection(name, schema, metadata)
}

func (cf *Config) GetSchema(collectionName string) (map[string]interface{}, error) {
	return cf.datastore.GetSchema(collectionName)
}

func (cf *Config) SaveItem(collectionName string, item map[string]interface{}) error {
	itemID := bson.NewObjectId().Hex()
	// TODO(tonyalaribe): investigate how to handle slugs, and indexing.

	return cf.datastore.SaveItem(collectionName, itemID, item)
}

func (cf *Config) GetItem(collectionName, itemID string) (item map[string]interface{}, err error) {
	return cf.datastore.GetItem(collectionName, itemID)
}

func (cf *Config) GetItems(collectionName string, queryMeta datalayer.QueryMeta) (items []map[string]interface{}, respInfo datalayer.ItemsResponseInfo, err error) {
	return cf.datastore.GetItems(collectionName, queryMeta)
}
