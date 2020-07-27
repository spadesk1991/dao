package mongo

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2"
)

type Demo struct {
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}

func TestRun(t *testing.T) {
	NewMongoOptions(Url("mongodb://localhost:27017/test"), DB("test"), Mode(mgo.Monotonic), Refresh(true)).Dial()

	if err := Model(Demo{}).Insert(&Demo{
		Name: "jack1",
		Age:  12,
	}); err != nil {
		t.Fatal(err)
	}
	res := make([]Demo, 0)
	if err := Model(Demo{}).Find(nil).All(&res); err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	t.Log("ok")
}
