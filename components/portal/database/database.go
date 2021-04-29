package database

import (
  "context"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

const (
  ConfigEntryKey = "database"
  TyInflux       = "influx"
  TyMongo        = "mongo"
)

type Database interface {
  NewConnection() error
  UpdateDatabase(string)
  UpdateTable(string)
  Insert(interface{}) error
  Commit() error
  Query(interface{}) (interface{}, error)
}
type MongoDatabase interface {
  Database
  InsertMany([]interface{}) error
  Update(interface{}, interface{}) error
  UpdateForResult(interface{}, interface{}) (*mongo.UpdateResult, error)
  ReplaceOne(interface{}, interface{}) error
  Delete(interface{}) error
  QueryCount(interface{}) (int64, error)
  QueryOne(interface{}) *mongo.SingleResult
  QueryAll(interface{}, *options.FindOptions) (*mongo.Cursor, context.Context, error)
  RemoveOne(i interface{}) error
  RemoveMany(i interface{}) error
  UpdateOne(i interface{}, condition interface{}) error
}

func NewDatabase(t string) (Database, error) {
  switch t {
  case TyInflux:
    return newInfluxDB(), nil
  case TyMongo:
    return newMongoDB(), nil
  default:
    panic("unsupported database")
  }
}
