package v1

type Tensor struct {
  Output []string `json:"output"`
  Input  []string `json:"input"`
}

type AuthId struct {
  ModelMd5 string `json:"modelMd5" binding:"required"`
}

type ModelId struct {
  ModelId string `json:"modelId" binding:"required"`
}

type ModelDownLoad struct {
  ModelMd5 string `json:"modelMd5" binding:"required"`
  LabelMd5 string `json:"labelMd5" binding:"required"`
}

type ModelBuiltInConfig struct {
  Tensors   Tensor    `json:"tensors"`
  Labels    []string  `json:"labels"`
  Threshold []float64 `json:"threshold"`
  Mapping   [] string `json:"mapping"`
  Timestamp int64     `json:"timestamp"`
  Md5       string    `json:"md5"`
}

type SubModelRequestBody struct {
  ModelId   string    `json:"modelId" binding:"required"`
  Threshold []float64 `json:"threshold" binding:"required"`
  Mapping   [] string `json:"mapping" binding:"required"`
}

type ModelAttributeBody struct {
  TargetIp  string    `json:"targetIp" binding:"required"`
  ModelPath string    `json:"modelPath" binding:"required"`
  Threshold []float64 `json:"threshold" binding:"required"`
  Mapping   [] string `json:"mapping" binding:"required"`
  Timestamp int64     `json:"timestamp" binding:"required"`
  Md5       string    `json:"md5" binding:"required"`
}

type QueryListModelBody struct {
  TargetIp string `json:"targetIp" binding:"required"`
}
type CreateModelBody struct {
  TargetIp string `json:"targetIp" binding:"required"`
  Bid      string `json:"bid"`
  Btype    string `json:"btype"`
  Model    string `json:"model"`
  Version  string `json:"version"`
  Mode     string `json:"model"`
  Extra    string `json:"extra"`
}

