package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

/**
  插入 model 数据
*/
func InsertModelData(data interfaces.Model) (string, error) {
  id := primitive.NewObjectID()
  data.Id = id
  db.UpdateTable("model")
  return id.Hex(), db.Insert(data)
}

/**
  更新 model 数据
*/
func UpdateLabelSetDataByMd5(modelMd5 string) error {
  selector := bson.M{"modelmd5": modelMd5}
  data := bson.D{{"$set", bson.D{{"ispublished", true},}},}
  db.UpdateTable("model")
  return db.Update(selector, data)
}

/**
  delete model
*/
func DeleteModelInfo(modelMd5 string) error {
  selector := bson.M{"modelmd5": modelMd5}
  db.UpdateTable("model")
  return db.Delete(selector)
}

/**
  查询 [] model
*/
func QueryModelsData() ([]interfaces.Model, error) {
  var models []interfaces.Model
  db.UpdateTable("model")
  cursor, ctx, err := db.QueryAll(bson.M{}, nil)
  if err != nil {
    return models, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &models)
  return models, err
}

/**
  查询模型是否存在
*/
func QueryModelIsExistByMd5(modelMd5 string) (int64, error) {
  db.UpdateTable("model")
  selector := bson.M{"modelmd5": modelMd5}
  return db.QueryCount(selector)
}

/**
  查询model by modelMd5
*/
func QueryModelInfo(modelMd5 string) (error, interfaces.Model) {
  var data interfaces.Model
  db.UpdateTable("model")
  selector := bson.M{"modelmd5": modelMd5}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return err, data
}

/**
  查询Node节点信息
*/
func QueryNodeInfo(name string) (error, interfaces.NetInfo) {
  var data interfaces.NetInfo
  db.UpdateTable("node")
  selector := bson.M{"name": name}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return err, data
}

func QueryModelKey(username string) ([]interfaces.ModelKey, error) {
  var models []interfaces.ModelKey
  db.UpdateTable("modelKey")
  cursor, ctx, err := db.QueryAll(bson.M{"username":username}, nil)
  if err != nil {
    return models, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &models)
  return models, err
}
