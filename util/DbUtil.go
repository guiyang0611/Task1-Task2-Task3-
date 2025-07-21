package util

import (
	"fmt"
	"gorm.io/gorm"
)

func UpdateByModel(db *gorm.DB, model interface{}, whereMap map[string]interface{}, updateMap map[string]interface{}) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	tx := db.Debug().Model(model)

	if len(whereMap) > 0 {
		tx.Where(whereMap)
	}

	if len(updateMap) > 0 {
		tx.Updates(updateMap)
	}

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DeleteByModel(db *gorm.DB, model interface{}, query string, params ...interface{}) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if len(query) == 0 {
		return fmt.Errorf("sql is empty")
	}

	tx := db.Debug().Where(query, params...).Delete(model)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func CreateByModel(db *gorm.DB, module interface{}) {
	tx := db.Create(module)
	if tx.Error != nil {
		panic(tx.Error)
	}
}
