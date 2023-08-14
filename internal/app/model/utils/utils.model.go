package utils

import (
	"context"
	"transfer/api/types/schema"

	"gorm.io/gorm"
)

func WrapPageQuery(ctx context.Context, db *gorm.DB, pp schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {

	total, err := FindPage(ctx, db, pp, out)
	if err != nil {
		return nil, err
	}

	return &schema.PaginationResult{
		Total:    total,
		PageNo:   pp.GetPageNo(),
		PageSize: pp.GetPageSize(),
	}, nil
}

func FindPage(ctx context.Context, db *gorm.DB, pp schema.PaginationParam, out interface{}) (int64, error) {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	} else if count == 0 {
		return count, nil
	}

	current, pageSize := pp.GetPageNo(), pp.GetPageSize()
	if current > 0 && pageSize > 0 {
		db = db.Offset((current - 1) * pageSize).Limit(pageSize)
	} else if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	err = db.Find(out).Error
	return count, err
}
