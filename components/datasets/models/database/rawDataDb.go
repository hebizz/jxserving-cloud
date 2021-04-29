package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo/options"
)

/**
  插入raw data数据
*/
func InsertRawData(data interfaces.Data) (string, error) {
  objectId := primitive.NewObjectID()
  data.Id = objectId
  db.UpdateTable("data")
  err := db.Insert(data)
  return objectId.Hex(), err
}

/**
  通过id更新
*/
func UpdateRawData(data interfaces.Data) error {
  db.UpdateTable("data")
  selector := bson.M{"$set": bson.M{"_id": data.Id}}
  return db.Update(selector, data)
}

/**
  通过id更新
*/
func ReplaceRawData(data interfaces.Data) error {
  db.UpdateTable("data")
  selector := bson.M{"_id": data.Id}
  return db.ReplaceOne(selector, data)
}

/**
  label 为空的数据个数
*/
func QueryAllNoneDataCount(name string) (int64, error) {
  selector := bson.M{"label": bson.M{"$size": 0}, "labelname": name}
  db.UpdateTable("data")
  count, err := db.QueryCount(selector)
  if err != nil {
    return 0, err
  }
  return count, err
}

/**
  通过name 查询该集合里所有数据 & 分页加载
*/
func QueryRawDataListByIndex(startIndex int64, offset int64, sort string, idList []string) ([]interfaces.Data, error) {
  var rawDataList []interfaces.Data
  var idList2 []primitive.ObjectID
  for _, idStr := range idList {
    id, _ := primitive.ObjectIDFromHex(idStr)
    idList2 = append(idList2, id)
  }
  selector := bson.M{"_id": bson.M{"$in": idList2}}
  db.UpdateTable("data")
  sortM := bson.D{{sort, -1}} // -1 降序排列, 1 升序
  findOptions := options.Find().SetSort(sortM).SetLimit(offset).SetSkip(startIndex)
  cursor, ctx, err := db.QueryAll(selector, findOptions)
  if err != nil {
    return rawDataList, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &rawDataList)
  return rawDataList, err
}

/**
  查询所有data数据 by idList
*/
func QueryRawDataByIdList(idList []string) ([]interfaces.Data, error) {
  var rawDataList []interfaces.Data
  var idList2 []primitive.ObjectID
  for _, idStr := range idList {
    id, _ := primitive.ObjectIDFromHex(idStr)
    idList2 = append(idList2, id)
  }
  selector := bson.M{"_id": bson.M{"$in": idList2}}
  db.UpdateTable("data")
  cursor, ctx, err := db.QueryAll(selector, nil)
  if err != nil {
    return rawDataList, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &rawDataList)
  return rawDataList, err
}

/**
  查询所有data数据 by timestamp offset
*/
func QueryRawDataByTimeOffset(startTime int64, endTime int64, idList []string) ([]interfaces.Data, error) {
  var rawDataList []interfaces.Data
  var idList2 []primitive.ObjectID
  for _, idStr := range idList {
    id, _ := primitive.ObjectIDFromHex(idStr)
    idList2 = append(idList2, id)
  }
  selector := bson.M{"_id": bson.M{"$in": idList2}, "timestamp": bson.M{"$gte": startTime, "$lte": endTime}}
  db.UpdateTable("data")
  cursor, ctx, err := db.QueryAll(selector, nil)
  if err != nil {
    return rawDataList, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &rawDataList)
  return rawDataList, err
}
