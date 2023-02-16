package gomongo_test

import (
	"context"
	"github.com/aluka-7/common"
	"github.com/aluka-7/configuration"
	"github.com/aluka-7/configuration/backends"
	"github.com/aluka-7/gomongo"
	"github.com/aluka-7/gomongo/base"
	"github.com/aluka-7/gomongo/search"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestGomongo(t *testing.T) {
	conf := configuration.MockEngine(t, backends.StoreConfig{Exp: map[string]string{
		"/system/base/mongo/1000":   "{\"uri\":\"mongodb://localhost:27017/\",\"database\":\"ab_1000\"}",
		"/system/base/mongo/common": "{\"timeOut\":2,\"maxPoolSize\":100,\"minPoolSize\":10,\"maxConnecting\":50,\"maxConnIdleTime\":300}",
	}})

	db := gomongo.Engine(conf, "1000").Connection("")
	br := base.NewBaseRepository(db, map[string]search.Filter{
		// 检索
		"name": {FieldName: "name", Operator: search.EQ},
		"age":  {FieldName: "age", Operator: search.LTE},

		// 排序
		"ageSort": {FieldName: "age"},
	})
	ctx := context.Background()

	// insert
	insertedID, err := br.Save(ctx, "test", struct {
		Name  string `bson:"name"`
		Age   uint   `bson:"age"`
		Email string `bson:"email"`
	}{
		"uio", 13, "efv@ijn.com",
	})
	if err != nil {
		t.Errorf("Save Error: %e", err)
	}
	t.Logf("InsertId: %s", insertedID.(primitive.ObjectID).Hex())

	// update
	affectedRows, err := br.Update(ctx, "test", insertedID, bson.M{"$set": struct {
		Name  string `bson:"name"`
		Age   uint   `bson:"age"`
		Email string `bson:"email"`
	}{
		"123456", 13, "SDF@RFV.com",
	}})
	if err != nil {
		t.Errorf("Update Error: %e", err)
	}
	t.Logf("Affected Rows: %d", affectedRows)

	// get one
	var bean struct {
		Id    primitive.ObjectID `bson:"_id"`
		Name  string             `bson:"name"`
		Age   uint               `bson:"age"`
		Email string             `bson:"email"`
	}
	err = br.ReadById(ctx, "test", insertedID, &bean)
	if err != nil {
		t.Errorf("Update Error: %e", err)
	}
	t.Logf("Read Row: %v", bean)

	// query
	var (
		cq = common.Query{
			PageSize: 5,
			Page:     1,
		}
		list []struct {
			Id    primitive.ObjectID `bson:"_id"`
			Name  string             `bson:"name"`
			Age   uint               `bson:"age"`
			Email string             `bson:"email"`
		}
	)
	cq.Sorted = append(cq.Sorted, struct {
		Id   string `json:"id"`
		Desc bool   `json:"desc"`
	}{Id: "ageSort", Desc: true})
	cq.Filtered = append(cq.Filtered, struct {
		Id    string      `json:"id"`
		Value interface{} `json:"value"`
	}{Id: "age", Value: 10})
	page, err := br.Query(ctx, cq, "test", &list)
	if err != nil {
		t.Errorf("Query Error: %e", err)
	}
	t.Logf("List: %v, Page: %v", list, page)

}
