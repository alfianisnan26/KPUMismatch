package dao

import "kawalrealcount/internal/data/model"

type Database interface {
	PutReplaceData(entity model.HHCWEntity) error
}
