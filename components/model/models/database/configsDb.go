package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
)

/**
  查询label信息
*/
func QueryLabelsInfoById(labelId string) (interfaces.Target, error) {
  var data interfaces.Target
  db.UpdateTable("target")
  ids, err := primitive.ObjectIDFromHex(labelId)
  if err != nil {
    return data, err
  }
  selector := bson.M{"_id": ids}
  singleResult := db.QueryOne(selector)
  err = singleResult.Decode(&data)
  return data, err
}

/**
  插入modelConfigs信息
*/
func InsertModelConfigsInfo(data interfaces.ModelConfig) error {
  db.UpdateTable("modelConfig")
  return db.Insert(data)
}

/**
  更新modelConfigs信息
*/
func UpdateModelConfigsInfo(modelId string, config interfaces.Config) (*mongo.UpdateResult, error) {
  db.UpdateTable("modelConfig")
  selector := bson.M{"modelid": modelId}
  data := bson.M{"$push": bson.M{"configs": config}}
  return db.UpdateForResult(selector, data)
}

/**
  查询modelConfigs信息
*/
func QueryModelConfigsInfo(modelId string) (interfaces.ModelConfig, error) {
  var data interfaces.ModelConfig
  db.UpdateTable("modelConfig")
  selector := bson.M{"modelid": modelId}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return data, err
}
