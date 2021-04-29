package database

import (
  "context"
  "os"
  "time"

  "github.com/spf13/cast"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  log "k8s.io/klog"
)

type MongoDB struct {
  Host      string
  Username  string
  Password  string
  Database  string
  TableName string
  client    *mongo.Client
}

func newMongoDB() *MongoDB {
  if ret := utils.ReadMap(ConfigEntryKey); ret != nil {
    value := cast.ToStringMapString(ret[TyMongo])
    log.Info("msg::", value["host"], value["database"], value["tablename"])
    return &MongoDB{
      Host:      value["host"],
      Database:  value["database"],
      TableName: value["tablename"],
      Username:  os.Getenv("JXS_ANA_MONGO_USR"),
      Password:  os.Getenv("JXS_ANA_MONGO_PWD"),
      client:    nil,
    }
  } else {
    panic("unable to find MongoDB config")
  }
}

func (db *MongoDB) CloseConnection() error {
  err := db.client.Disconnect(context.TODO())
  if err != nil {
    log.Error("disconnect mongo collection error:", err)
    return err
  }
  return nil
}

func (db *MongoDB) NewConnection() error {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.Host).SetMaxPoolSize(1024))
  if err != nil {
    return err
  }
  db.client = client
  return nil
}

func (db *MongoDB) CreateCollection() (context.Context, *mongo.Collection) {
  ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
  c := db.client.Database(db.Database).Collection(db.TableName)
  return ctx, c
}

func (db *MongoDB) UpdateDatabase(database string) {
  db.Database = database
}

func (db *MongoDB) UpdateTable(table string) {
  db.TableName = table
}

func (db *MongoDB) Insert(i interface{}) error {
  ctx, c := db.CreateCollection()
  _, err := c.InsertOne(ctx, i)
  return err
}

func (db *MongoDB) InsertMany(i []interface{}) error {
  ctx, c := db.CreateCollection()
  _, err := c.InsertMany(ctx, i)
  return err
}

func (db *MongoDB) Delete(i interface{}) error {
  ctx, c := db.CreateCollection()
  _, err := c.DeleteOne(ctx, i)
  return err
}

func (db *MongoDB) Update(filter interface{}, data interface{}) error {
  ctx, c := db.CreateCollection()
  _, err := c.UpdateOne(ctx, filter, data)
  return err
}

func (db *MongoDB) UpdateForResult(filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
  ctx, c := db.CreateCollection()
  result, err := c.UpdateOne(ctx, filter, data)
  return result, err
}

func (db *MongoDB) ReplaceOne(filter interface{}, data interface{}) error {
  ctx, c := db.CreateCollection()
  _, err := c.ReplaceOne(ctx, filter, data)
  return err
}

func (db *MongoDB) Commit() error {
  return nil
}

func (db *MongoDB) QueryCount(filter interface{}) (int64, error) {
  ctx, c := db.CreateCollection()
  count, err := c.CountDocuments(ctx, filter)
  return count, err
}

func (db *MongoDB) QueryOne(filter interface{}) *mongo.SingleResult {
  ctx, c := db.CreateCollection()
  singleResult := c.FindOne(ctx, filter)
  return singleResult
}

func (db *MongoDB) QueryAll(filter interface{}, options *options.FindOptions) (*mongo.Cursor, context.Context, error) {
  ctx, c := db.CreateCollection()
  cur, err := c.Find(ctx, filter, options)
  return cur, ctx, err
}

func (db *MongoDB) Query(filter interface{}) (interface{}, error) {
  ctx, c := db.CreateCollection()
  cur, err := c.Find(ctx, filter)
  if err != nil {
    return "", err
  }
  defer cur.Close(ctx)

  var ret []bson.M
  for cur.Next(ctx) {
    var result bson.M
    err := cur.Decode(&result)
    if err != nil {
      return "", err
    }
    ret = append(ret, result)
  }
  if err := cur.Err(); err != nil {
    return "", err
  }
  return ret, nil
}

func (db *MongoDB) RemoveOne(i interface{}) error {
  ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
  c := db.client.Database(db.Database).Collection(db.TableName)
  _, err := c.DeleteOne(ctx, i)
  if err != nil {
    return err
  }
  return nil
}

func (db *MongoDB) RemoveMany(i interface{}) error {
  ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
  c := db.client.Database(db.Database).Collection(db.TableName)
  _, err := c.DeleteMany(ctx, i)
  if err != nil {
    return err
  }
  return nil
}

func (db *MongoDB) UpdateOne(i interface{}, condition interface{}) error {
  ctx, c := db.CreateCollection()
  _, err := c.UpdateOne(ctx, i, condition)
  if err != nil {
    return err
  }
  return nil
}
