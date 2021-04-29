package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
)

/**
  插入dataSets数据
*/
func InsertDataSetsData(data interfaces.Dataset) error {
  db.UpdateTable("dataSets")
  return db.Insert(data)
}

/**
  通过name 查询该集合里是否存在该文档并返回
*/
func QueryDataSetDataByName(name string) (*interfaces.Dataset, error) {
  data := new(interfaces.Dataset)
  db.UpdateTable("dataSets")
  selector := bson.M{"name": name}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return data, err
}

/**
  查询所有dataSets 数据
*/
func QueryAllDatasetData() ([]interfaces.Dataset, error) {
  var dataList [] interfaces.Dataset
  db.UpdateTable("dataSets")
  cursor, ctx, err := db.QueryAll(bson.M{}, nil)
  if err != nil {
    return dataList, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &dataList)
  return dataList, err
}

/**
  更新dataSet 数据库: 向sets[]添加数据
*/
func UpdateDataSetsData(name string, dataId string) error {
  selector := bson.M{"name": name}
  data := bson.M{"$push": bson.M{"sets": dataId}}
  db.UpdateTable("dataSets")
  return db.Update(selector, data)
}

/**
  删除dataSet元数据
*/
func DeleteDataSetsData(name string) error {
  selector := bson.M{"name": name}
  db.UpdateTable("dataSets")
  return db.Delete(selector)
}
