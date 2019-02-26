package zgo_db_mysql

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 获取数据库链接对象
func GORMDB(name string) *gorm.DB {
	return dbs[name]
}

// 添加
func Insert(db *gorm.DB, in interface{}) error {
	if ok := db.NewRecord(in); ok {
		if err := db.Create(&in).Error; err != nil {
			return err
		}
	} else {
		return errors.New("NewRecord false 添加失败")
	}
	return nil
}

// 修改
func Update(db *gorm.DB, in interface{}) error {
	if ok := db.NewRecord(in); ok {
		if err := db.Update(&in).Error; err != nil {
			return err
		}
	} else {
		return errors.New("NewRecord false 添加失败")
	}
	return nil
}

// 删除

// 查询

// 查询第一个
func First(ctx context.Context, db *gorm.DB, out interface{}) chan interface{} {
	db.Table("spider.table")
	outch := make(chan interface{})
	go func(ctx context.Context, out interface{}) {
		db.First(out)
		outch <- out
	}(ctx, out)
	return outch
}
