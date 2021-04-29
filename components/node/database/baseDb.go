package database

import "gitlab.jiangxingai.com/jxserving/components/portal/database"

var db database.MongoDatabase

func init() {
  rdb, _ := database.NewDatabase(database.TyMongo)
  dbClient := rdb.(database.MongoDatabase)
  err := dbClient.NewConnection()
  if err != nil {
    panic(err)
  }
  db = dbClient
}
