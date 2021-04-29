package v1

import (
  "encoding/json"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/model/models/database"
  "gitlab.jiangxingai.com/jxserving/components/model/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/model/router/api/v1Branch"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

/**
  查询models列表
*/
func QueryModelsHandler(c *gin.Context) {
  var body []interfaces.ModelInfo
  var modelInfo interfaces.ModelInfo
  models, err := database.QueryModelsData()
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  for _, model := range models {
    model.ModelId = model.Id.Hex()
    modelInfo.Model = model
    labelSet, err := database.QueryLabelsInfoById(model.TargetId)
    if err != nil {
      log.Error("query labels info by id error:", err)
      continue
    }
    if labelSet.Name == "" {
      log.Error("label name is nil")
      continue
    }
    modelInfo.Label = labelSet
    body = append(body, modelInfo)
  }
  response.Success(c, "query model list success", body)
}

/**
  查询模型子模型
*/
func QuerySubModelListHandler(c *gin.Context) {
  var body ModelId
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  modelConfig, err := database.QueryModelConfigsInfo(body.ModelId)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  response.Success(c, "query sub model list success", modelConfig.Configs)
}

func QueryModelKeyHandler(c *gin.Context) {
  var body interfaces.ModelKey
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  data, err := database.QueryModelKey(body.Username)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  response.Success(c, "query  model key success", data)
}

/**
  上传模型信息
*/
func UploadModelHandler(c *gin.Context) {

  header, err := c.FormFile("file")
  framework := c.Request.FormValue("framework")
  notes := c.Request.FormValue("notes")
  version := c.Request.FormValue("version")
  targetId := c.Request.FormValue("targetId")
  threshold := c.Request.FormValue("threshold")
  mapping := c.Request.FormValue("mapping")

  if err != nil || framework == "" || notes == "" || version == "" {
    log.Error("upload file error:", err)
    response.ServerError(c, http.StatusBadRequest, "request params error", nil, "")
    return
  }
  labelSet, err := database.QueryLabelsInfoById(targetId)
  if err != nil || labelSet.Name == "" {
    response.ServerError(c, http.StatusBadRequest, "targetId is wrong", nil, "")
    return
  }
  originName := header.Filename
  if !strings.HasSuffix(originName, ".tar.gz") {
    response.ServerError(c, http.StatusBadRequest, "Unsupported file format type, tar.gz only", nil, "")
    return
  }
  dir := config.ModelSavePath() + utils.GetMD5Str(time.Now().UnixNano()/1e6) + "/"
  path := dir + version + "_" + framework + ".tar.gz"
  err = os.MkdirAll(dir, os.ModePerm)
  if err != nil {
    log.Error("create file dir error:", err)
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  err = c.SaveUploadedFile(header, path)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    err := os.RemoveAll(dir)
    log.Error("remove dir error:", err)
    return
  }
  modelMd5, _ := v1Branch.CalculateHash(path)
  existByMd5, err := database.QueryModelIsExistByMd5(modelMd5)
  if err != nil || existByMd5 > 0 {
    response.ServerError(c, http.StatusForbidden, "model already exists", nil, modelMd5)
    err := os.RemoveAll(dir)
    if err != nil {
      log.Error("remove dir error:", err)
    }
    return
  }
  fileInfo, _ := os.Stat(path)
  size := fileInfo.Size()

  model := interfaces.Model{
    TargetId:    targetId,
    RawName:     originName,
    FrameWork:   framework,
    Version:     version,
    ModelMd5:    modelMd5,
    ModelPath:   path,
    IsPublished: false,
    Notes:       notes,
    Size:        size,
    Timestamp:   time.Now().Unix(),
  }
  modelId, err := database.InsertModelData(model)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  // 插入configs信息
  configMd5 := utils.GetMD5Str(time.Now().UnixNano() / 1e6)
  mps := strings.Split(mapping, ",")
  ths := strings.Split(threshold, ",")
  var thresholds []float64
  for _, th := range ths {
    f, _ := strconv.ParseFloat(th, 64)
    thresholds = append(thresholds, f)
  }
  cf := interfaces.Config{Threshold: thresholds, Mapping: mps, Md5: configMd5, Timestamp: time.Now().Unix()}
  modelConfig := interfaces.ModelConfig{ModelId: modelId, Configs: []interfaces.Config{cf}}
  err = database.InsertModelConfigsInfo(modelConfig)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }

  response.Success(c, "upload model success", modelMd5)
}

/**
  迭代子模型
*/
func UpdateSubModelHandler(c *gin.Context) {
  var body SubModelRequestBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  cf := interfaces.Config{Threshold: body.Threshold, Mapping: body.Mapping, Md5: utils.GetMD5Str(time.Now().UnixNano()), Timestamp: time.Now().Unix()}
  updateResult, err := database.UpdateModelConfigsInfo(body.ModelId, cf)
  if err != nil || updateResult.MatchedCount == 0 {
    response.ServerError(c, http.StatusBadRequest, "invalid modelId", nil, "")
    return
  }
  response.Success(c, "update sub model success", "")
}

/**
  发布模型
*/
func PublishModelHandler(c *gin.Context) {
  var body AuthId
  err := c.ShouldBindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  err = database.UpdateLabelSetDataByMd5(body.ModelMd5)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "publish success", "")
}

/**
  下载离线模型
*/
func DownloadModelOffLineHandler(c *gin.Context) {
  var body ModelDownLoad
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  err, model := database.QueryModelInfo(body.ModelMd5)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "modelMd5 is invalid", nil, "")
    return
  }
  // 文件解压缩,然后重新配置distros.json
  destPath := config.ModelSavePath() + "unTar/"
  filePath, err := v1Branch.DeCompressTar(model.ModelPath, destPath)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  path := destPath + strings.Split(strings.Split(filePath, "unTar/")[1], "/")[0]
  // 创建distors.json
  modelPath := string([]rune(filePath)[0:strings.LastIndex(filePath, "/")])
  configJson := modelPath + "/config.json"

  if !utils.PathOrFileExists(configJson) {
    response.ServerError(c, http.StatusInternalServerError, "config.json not found", err, "")
    return
  }
  line, err := v1Branch.ReadJsonToStr(configJson)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  var configJsonBody ModelBuiltInConfig
  err = json.Unmarshal([]byte(line), &configJsonBody)
  if err != nil {
    log.Error("unMarshal error", err)
  }
  dbModelConfig, err := database.QueryModelConfigsInfo(model.Id.Hex())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  var matchConfig interfaces.Config
  for _, config := range dbModelConfig.Configs {
    if config.Md5 == body.LabelMd5 {
      matchConfig = config
      break
    }
  }
  configJsonBody.Timestamp = matchConfig.Timestamp
  configJsonBody.Md5 = matchConfig.Md5
  configJsonBody.Mapping = matchConfig.Mapping
  configJsonBody.Threshold = matchConfig.Threshold
  distrosPath := modelPath + "/distros.json"
  data, _ := json.Marshal(configJsonBody)
  err = v1Branch.WriteStrToJson(distrosPath, data)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  tarPath := config.ModelSavePath() + "tar/"
  err = v1Branch.TarFilesDirs(path, tarPath, model.RawName)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "download model success", tarPath+model.RawName)
}

/**
  删除模型
*/
func DeleteModelHandler(c *gin.Context) {
  var body AuthId
  err := c.ShouldBindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
  } else {
    err, model := database.QueryModelInfo(body.ModelMd5)
    if model.ModelPath == "" || err != nil {
      response.ServerError(c, http.StatusBadRequest, "invalid modelMd5", nil, "")
      return
    }
    modelPath := string([]rune(model.ModelPath)[0:strings.LastIndex(model.ModelPath, "/")])
    err = os.RemoveAll(modelPath)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    err = database.DeleteModelInfo(body.ModelMd5)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    err = database.DeleteFromPerformanceData(body.ModelMd5)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    response.Success(c, "delete model success", "")
  }
}
