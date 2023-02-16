package base

import (
	"context"
	"github.com/aluka-7/common"
	"github.com/aluka-7/gomongo/search"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IBaseRepository interface {
	Save(ctx context.Context, collection string, bean interface{}) (interface{}, error)
	Update(ctx context.Context, collection string, id, update interface{}) (int64, error)
	ReadById(ctx context.Context, collection string, id, bean interface{}) error
	Query(ctx context.Context, cq common.Query, collection string, list interface{}) (page *common.Pagination, err error)
}

func NewBaseRepository(gmc *mongo.Database, column map[string]search.Filter) BaseRepository {
	return BaseRepository{gmc: gmc, column: column}
}

type BaseRepository struct {
	gmc    *mongo.Database
	column map[string]search.Filter
}

func (b *BaseRepository) Save(ctx context.Context, collection string, bean interface{}) (interface{}, error) {
	res, err := b.gmc.Collection(collection).InsertOne(ctx, bean)
	return res.InsertedID, err
}

func (b BaseRepository) Update(ctx context.Context, collection string, id, update interface{}) (int64, error) {
	res, err := b.gmc.Collection(collection).UpdateByID(ctx, id, update)
	return res.ModifiedCount, err
}

func (b BaseRepository) ReadById(ctx context.Context, collection string, id, bean interface{}) error {
	return b.gmc.Collection(collection).FindOne(ctx, bson.D{{"_id", id}}).Decode(bean)
}

func (b BaseRepository) Query(ctx context.Context, cq common.Query, collection string, list interface{}) (page *common.Pagination, err error) {
	query := search.NewQuery(cq)
	filter := query.MarkFiltered(b.column)
	option := query.MarkSort(b.column)
	page = query.MarkPage()
	limit, skip := page.Limit()
	option.SetLimit(int64(limit))
	option.SetSkip(int64(skip))
	cursor, err := b.gmc.Collection(collection).Find(ctx, filter, option)
	defer cursor.Close(ctx)
	if err != nil {
		return
	}
	if err = cursor.All(ctx, list); err != nil {
		return
	}
	var total int64
	if total, err = b.gmc.Collection(collection).CountDocuments(ctx, filter); err == nil {
		page.SetTotalRecord(int(total))
	}
	return
}
