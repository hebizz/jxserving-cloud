package interfaces

import (
  "time"

  influx "github.com/influxdata/influxdb1-client/v2"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Label struct {
  N string   `json:"name"` // label name like fire crane..
  T int      `json:"type"` // 1 rectangle 2 cycle 3 ..
  D []string `json:"position"`
}

type Data struct {
  Id        primitive.ObjectID `bson:"_id" json:"id"`
  Name      string             `json:"name"`
  TimeStamp int64              `json:"timestamp"`
  Path      string             `json:"path"`
  Label     []Label            `json:"label"`
  LabelName string             `json:"labelName"`
}

type AlertPosition struct {
  LeftX  string `json:"left_x" binding:"required"`
  LeftY  string `json:"left_y" binding:"required"`
  RightX string `json:"right_x" binding:"required"`
  RightY string `json:"right_y" binding:"required"`
}

type AppInfo struct {
  AppName    string `json:"app_name" binging:"required"`
  AppVersion string `json:"app_version" binging:"required"`
}

type ReportData struct {
  Id            primitive.ObjectID `bson:"_id"`
  DataId        string             `json:"data_id"`
  Project     string             `json:"cluster_id"`
  Title         string             `json:"title"`
  AlertType     string             `json:"alert_type"`
  AlertPosition []AlertPosition    `json:"alert_position" binding:"required"`
  EventId       string             `json:"event_id" binding:"required"`
  DeviceId      string             `json:"device_id"`
  AlertTime     int64              `json:"time"`
  AppInfo       AppInfo            `json:"app_info" `
  Reliability   float64            `json:"reliability"`
  Description   string             `json:"description"`
  ImagePath     string             `json:"image_path"`
  Image         string             `json:"image"`
  AppId         string             `json:"app_id" `
  CreateTime    int64              `json:"created_time"`
  HostIp        string             `json:"host_ip"`
  ReMark        bool               `json:"reMark"`
}

type Project struct {
  Name string
  Api  string
  Host string
  Port string
}

type AnalysisData struct {
  RawDataId string
  Value     float64
  Label     map[string]string
}

func (d *AnalysisData) ToInfluxPoint(measurement string) (*influx.Point, error) {
  tags := map[string]string{
    "id": d.RawDataId,
  }
  fields := map[string]interface{}{
    "value": d.Value,
  }

  return influx.NewPoint(measurement, tags, fields, time.Now())
}
