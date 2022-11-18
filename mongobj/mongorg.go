// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

/*
Package mongobj implements a simple interface to access mongodb.
將自定義的資料結構藉由Primary Key存入與取出mongodb，而不使用內建的ObjectID(_id)。

	TODO: 將此package移動到獨立的repo
*/
package mongobj

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"reflect"
	"strings"
)

// Mongorg means mongo object organizer
type Mongorg struct {
	client *mongo.Client
}

func New(dbUri string) (*Mongorg, error) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		return nil, err
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return &Mongorg{
		client: client,
	}, nil
}

func (m *Mongorg) Destroy() error {
	if m.client != nil {
		return m.client.Disconnect(context.TODO())
	}
	return nil
}

func (m *Mongorg) Collection(dbName string, collName string) *Mongolection {
	coll := m.client.Database(dbName).Collection(collName)

	mongolection := &Mongolection{
		mongorg: m,
		coll:    coll,
	}

	return mongolection
}

// Mongolection means mongo object collection
type Mongolection struct {
	mongorg *Mongorg
	coll    *mongo.Collection
	//primaryKeys []string
}

/*
InitializeUniqueKey
用來初始化Primary Key，多次設定相同的schema不會造成問題，這個package目前基於pk來存取物件。

	TODO: 怎麼有效讓程式碼只設定一次。因為檢查與設定大概都會消耗存取次數，基本上應該沒什麼太大差別？
*/
func (m *Mongolection) InitializeUniqueKey(keys ...string) error {
	idxView := m.coll.Indexes()
	//opts := options.ListIndexes() //.SetMaxTime(2 * time.Second)
	//cursor, err := idxView.List(context.TODO(), opts)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// Get a slice of all indexes returned and print them out.
	//var results []bson.M
	//if err = cursor.All(context.TODO(), &results); err != nil {
	//	log.Fatal(err)
	//}
	//for _, r := range results {
	//	fmt.Println(r)
	//}

	keysDoc := bsonx.Doc{}
	var indexName string
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			// desc
			key = strings.TrimLeft(key, "-")
			keysDoc = keysDoc.Append(key, bsonx.Int32(-1))
		} else {
			// asc
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
		indexName += "-"
		indexName += key
	}

	realName, err := idxView.CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(true),
		},
		nil, //options.CreateIndexes().SetMaxTime(10*time.Second),
	)

	fmt.Printf("Index Name: %s\n", realName)

	return err
}

/*
Get will find values for pk.
If there are no values for pk, nil without error will be returned.
藉由primary key取回自定義資料結構，取回為自定義資料結構的指標。

	Ex.
		_ = m.Get(bson.M{"PK": "1234567"}, reflect.TypeOf(MyStructure{})).(*MyStructure)
*/
func (m *Mongolection) Get(pk bson.M, resultType reflect.Type) (interface{}, error) {
	result := reflect.New(resultType)
	//initializeStruct(resultType, result.Elem())
	resultInterface := result.Interface()
	//fmt.Println("resultType:", resultType)
	//fmt.Println("result:", result)
	//fmt.Println("result.Elem():", result.Elem())
	//fmt.Println("reflect.TypeOf(result):", reflect.TypeOf(resultInterface))
	//tmp := bson.M{}
	err := m.coll.FindOne(context.TODO(), pk).Decode(resultInterface)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, nil
	}
	//bson.Unmarshal()
	fmt.Println(resultInterface)
	return resultInterface, nil
}

/*
Set will upsert values for pk, values should include pk.

	Ex.
		type MyStructure struct {
			ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
			MyPk string `json:"myPk" binding:"required" bson:"myPk"`
			MyName string `json:"myName" binding:"required" bson:"myName"`
		}

		values := MyStructure{MyPk:"1234567", MyName:"Honey"}
		_ = m.Set(bson.M{"myPk": "1234567", values}, values)
*/
func (m *Mongolection) Set(pk bson.M, values interface{}) error {
	opts := options.Replace().SetUpsert(true)
	_, err := m.coll.ReplaceOne(context.TODO(), pk, values, opts)
	return err
}

func (m Mongolection) SetBulk(pks []bson.M, bulkValues []interface{}) error {
	if len(pks) != len(bulkValues) || len(pks) == 0 {
		return errors.New("number of pks or values is incorrect")
	}
	models := make([]mongo.WriteModel, 0, len(pks))
	for i, pk := range pks {
		values := bulkValues[i]
		models = append(models, mongo.NewReplaceOneModel().SetFilter(pk).SetUpsert(true).SetReplacement(values))
	}
	opts := options.BulkWrite().SetOrdered(true)
	_, err := m.coll.BulkWrite(context.TODO(), models, opts)
	return err
}

//func initializeStruct(t reflect.Type, v reflect.Value) {
//	for i := 0; i < v.NumField(); i++ {
//		f := v.Field(i)
//		ft := t.Field(i)
//		switch ft.Type.Kind() {
//		case reflect.Map:
//			f.Set(reflect.MakeMap(ft.Type))
//		case reflect.Slice:
//			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
//		case reflect.Chan:
//			f.Set(reflect.MakeChan(ft.Type, 0))
//		case reflect.Struct:
//			initializeStruct(ft.Type, f)
//		case reflect.Ptr:
//			fv := reflect.New(ft.Type.Elem())
//			initializeStruct(ft.Type.Elem(), fv.Elem())
//			f.Set(fv)
//		default:
//		}
//	}
//}
