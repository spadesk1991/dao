package mysql

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type config struct {
	url          string
	debug        bool
	maxIdleConns int
	maxOpenConns int
	logger       logger
}

type logger interface {
	Print(v ...interface{})
}

func DefaultConfig() *config {
	return &config{
		maxIdleConns: -1,
		maxOpenConns: -1,
		logger:       log.New(os.Stdout, "\r\n", log.LstdFlags),
	}
}

func (c *config) Url(str string) *config {
	c.url = str
	return c
}

func (c *config) Debug() *config {
	c.debug = true
	return c
}

func (c *config) MaxIdleConns(n int) *config {
	c.maxIdleConns = n
	return c
}

func (c *config) MaxOpenConns(n int) *config {
	c.maxOpenConns = n
	return c
}

func (c *config) SetLogOut(l logger) *config {
	c.logger = l
	return c
}

func (c *config) Dial() {
	var err error
	db, err = gorm.Open("mysql", c.url)
	if err != nil {
		panic(err)
	}
	if c.debug { // 生产环境关闭log
		db.LogMode(true)
	}
	db.SetLogger(c.logger)
	db.DB().SetMaxIdleConns(c.maxIdleConns)
	db.DB().SetMaxOpenConns(c.maxOpenConns)
}

func GetDB() *gorm.DB {

	return db
}
