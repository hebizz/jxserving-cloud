package v1

import "gitlab.jiangxingai.com/jxserving/components/interfaces"

type ImageDataBody struct {
  Id        string             `json:"id"binding:"required"`
  TimeStamp int64              `json:"timestamp"binding:"required"`
  ImageStr  string             `json:"image"`
  Path      string             `json:"path"binding:"required"`
  Name      string             `json:"name" binding:"required"`
  LabelName string             `json:"labelName" binding:"required"`
  Label     []interfaces.Label `json:"label"`
}

type LabelSetDeleteBody struct {
  Id string `json:"id" binding:"required"`
}

type DataSetDownLoadBody struct {
  Name      string `json:"name" binding:"required"`
  Type      int64  `json:"type" binding:"required"`     // type=1 : voc / type=2: coco
  DownType  int    `json:"downType" binding:"required"` // downType=1 : all / downType=2 : apart
  StartTime string `json:"startTime"`
  EndTime   string `json:"endTime"`
}

type DownLoadResponseBody struct {
  Path string `json:"path"`
}

type DataSetBody struct {
  Name   string `json:"name"`
  Marker int64  `json:"marker"`
  Total  int64  `json:"total"`
}

type DeleteDataSetBody struct {
  Name string `json:"name" binding:"required"`
}

type DataSetAdminBody struct {
  Name       string `json:"name" binding:"required"`
  StartIndex int64  `json:"startIndex" binding:"required"`
  Offset     int64  `json:"offset" binding:"required"`
}
