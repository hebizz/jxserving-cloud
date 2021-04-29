package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/options"
)

/**
  查询修改后的data reMark 状态
*/
func QueryReportDataReMarkStatus(project string, eventId string) (bool, error) {
  var data interfaces.ReportData
  db.UpdateTable(project)
  selector := bson.M{"eventid": eventId}
  singleResult := db.QueryOne(selector)
  err := singleResult.Decode(&data)
  return data.ReMark, err
}

/**
  更新reportData reMark 状态
*/
func UpdateReportDataReMarkStatus(project string, eventId string) error {
  db.UpdateTable(project)
  selector := bson.M{"eventid": eventId}
  updateSelector := bson.M{"$set": bson.M{"remark": true}}
  return db.Update(selector, updateSelector)
}

/**
  查询project list
*/
func QueryProjectList() ([]interfaces.Project, error) {
  db.UpdateTable("project")
  var dataList [] interfaces.Project
  cursor, ctx, err := db.QueryAll(bson.M{}, nil)
  if err != nil {
    return dataList, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &dataList)
  return dataList, err
}

/**
  分页加载未标注的图片列表
*/
func QueryNoMarkedDataByProject(project string, startIndex int64, offset int64, alertTime int64, sort string) ([]interfaces.ReportData, error) {
  var dataList []interfaces.ReportData
  db.UpdateTable(project)
  selector := bson.M{"remark": false, "alerttime": bson.M{"$lte": alertTime}}
  sortM := bson.D{{sort, -1}} // -1 降序排列, 1 升序
  findOptions := options.Find().SetSort(sortM).SetLimit(offset).SetSkip(startIndex - 1)
  cursor, ctx, err := db.QueryAll(selector, findOptions)
  if err != nil {
    return dataList, err
  }
  defer cursor.Close(ctx)
  err = cursor.All(ctx, &dataList)
  return dataList, err
}
