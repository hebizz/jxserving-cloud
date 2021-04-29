package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

/**
  插入 target 数据
*/
func InsertTargetData(data interfaces.Target) (string, error) {
  objectId := primitive.NewObjectID()
  data.Id = objectId
  db.UpdateTable("target")
  err := db.Insert(data)
  return objectId.Hex(), err
}

/**
  更新 target 数据
*/
func UpdateTargetData(data interfaces.Target) error {
  id, _ := primitive.ObjectIDFromHex(data.LabelId)
  selector := bson.M{"_id": id}
  updateData := bson.M{"$set": bson.M{"name": data.Name, "d": data.D}}
  db.UpdateTable("target")
  return db.Update(selector, updateData)
}

/**
  delete target 数据
*/
func DeleteTargetData(id string) error {
  _id, _ := primitive.ObjectIDFromHex(id)
  selector := bson.M{"_id": _id}
  db.UpdateTable("target")
  return db.Delete(selector)
}

/**
  查询 [] target
*/
func QueryTargetData() ([]interfaces.Target, error) {
  var labelSets []interfaces.Target
  db.UpdateTable("target")
  cursor, ctx, err := db.QueryAll(bson.M{}, nil)
  if err != nil {
    return labelSets, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &labelSets)
  return labelSets, err
}

/**
  查询target
*/
func QueryTarget(id string) (error, *interfaces.Target) {
  data := new(interfaces.Target)
  db.UpdateTable("target")
  _id, _ := primitive.ObjectIDFromHex(id)
  selector := bson.M{"_id": _id}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return err, data
}

/**
  查询model document 中是否已经评价关联了该target
*/
func QueryTargetIsRelateByModel(label string) error {
  data := new(interfaces.Model)
  db.UpdateTable("model")
  selector := bson.M{"relatedlabel": label}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return err
}
