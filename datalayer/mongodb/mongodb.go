package mongodb

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"github.com/tonyalaribe/ninja/datalayer"
)

type Datastore struct {
	DB               *mgo.Database
	SchemaCollection string
}

const DriverName = "mongodb"

func init() {
	datalayer.Register(DriverName, &Datastore{})
}

func NewDatastore(config datalayer.DBConfig) (*Datastore, error) {
	ds := Datastore{}
	session, err := mgo.Dial(config.ConnectionString)
	if err != nil {
		return nil, err
	}
	ds.DB = session.DB(config.DatabaseName)
	ds.SchemaCollection = config.SchemaCollectionName
	return &ds, nil
}

func (ds *Datastore) Connect(config datalayer.DBConfig) (datalayer.DataStore, error) {
	return NewDatastore(config)
}

type collectionData struct {
	Name     string                 `bson:"_id"`
	Schema   map[string]interface{} `bson:"schema"`
	MetaData map[string]interface{} `bson:"metadata"`
}

func (ds *Datastore) CreateCollection(ctx context.Context, name string, schema, metadata map[string]interface{}) error {
	data := collectionData{}
	data.Name = name
	data.Schema = schema
	data.MetaData = metadata
	err := ds.DB.C(ds.SchemaCollection).Insert(data)
	if err != nil {
		return errors.Wrap(err, "mongoDB: unable to create collection")
	}

	// TODO(tonyalaribe): make use of metadata
	return nil
}

func (ds *Datastore) GetCollections(ctx context.Context) (collections []datalayer.CollectionVM, err error) {
	err = ds.DB.C(ds.SchemaCollection).Find(nil).All(&collections)
	return collections, err
}

func (ds *Datastore) GetSchema(ctx context.Context, collectionName string) (map[string]interface{}, error) {
	result := collectionData{}
	err := ds.DB.C(ds.SchemaCollection).FindId(collectionName).One(&result)
	return result.Schema, errors.Wrap(err, "mongoDB: unable to get schema")
}

func (ds *Datastore) SaveItem(ctx context.Context, collectionName, itemID string, item map[string]interface{}) error {
	item["_id"] = itemID
	err := ds.DB.C(collectionName).Insert(item)
	return errors.Wrap(err, "mongoDB: unable to save item")
}

func (ds *Datastore) GetItem(ctx context.Context, collectionName, itemID string) (item map[string]interface{}, err error) {
	err = ds.DB.C(collectionName).FindId(itemID).One(&item)
	return item, errors.Wrap(err, "mongoDB: unable to get item")
}

func (ds *Datastore) GetItems(ctx context.Context, collectionName string, queryData datalayer.QueryMeta) (items []map[string]interface{}, respInfo datalayer.ItemsResponseInfo, err error) {
	err = ds.DB.C(collectionName).Find(nil).All(&items)
	return items, respInfo, errors.Wrap(err, "mongoDB: unable to items")
}
