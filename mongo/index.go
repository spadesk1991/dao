package mongo

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/inflection"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

var s *mgo.Session
var o *options

type options struct {
	url     string
	db      string
	mode    mgo.Mode
	refresh bool
}

type Option func(*options)

func newOptions(opts ...Option) *options {
	o = &options{
		url:     "mongodb://localhost:27017/admin",
		db:      "admin",
		mode:    mgo.Monotonic,
		refresh: true,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func NewMongoOptions(opts ...Option) *options {
	return newOptions(opts...)
}

func Url(str string) Option {
	return func(o2 *options) {
		o2.url = str
	}
}

func DB(db string) Option {
	return func(o2 *options) {
		o2.db = db
	}
}

func Mode(m mgo.Mode) Option {
	return func(o2 *options) {
		o2.mode = m
	}
}

func Refresh(b bool) Option {
	return func(o2 *options) {
		o2.refresh = b
	}
}

func (o *options) Dial() {
	var err error
	s, err = mgo.Dial(o.url)
	if err != nil {
		panic(err)
	}
	s.SetMode(o.mode, o.refresh)
}

type myC struct {
	*mgo.Session
	*mgo.Collection
}

func getModelName(m interface{}) string {
	t := reflect.TypeOf(m)
	return inflection.Plural(defaultName(t.Name()))
}
func GetDB() (*mgo.Session, *mgo.Database) {
	ms := s.Copy()
	return ms, ms.DB(o.db)
}

func Do(fn func(db *mgo.Database) error) (err error) {
	ms := s.Copy()
	defer ms.Close()

	db := ms.DB(o.db)
	return fn(db)
}

func Model(m interface{}) *myC {
	name := getModelName(m)
	ms := s.Copy()
	c := ms.DB(o.db).C(name)
	return &myC{
		Session:    ms,
		Collection: c,
	}
}

func (c *myC) EnsureIndexKey(key ...string) error {
	defer c.Session.Close()
	return c.Collection.EnsureIndexKey(key...)
}

func (c *myC) EnsureIndex(index mgo.Index) error {
	defer c.Session.Close()
	return c.Collection.EnsureIndex(index)
}

func (c *myC) DropIndex(key ...string) error {
	defer c.Session.Close()
	return c.Collection.DropIndex(key...)
}

func (c *myC) DropIndexName(name string) error {
	defer c.Session.Close()
	return c.Collection.DropIndexName(name)
}

func (c *myC) Indexes() (indexes []mgo.Index, err error) {
	defer c.Session.Close()
	return c.Collection.Indexes()
}

type myQuery struct {
	*mgo.Query
	*mgo.Session
}

func (c *myC) Find(query interface{}) *myQuery {
	q := c.Collection.Find(query)
	mq := &myQuery{q, c.Session}
	return mq
}

func (c *myC) FindId(id interface{}) *myQuery {
	q := c.Collection.FindId(id)
	mq := &myQuery{q, c.Session}
	return mq
}

type myPipe struct {
	*mgo.Pipe
}

func (c *myC) Pipe(pipeline interface{}) *myPipe {
	p := c.Collection.Pipe(pipeline)
	mp := &myPipe{p}
	return mp
}

func (c *myC) Insert(docs ...interface{}) error {
	defer c.Session.Close()
	return c.Collection.Insert(docs...)
}

func (c *myC) Update(selector interface{}, update interface{}) error {
	defer c.Session.Close()
	return c.Collection.Update(selector, update)
}

func (c *myC) UpdateId(id interface{}, update interface{}) error {
	defer c.Session.Close()
	return c.Collection.UpdateId(id, update)
}

func (c *myC) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	defer c.Session.Close()
	return c.Collection.UpdateAll(selector, update)
}

func (c *myC) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	defer c.Session.Close()
	return c.Collection.Upsert(selector, update)
}

func (c *myC) UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	defer c.Session.Close()
	return c.Collection.UpsertId(id, update)
}

func (c *myC) Remove(selector interface{}) error {
	defer c.Session.Close()
	return c.Collection.Remove(selector)
}

func (c *myC) RemoveId(id interface{}) error {
	defer c.Session.Close()
	return c.Collection.RemoveId(id)
}

func (c *myC) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	defer c.Session.Close()
	return c.Collection.RemoveAll(selector)
}

func (c *myC) DropCollection() error {
	defer c.Session.Close()
	return c.Collection.DropCollection()
}

func (c *myC) Create(info *mgo.CollectionInfo) error {
	defer c.Session.Close()
	return c.Collection.Create(info)
}

func (q *myQuery) Batch(n int) *myQuery {
	q.Query = q.Query.Batch(n)
	return q
}

func (q *myQuery) Prefetch(p float64) *myQuery {
	q.Query = q.Query.Prefetch(p)
	return q
}

func (q *myQuery) Skip(n int) *myQuery {
	q.Query = q.Query.Skip(n)
	return q
}

func (q *myQuery) Limit(n int) *myQuery {
	q.Query = q.Query.Limit(n)
	return q
}

func (q *myQuery) Select(selector interface{}) *myQuery {
	q.Query = q.Query.Select(selector)
	return q
}

func (q *myQuery) Sort(fields ...string) *myQuery {
	q.Query = q.Query.Sort(fields...)
	return q
}

func (q *myQuery) Explain(result interface{}) error {
	defer q.Session.Close()
	return q.Query.Explain(&result)
}

func (q *myQuery) Hint(indexKey ...string) *myQuery {
	q.Query = q.Query.Hint(indexKey...)
	return q
}

func (q *myQuery) SetMaxScan(n int) *myQuery {
	q.Query = q.Query.SetMaxScan(n)
	return q
}

func (q *myQuery) SetMaxTime(d time.Duration) *myQuery {
	q.Query = q.Query.SetMaxTime(d)
	return q
}

func (q *myQuery) Snapshot() *myQuery {
	q.Query = q.Query.Snapshot()
	return q
}

func (q *myQuery) Comment(comment string) *myQuery {
	q.Query = q.Query.Comment(comment)
	return q
}

func (q *myQuery) LogReplay() *myQuery {
	q.Query = q.Query.LogReplay()
	return q
}

func (q *myQuery) One(result interface{}) (err error) {
	defer q.Session.Close()
	data := bson.M{}
	if err = q.Query.One(&data); err != nil {
		return
	}
	bt, err := json.Marshal(&data)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bt, result); err != nil {
		return
	}
	return
}

func (q *myQuery) All(result interface{}) (err error) {
	defer q.Session.Close()
	data := make([]bson.M, 0)
	if err = q.Query.All(&data); err != nil {
		return
	}
	bt, err := json.Marshal(&data)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bt, result); err != nil {
		return
	}
	return
}

func (q *myQuery) Count() (n int, err error) {
	defer q.Session.Close()
	return q.Query.Count()
}

func (c *myC) Count() (n int, err error) {
	defer c.Session.Close()
	return c.Collection.Count()
}
func (q *myQuery) Distinct(key string, result interface{}) error {
	defer q.Session.Close()
	return q.Query.Distinct(key, result)
}

func (q *myQuery) MapReduce(job *mgo.MapReduce, result interface{}) (info *mgo.MapReduceInfo, err error) {
	defer q.Session.Close()
	return q.Query.MapReduce(job, result)
}

func (q *myQuery) Apply(change mgo.Change, result interface{}) (info *mgo.ChangeInfo, err error) {
	defer q.Session.Close()
	return q.Query.Apply(change, result)
}

func defaultName(name string) string {
	const (
		lower = false
		upper = true
	)

	if name == "" {
		return ""
	}

	var (
		value                                    = name
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	return s
}
