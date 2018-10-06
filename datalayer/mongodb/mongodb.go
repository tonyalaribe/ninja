package mongodb

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"github.com/tonyalaribe/ninja/datalayer"
)

type datastore struct {
	db               *mgo.Database
	schemaCollection string
}

const DriverName = "mongodb"

func init() {
	datalayer.Register(DriverName, &datastore{})
}

func NewDatastore(config datalayer.DBConfig) (*datastore, error) {
	ds := datastore{}
	session, err := mgo.Dial(config.ConnectionString)
	if err != nil {
		return nil, err
	}
	ds.db = session.DB(config.DatabaseName)
	ds.schemaCollection = config.SchemaCollectionName
	return &ds, nil
}

func (ds *datastore) Connect(config datalayer.DBConfig) (datalayer.DataStore, error) {
	return NewDatastore(config)
}

func (ds *datastore) CreateCollection(name string, schema, metadata map[string]interface{}) error {
	schema["_id"] = name
	err := ds.db.C(ds.schemaCollection).Insert(schema)
	if err != nil {
		return errors.Wrap(err, "mongodb: unable to create collection")
	}

	// TODO(tonyalaribe): make use of metadata
	return nil
}

func (ds *datastore) GetSchema(collectionName string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	err := ds.db.C(ds.schemaCollection).FindId(collectionName).One(&result)
	return result, errors.Wrap(err, "mongodb: unable to get schema")
}

func (ds *datastore) SaveItem(collectionName, itemID string, item map[string]interface{}) error {
	item["_id"] = itemID
	err := ds.db.C(collectionName).Insert(item)
	return errors.Wrap(err, "mongodb: unable to save item")
}

func (ds *datastore) GetItem(collectionName, itemID string) (item map[string]interface{}, err error) {
	err = ds.db.C(collectionName).FindId(itemID).One(&item)
	return item, errors.Wrap(err, "mongodb: unable to get item")
}

func (ds *datastore) GetItems(collectionName string, queryData datalayer.QueryMeta) (items []map[string]interface{}, respInfo datalayer.ItemsResponseInfo, err error) {
	err = ds.db.C(collectionName).Find(bson.M{}).All(&items)
	return items, respInfo, errors.Wrap(err, "mongodb: unable to items")
}
