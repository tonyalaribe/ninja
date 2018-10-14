package datalayer

import (
	"context"
	"errors"
	"log"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]DataStore)
)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver DataStore) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		log.Panic("datalayer: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		log.Panic("datalayer: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Connect(name string, dbConfig DBConfig) (DataStore, error) {
	driver := drivers[name]
	if driver == nil {
		return nil, errors.New("datalayer: No such driver available")
	}
	return driver.Connect(dbConfig)
}

type DBConfig struct {
	DriverType           string `mapstructure:"driver_type"` // eg mongodb, etc
	ConnectionString     string `mapstructure:"connection_string"`
	DatabaseName         string `mapstructure:"database_name"`
	SchemaCollectionName string `mapstructure:"schema_collection_name"` // where schemas will be stored.
}

type CollectionVM struct {
	Name   string
	Schema map[string]interface{}
	Meta   map[string]interface{}
}

//go:generate mockgen -destination=./mock/mock_datastore.go -package=mock github.com/tonyalaribe/ninja/datalayer DataStore
type DataStore interface {
	Connect(dbConfig DBConfig) (datastore DataStore, err error)
	CreateCollection(ctx context.Context, name string, schema, metadata map[string]interface{}) error
	GetCollections(ctx context.Context) (collections []CollectionVM, err error)
	GetSchema(ctx context.Context, collectionName string) (map[string]interface{}, error)
	SaveItem(ctx context.Context, collectionName, itemID string, item map[string]interface{}) error
	GetItem(ctx context.Context, collectionName, itemID string) (item map[string]interface{}, err error)
	GetItems(ctx context.Context, collectionName string, queryMeta QueryMeta) (items []map[string]interface{}, respInfo ItemsResponseInfo, err error)
}

type QueryMeta struct {
	Page        int
	Count       int
	QueryString string
}

type ItemsResponseInfo struct{}
