package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

type options struct {
	addr string
	pwd  string
	db   int
}

type Option func(*options)

func newOptions(opts ...Option) options {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Addr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

func PWD(pwd string) Option {
	return func(o *options) {
		o.pwd = pwd
	}
}

func DB(db int) Option {
	return func(o *options) {
		o.db = db
	}
}

func NewRedisOptions(opts ...Option) options {
	return newOptions(opts...)
}

func (o options) Dial() *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:     o.addr,
		Password: o.pwd, // no password set
		DB:       o.db,  // use default DB
	})
	fmt.Println(db.Ping(context.TODO()))
	rdb = db // 全局db
	return db
}

func GetDB() *redis.Client {
	return rdb
}
