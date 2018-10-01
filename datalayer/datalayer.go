package datalayer

type DataStore interface {
	CreateCollection(name string, schema, metadata map[string]interface{}) error
	GetSchema(collectionName string) (map[string]interface{}, error)
	SaveItem(collectionName, itemID string, item map[string]interface{}) error
	GetItem(collectionName, itemID string) (item map[string]interface{}, err error)
	GetItems(collectionName string, queryMeta QueryMeta) (items []map[string]interface{}, respInfo ItemsResponseInfo, err error)
}

type QueryMeta struct {
	Page        int
	Count       int
	QueryString string
}

type ItemsResponseInfo struct {
}
