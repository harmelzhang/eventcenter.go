package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var database *mongo.Database

// InitDB 初始化数据库连接
func InitDB(db *mongo.Database) {
	database = db
}

// 自定义操作

type Collection struct {
	ctx context.Context
	*mongo.Collection
}

type QuerySet struct {
	Collection *Collection
	Filter     bson.D
	Options    options.FindOptions
}

func (c *Collection) QuerySet() *QuerySet {
	return &QuerySet{Collection: c, Filter: make(bson.D, 0)}
}

func (qs *QuerySet) Q(key string, value any) *QuerySet {
	qs.Filter = append(qs.Filter, bson.E{Key: key, Value: value})
	return qs
}

func (qs *QuerySet) Sort(sort ...bson.E) *QuerySet {
	qs.Options.Sort = bson.D(sort)
	return qs
}

func (qs *QuerySet) Skip(i int64) *QuerySet {
	qs.Options.SetSkip(i)
	return qs
}

func (qs *QuerySet) Limit(i int64) *QuerySet {
	qs.Options.SetLimit(i)
	return qs
}

func (qs *QuerySet) All(results interface{}) error {
	cur, err := qs.Collection.Find(qs.Collection.ctx, qs.Filter, &qs.Options)
	if err != nil {
		return err
	}
	defer func() {
		err = cur.Close(qs.Collection.ctx)
		if err != nil {
			log.Printf("close cursor error: %v", err)
		}
	}()

	err = cur.All(qs.Collection.ctx, results)
	if err != nil {
		return err
	}
	return nil
}

func (qs *QuerySet) One(result interface{}) error {
	cur, err := qs.Collection.Find(qs.Collection.ctx, qs.Filter, &qs.Options)
	if err != nil {
		return err
	}
	defer func() {
		err = cur.Close(qs.Collection.ctx)
		if err != nil {
			log.Printf("close cursor error: %v", err)
		}
	}()

	if cur.Current == nil {
		if cur.TryNext(qs.Collection.ctx) {
			err = cur.Decode(result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DB(ctx context.Context, table string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{ctx, database.Collection(table, opts...)}
}
