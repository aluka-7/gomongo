package search

import (
	"github.com/aluka-7/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// {"pageSize":10,"page":0,"sorted":[{"id":"firstName","desc":false}],"filtered":[{"id":"firstName","value":"3"}]}
type Query struct {
	common.Query
}

func NewQuery(query common.Query) Query {
	return Query{query}
}
func (sp Query) MarkSort(column map[string]Filter) (option *options.FindOptions) {
	option = options.Find()
	if len(sp.Sorted) > 0 {
		for _, v := range sp.Sorted {
			if v.Desc {
				option.SetSort(bson.M{column[v.Id].FieldName: -1})
			} else {
				option.SetSort(bson.M{column[v.Id].FieldName: 1})
			}
		}
	}
	return
}

func (sp Query) MarkFiltered(column map[string]Filter) bson.D {
	var filter = bson.D{}
	if len(sp.Filtered) > 0 {
		for _, v := range sp.Filtered {
			if k, ok := column[v.Id]; ok {
				switch k.Operator {
				case NE:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$ne": v.Value}})
				case LIKE:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$regex": v.Value}})
				case GT:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$gt": v.Value}})
				case LT:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$lt": v.Value}})
				case GTE:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$gte": v.Value}})
				case LTE:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$lte": v.Value}})
				case IN:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$in": v.Value}})
				case NI:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$nin": v.Value}})
				default:
					filter = append(filter, bson.E{k.FieldName, bson.M{"$eq": v.Value}})
				}
			}
		}
	}
	return filter
}
