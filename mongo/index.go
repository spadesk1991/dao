package mongo

import (
	"reflect"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

var s *mgo.Session
var cfg *config

const mongodb = "mws"

type config struct {
	url     string
	db      string
	mode    mgo.Mode
	refresh bool
}

func DefaultConfig() *config {
	cfg = &config{
		mode:    mgo.Monotonic,
		refresh: true,
	}
	return cfg
}

func (c *config) Url(str string) *config {
	c.url = str
	return c
}

func (c *config) DB(db string) *config {
	c.db = db
	return c
}

func (c *config) Mode(m mgo.Mode) *config {
	c.mode = m
	return c
}

func (c *config) Refresh(b bool) *config {
	c.refresh = b
	return c
}

func (c *config) Dial() {
	var err error
	s, err = mgo.Dial(c.url)
	if err != nil {
		panic(err)
	}
	s.SetMode(c.mode, c.refresh)
}

type myC struct {
	*mgo.Session
	*mgo.Collection
}

func getModelName(m interface{}) string {
	v := reflect.ValueOf(m)
	ss := strings.Split(v.String(), ".")
	return ss[len(ss)-1]
}
func GetDB() (*mgo.Session, *mgo.Database) {
	ms := s.Copy()
	return ms, ms.DB(mongodb)
}

func Model(m interface{}) *myC {
	name := getModelName(m)
	ms := s.Copy()
	c := ms.DB(cfg.db).C(name)
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
	return q.Explain(result)
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
	return q.Query.One(&result)
}

func (q *myQuery) All(result interface{}) error {
	defer q.Session.Close()
	return q.Query.All(result)
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

//func Count(collection string, query interface{}) (int, error) {
//	ms, c := connect(collection)
//	defer ms.Close()
//	return c.Find(query).Count()
//}
//
//func Insert(collection string, docs ...interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//	return c.Insert(docs...)
//}
//
//func FindOne(collection string, query, selector, result interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	return c.Find(query).Select(selector).One(result)
//}
//
//func FindAll(collection string, query, selector, result interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	return c.Find(query).Select(selector).All(result)
//}
//
//func FindPage(collection string, page, limit int, query, selector, result interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	return c.Find(query).Select(selector).Skip((page - 1) * limit).Limit(limit).All(result)
//}
//
//func FindIter(collection string, query interface{}) *mgo.Iter {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	return c.Find(query).Iter()
//}
//
//func Update(collection string, selector, update interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	return c.Update(selector, update)
//}
//
//func Upsert(collection string, selector, update interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	_, err := c.Upsert(selector, update)
//	return err
//}
//
//func UpdateAll(collection string, selector, update interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	_, err := c.UpdateAll(selector, update)
//	return err
//}
//
//func Remove(collection string, selector interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	return c.Remove(selector)
//}
//
//func RemoveAll(collection string, selector interface{}) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	_, err := c.RemoveAll(selector)
//	return err
//}
//
////insert one or multi documents
//func BulkInsert(collection string, docs ...interface{}) (*mgo.BulkResult, error) {
//	ms, c := connect(collection)
//	defer ms.Close()
//	bulk := c.Bulk()
//	bulk.Insert(docs...)
//	return bulk.Run()
//}
//
//func BulkRemove(collection string, selector ...interface{}) (*mgo.BulkResult, error) {
//	ms, c := connect(collection)
//	defer ms.Close()
//
//	bulk := c.Bulk()
//	bulk.Remove(selector...)
//	return bulk.Run()
//}
//
//func BulkRemoveAll(collection string, selector ...interface{}) (*mgo.BulkResult, error) {
//	ms, c := connect(collection)
//	defer ms.Close()
//	bulk := c.Bulk()
//	bulk.RemoveAll(selector...)
//	return bulk.Run()
//}
//
//func BulkUpdate(collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
//	ms, c := connect(collection)
//	defer ms.Close()
//	bulk := c.Bulk()
//	bulk.Update(pairs...)
//	return bulk.Run()
//}
//
//func BulkUpdateAll(collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
//	ms, c := connect(collection)
//	defer ms.Close()
//	bulk := c.Bulk()
//	bulk.UpdateAll(pairs...)
//	return bulk.Run()
//}
//
//func BulkUpsert(collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
//	ms, c := connect(collection)
//	defer ms.Close()
//	bulk := c.Bulk()
//	bulk.Upsert(pairs...)
//	return bulk.Run()
//}
//
//func PipeAll(collection string, pipeline, result interface{}, allowDiskUse bool) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//	var pipe *mgo.Pipe
//	if allowDiskUse {
//		pipe = c.Pipe(pipeline).AllowDiskUse()
//	} else {
//		pipe = c.Pipe(pipeline)
//	}
//	return pipe.All(result)
//}
//
//func PipeOne(collection string, pipeline, result interface{}, allowDiskUse bool) error {
//	ms, c := connect(collection)
//	defer ms.Close()
//	var pipe *mgo.Pipe
//	if allowDiskUse {
//		pipe = c.Pipe(pipeline).AllowDiskUse()
//	} else {
//		pipe = c.Pipe(pipeline)
//	}
//	return pipe.One(result)
//}
