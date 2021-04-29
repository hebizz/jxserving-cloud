package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
  Id          primitive.ObjectID `json:"id" bson:"_id"`
  ModelId     string             `json:"modelId"`
  TargetId    string             `json:"targetId"`
  FrameWork   string             `json:"framework" binding:"required"`
  Version     string             `json:"version" binding:"required"`
  Notes       string             `json:"notes" binding:"required"`
  RawName     string             `json:"rawName"`
  ModelMd5    string             `json:"modelMd5"`
  ModelPath   string             `json:"modelPath"`
  IsPublished bool               `json:"isPublished"`
  Size        int64              `json:"size"`
  Timestamp   int64              `json:"timestamp"`
}

type Config struct {
  Threshold []float64 `json:"threshold"`
  Mapping   [] string `json:"mapping"`
  Timestamp int64     `json:"timestamp"`
  Md5       string    `json:"md5"`
}

type ModelConfig struct {
  ModelId string    `json:"modelId"`
  Configs [] Config `json:"configs"`
}

type ModelInfo struct {
  Label Target `json:"label"`
  Model Model    `json:"model"`
}

type SubModelList struct {
  SubModelList [] ModelConfig `json:"subModelList"`
}

type ModelPerformance struct {
  ModelMd5         string
  EvalDataset      string
  Datesettype      int64
  ErrorRate        float64
  LeakRate         float64
  MeanAp           float64
  CreatedTimestamp int64
}

type NetInfo struct {
  Cidr string `json:"cidr"`
}

type Distribute struct {
  Name     string `json:"name" binding:"required"`
  FilePath string `json:"filePath" binding:"required"`
}

type Performance struct {
  ErrorRate float64
  LeakRate  float64
  Map       float64
}

type Modeldate struct {
  Modelmd5 string `json:"modelmd5"`
  Name     string `json:"name"`
  Type     int64  `json:"type"`
  DownType int64  `json:"downtype"`
}

type ModelKey struct {
  Username string `json:"username"`
  Rsa      string `json:"rsa"`
  Private  string `json:"private"`
}