package data

import (
	"context"
	"fmt"
	"github.com/leondevpt/wallet/trxservice/pkg/setting"

	"database/sql"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedis, NewDB, NewTrxRepo)

type contextTxKey struct{}

// Data .
type Data struct {
	// TODO wrapped database client
	db  *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(db *gorm.DB, rdb *redis.Client) (*Data, error) {
	return &Data{db: db, rdb: rdb}, nil
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func NewDB(c *setting.Config) *gorm.DB {
	newLogger := zapgorm2.New(zap.L())
	newLogger.SetAsDefault()
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.DbName)
	_, err := tryOpenOrCreateDB(c.DB.Driver, source, c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.DbName)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名是否加 s
		},
	})

	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	//InitDB(db)
	return db
}

/*
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&biz.Tx{}); err != nil {
		panic(err)
	}
}
*/

func NewRedis(c *setting.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       int(c.Redis.DB),
	})
	rdb.AddHook(redisotel.TracingHook{})
	if err := rdb.Close(); err != nil {
		fmt.Printf("redis close err:%s\n", err.Error())
	}
	return rdb
}

func tryOpenOrCreateDB(driver string, url, user, passwd string, host string, port int, dbName string) (*sql.DB, error) {
	if user == "" {
		user = "root"
	}
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, passwd, host, port)
	db, err := sql.Open(driver, source)
	if err != nil {
		zap.S().Error(driver, "Open database error", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		zap.S().Error(driver, "ping error", err.Error())
		return nil, err
	}
	checkOrCreateCmd := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", dbName)
	zap.S().Info("OpenOrCreateDB", "Cmd", checkOrCreateCmd)
	_, err = db.Exec(checkOrCreateCmd)
	if err != nil {
		zap.S().Error("OpenOrCreateDB", "create database error", err.Error())
		return nil, err
	}

	if db, err = sql.Open(driver, url); err != nil {
		zap.S().Error(driver, "retry open database error", err.Error())
		return nil, err
	}
	if err = db.Ping(); err != nil {
		zap.S().Error(driver, "retry ping error", err.Error())
		return nil, err
	}
	return db, nil
}
