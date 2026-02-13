package store

import "github.com/lfcypo/scrkun/model"

func autoMigrate() error {
	return db.AutoMigrate(model.Models...)
}
