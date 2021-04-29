package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
)

/**
  插入 performance 数据
*/
func InsertPerformanceData(data interfaces.ModelPerformance) error {
  db.UpdateTable("performance")
  return db.Insert(data)
}

/**
  delete models from performance
*/
func DeleteFromPerformanceData(modelMd5 string) error {
  selector := bson.M{"modelmd5": modelMd5}
  db.UpdateTable("performance")
  return db.RemoveMany(selector)
}

/**
  query performance data
*/
func QueryPerformanceData(modelMd5 string) ([]interfaces.ModelPerformance, error) {
  var models []interfaces.ModelPerformance
  selector := bson.M{"modelmd5": modelMd5}
  db.UpdateTable("performance")
  cursor, ctx, err := db.QueryAll(selector, nil)
  if err != nil {
    return models, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &models)
  return models, err
}
