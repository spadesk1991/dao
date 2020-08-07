package mysql

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type options struct {
	url          string
	debug        bool
	maxIdleConns int
	maxOpenConns int
	logger       logger
}

type Option func(*options)

func NewMysqlOptions(opts ...Option) *options {
	return newOptions(opts...)
}

func newOptions(opts ...Option) *options {
	o := &options{
		maxIdleConns: -1,
		maxOpenConns: -1,
		logger:       log.New(os.Stdout, "\r\n", log.LstdFlags),
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type logger interface {
	Print(v ...interface{})
}

func Url(str string) Option {
	return func(o *options) {
		o.url = str
	}
}

func Debug() Option {
	return func(o *options) {
		o.debug = true
	}
}

func MaxIdleConns(n int) Option {
	return func(o *options) {
		o.maxIdleConns = n
	}
}

func MaxOpenConns(n int) Option {
	return func(o *options) {
		o.maxOpenConns = n
	}
}

func SetLog(l logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

func (o *options) Dial() *gorm.DB {
	newdb, err := gorm.Open("mysql", o.url)
	if err != nil {
		panic(err)
	}
	if o.debug { // 生产环境关闭log
		newdb.LogMode(true)
	}
	newdb.SetLogger(o.logger)
	newdb.DB().SetMaxIdleConns(o.maxIdleConns)
	newdb.DB().SetMaxOpenConns(o.maxOpenConns)
	db = newdb // 赋值全局db对象
	return newdb
}

func GetDB() *gorm.DB {
	return db
}
