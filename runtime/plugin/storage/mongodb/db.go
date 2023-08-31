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
	collection *Collection
	filter     bson.D
	options    options.FindOptions
}

func (c *Collection) QuerySet() *QuerySet {
	return &QuerySet{collection: c, filter: make(bson.D, 0)}
}

func (qs *QuerySet) Q(key string, value any) *QuerySet {
	qs.filter = append(qs.filter, bson.E{Key: key, Value: value})
	return qs
}

func (qs *QuerySet) Filter(filter bson.D) *QuerySet {
	qs.filter = append(qs.filter, filter...)
	return qs
}

func (qs *QuerySet) Sort(sort ...bson.E) *QuerySet {
	qs.options.Sort = bson.D(sort)
	return qs
}

func (qs *QuerySet) Skip(i int64) *QuerySet {
	qs.options.SetSkip(i)
	return qs
}

func (qs *QuerySet) Limit(i int64) *QuerySet {
	qs.options.SetLimit(i)
	return qs
}

func (qs *QuerySet) Count() (cnt int64, err error) {
	return qs.collection.CountDocuments(qs.collection.ctx, qs.filter, &options.CountOptions{
		Limit: qs.options.Limit,
		Skip:  qs.options.Skip,
	})
}

func (qs *QuerySet) All(results interface{}) error {
	cur, err := qs.collection.Find(qs.collection.ctx, qs.filter, &qs.options)
	if err != nil {
		return err
	}
	defer func() {
		err = cur.Close(qs.collection.ctx)
		if err != nil {
			log.Printf("close cursor error: %v", err)
		}
	}()

	err = cur.All(qs.collection.ctx, results)
	if err != nil {
		return err
	}
	return nil
}

func (qs *QuerySet) One(result interface{}) error {
	cur, err := qs.collection.Find(qs.collection.ctx, qs.filter, &qs.options)
	if err != nil {
		return err
	}
	defer func() {
		err = cur.Close(qs.collection.ctx)
		if err != nil {
			log.Printf("close cursor error: %v", err)
		}
	}()

	if cur.Current == nil {
		if cur.TryNext(qs.collection.ctx) {
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
