package mapping

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type MongoCollectionToSQLDBTable struct {
	Name string `json:"name"`
	// MongoDBName string `json:"mongodb_dbname"`
	Collection string `json:"mongodb_collection"`
	// SQLDBName string `json:"sqldb_dbname"`
	Table string `json:"sqldb_table"`
	ColumnMapping map[string]string `json:"column_mapping"`
}

func FindItemByCollection(mappings map[string]*MongoCollectionToSQLDBTable, collection string) *MongoCollectionToSQLDBTable  {
	for _, item := range mappings {
		if item.Collection == collection {
			return item
		}
	}

	return nil
}

// LoadMapping load mapping from file
func LoadMapping(files []string) (map[string]*MongoCollectionToSQLDBTable, error) {
	mappingList := make(map[string]*MongoCollectionToSQLDBTable, len(files))
	for _, filePath := range files {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		mapping := new(MongoCollectionToSQLDBTable)
		err = json.Unmarshal(content, mapping)
		if err != nil {
			return nil, err
		}
		if _, exists := mappingList[mapping.Name]; exists {
			return nil, errors.New("duplicate mapping name: " + filePath)
		}
		if len(mapping.ColumnMapping) == 0 {
			return nil, errors.New("column_mapping length should be greater than 0: " + filePath)
		}
		mappingList[mapping.Name] = mapping
	}

	return mappingList, nil
}