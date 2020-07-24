package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

type config struct {
	addr string
	pwd  string
	db   int
}

func DefaultConfig() *config {
	return &config{}
}

func (c *config) Addr(addr string) *config {
	c.addr = addr
	return c
}

func (c *config) PWD(pwd string) *config {
	c.pwd = pwd
	return c
}

func (c *config) DB(db int) *config {
	c.db = db
	return c
}

func (c *config) Dial() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     c.addr,
		Password: c.pwd, // no password set
		DB:       c.db,  // use default DB
	})
	fmt.Println(rdb.Ping(context.TODO()))
}

func GetDB() *redis.Client {
	return rdb
}
