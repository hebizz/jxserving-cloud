package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
)

/**
  查询project 是否存在
*/
func QueryProjectInfoByName(name string) (error, interfaces.Project) {
  var data interfaces.Project
  db.UpdateTable("project")
  selector := bson.M{"name": name}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return err, data
}

/**
  插入Project
*/
func InsertProjectInfo(data interfaces.Project) error {
  db.UpdateTable("project")
  return db.Insert(data)
}
