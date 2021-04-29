package database

import (
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
)

/**
  插入上报数据到对应的table中
*/
func InsertReportData(data interfaces.ReportData) error {
  db.UpdateTable(data.Project)
  return db.Insert(data)
}
