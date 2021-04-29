package v1

import (
  "fmt"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"

  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/jxserving/components/datasets/models/database"
  "gitlab.jiangxingai.com/jxserving/components/datasets/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/datasets/router/api/v1Branch"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

const (
  TypeDatasetUploadVoc  = 1
  TypeDatasetUploadCoco = 2
)

/**
  上传voc zip 文件
*/
func UploadDataSetsForVocHandler(c *gin.Context) {
  file, header, err := c.Request.FormFile("file")
  typeValue := c.Request.FormValue("type")
  name := c.Request.FormValue("name")
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  val, _ := strconv.Atoi(typeValue)
  if val != TypeDatasetUploadVoc && val != TypeDatasetUploadCoco {
    response.ServerError(c, http.StatusBadRequest, "upload file type error,type must be 1 or 2", nil, "")
    return
  }

  fileName, err := v1Branch.UploadFile(val, header.Filename, file)
  if err == nil {
    // 解压文件
    unZipPath := config.VOCDataUploadUnzipPath()
    err := v1Branch.Unzip(fileName, unZipPath)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    // 存储到数据库
    spA := strings.Split(header.Filename, ".")
    xmlFilePath := unZipPath + spA[0]
    annotations, err := v1Branch.FormatXMLFiles(xmlFilePath)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    // 将数据插入dataSet
    for _, annotation := range annotations {
      // 插入rawData数据
      var raw interfaces.Data
      raw.Name = annotation.FileName
      //将图片保存到指定路径
      srcImageName := xmlFilePath + "/JPEGImages/" + annotation.FileName
      targetPath := config.ImageSavePath() + utils.GetDate()
      err = v1Branch.CopyImageToTarget(targetPath, srcImageName, annotation.FileName)
      if err != nil {
        log.Error("save image  error::", err)
        continue
      }
      raw.Path = targetPath + annotation.FileName
      var labelList []interfaces.Label
      for _, obj := range annotation.Object {
        var label interfaces.Label
        label.N = obj.Name
        label.T = 1
        var marks []string
        marks = append(marks, obj.BndBox.XMin, obj.BndBox.YMin, obj.BndBox.XMax, obj.BndBox.YMax)
        label.D = marks
        labelList = append(labelList, label)
      }
      raw.Label = labelList
      raw.TimeStamp = time.Now().Unix()
      raw.LabelName = name
      id, err := database.InsertRawData(raw)
      if err != nil {
        log.Error("insert raw data error::", err)
        continue
      }
      // 更新dataSet
      //name := annotation.Object[0].Name
      dataSets, _ := database.QueryDataSetDataByName(name)
      if dataSets.Name == "" {
        // 先插入数据
        dataSets.Name = name
        dataSets.CreatedTimestamp = time.Now().Unix()
        dataSets.Sets = append(dataSets.Sets, id)
        err := database.InsertDataSetsData(*dataSets)
        if err != nil {
          log.Error("insert dataSets error::", err)
          continue
        }
      } else {
        err := database.UpdateDataSetsData(name, id)
        if err != nil {
          log.Error("update dataSets error::", err)
          continue
        }
      }
    }
    response.Success(c, "upload file success", "")
  } else {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
  }
}

/**
  上传voc zip 文件,然后通过name进行分类
*/
func UploadDataSetsClassifyByNameHandler(c *gin.Context) {
  file, header, err := c.Request.FormFile("file")
  typeValue := c.Request.FormValue("type")
  name := c.Request.FormValue("name")
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  val, _ := strconv.Atoi(typeValue)
  fileName, errSave := v1Branch.UploadFile(val, header.Filename, file)
  if errSave == nil {
    // 解压文件
    unZipPath := config.VOCDataUploadUnzipPath()
    errSave = v1Branch.Unzip(fileName, unZipPath)
    // 存储到数据库
    spA := strings.Split(header.Filename, ".")
    xmlFilePath := unZipPath + spA[0]
    annotations, err := v1Branch.FormatXMLFiles(xmlFilePath)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    // 将数据插入dataSet
    for _, annotation := range annotations {
      // 插入rawData数据
      var raw interfaces.Data
      raw.Name = annotation.FileName
      //将图片保存到指定路径
      srcImageName := xmlFilePath + "/JPEGImages/" + annotation.FileName
      targetPath := config.ImageSavePath() + utils.GetDate()
      err = v1Branch.CopyImageToTarget(targetPath, srcImageName, annotation.FileName)
      if err != nil {
        log.Error("save image  error::", err)
        continue
      }
      raw.Path = targetPath + annotation.FileName
      var labelList []interfaces.Label
      for _, obj := range annotation.Object {
        var label interfaces.Label
        label.N = obj.Name
        label.T = 1
        var marks []string
        marks = append(marks, obj.BndBox.XMin, obj.BndBox.YMin, obj.BndBox.YMin, obj.BndBox.YMax)
        label.D = marks
        labelList = append(labelList, label)
      }
      raw.Label = labelList
      raw.LabelName = name
      raw.TimeStamp = time.Now().Unix()
      id, err := database.InsertRawData(raw)
      if err != nil {
        log.Error("insert raw data error::", err)
        continue
      }
      // 更新dataSet
      // todo 处理多个标注的场景
      name := annotation.Object[0].Name
      dataSets, _ := database.QueryDataSetDataByName(name)
      if dataSets.Name == "" {
        // 先插入数据
        dataSets.Name = name
        dataSets.CreatedTimestamp = time.Now().Unix()
        dataSets.Sets = append(dataSets.Sets, id)
        err := database.InsertDataSetsData(*dataSets)
        if err != nil {
          log.Error("insert dataSets error::", err)
          continue
        }
      } else {
        err = database.UpdateDataSetsData(name, id)
        if err != nil {
          log.Error("update dataSets error::", err)
          continue
        }
      }
    }
  }
  response.Success(c, "upload file success", "")
}

/**
  dataset 下载
*/
func DownloadDataSetHandler(c *gin.Context) {
  var body DataSetDownLoadBody
  err := c.BindJSON(&body)
  log.Info("down load all file request: ", body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  // 通过传入的name获取dataSet
  dataSets, err := database.QueryDataSetDataByName(body.Name)
  if err != nil {
    log.Error("query dataset by name error:", err)
  }
  // 取出数据并生成压缩包
  if len(dataSets.Sets) == 0 {
    log.Info(fmt.Sprintf("%s raw data not found, maybe name is wrong", body.Name))
    response.Success(c, fmt.Sprintf("download %s file success", body.Name), &DownLoadResponseBody{Path: ""})
    return
  }
  var rawDataList []interfaces.Data
  if body.DownType == 1 {
    // 生成voc格式数据包
    rawDataList, err = database.QueryRawDataByIdList(dataSets.Sets)
    if err != nil {
      log.Error("query rawData by idList error::", err)
    }
  } else {
    ts := utils.DateTimeStr2TimeStamp(body.StartTime)
    te := utils.DateTimeStr2TimeStamp(body.EndTime)
    rawDataList, err = database.QueryRawDataByTimeOffset(ts, te, dataSets.Sets)
  }
  if len(rawDataList) == 0 {
    log.Info(fmt.Sprintf("%s raw data not found, maybe timestamp is wrong", body.Name))
    response.Success(c, fmt.Sprintf("download %s file success", body.Name), &DownLoadResponseBody{Path: ""})
    return
  }
  if body.Type == 1 {
    var xmlNameList []string
    var filePath string
    var xmlName string
    timeStampPath := utils.GetMD5Str(time.Now().Unix())
    for i := 0; i < len(rawDataList); i++ {
      xmlName, err = v1Branch.GenerateXML(body.Name, rawDataList[i].Name, timeStampPath, rawDataList[i].Path, rawDataList[i].Label)
      if err == nil || xmlName == "" {
        xmlNameList = append(xmlNameList, xmlName)
      } else {
        log.Error("generateXM  error::", err)
      }
      // 保存图片到
      targetPath := config.VOCDataSavePath() + timeStampPath + "/JPEGImages/"
      err = v1Branch.CopyImageToTarget(targetPath, rawDataList[i].Path, rawDataList[i].Name)
      if err != nil {
        log.Error("save image error:: ", err)
      }
    }
    // 保存file name 到 val.txt
    if err := v1Branch.SaveXmlFileNameToTxt(timeStampPath, xmlNameList); err != nil {
      log.Error("save xml name to txt error:", err)
    }

    // 压缩文件
    filePath, err = v1Branch.Zip(timeStampPath)
    if err != nil {
      response.ServerError(c, http.StatusInternalServerError, "", err, "")
      return
    }
    // 删除生成的文件目录,二十分钟后删除.zip文件
    if filePath != "" {
      removePaths := strings.Split(filePath, ".")
      err := os.RemoveAll(removePaths[0])
      if err != nil {
        log.Error("remove file error::", err)
      }
    }
    response.Success(c, fmt.Sprintf("download %s file success", body.Name), &DownLoadResponseBody{Path: filePath})
  } else {
    // 生成coco格式数据包
    response.Success(c, "目前暂不支持下载coco数据包", &DownLoadResponseBody{Path: ""})
  }
}

/**
  获取数据集列表数据
*/
func QueryDataSetsHandler(c *gin.Context) {
  dataSets, err := database.QueryAllDatasetData()
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "/")
    return
  }
  var dataSetList []DataSetBody
  for _, dataSet := range dataSets {
    total := int64(len(dataSet.Sets))
    // 查询dataSet.Name 对应的元数据中未标注的个数
    count, err := database.QueryAllNoneDataCount(dataSet.Name)
    if err != nil {
      log.Error("query all none data count error:", err)
    }
    marker := total - count
    if marker < 0 {
      marker = 0
    }
    dataSetBody := &DataSetBody{Name: dataSet.Name, Marker: marker, Total: total}
    dataSetList = append(dataSetList, *dataSetBody)
  }
  response.Success(c, "query dataset list success", dataSetList)
}

/**
  删除数据集
*/
func DeleteDataSetsHandler(c *gin.Context) {
  var body DeleteDataSetBody
  err := c.BindJSON(&body)
  if err != nil {
    response.ServerError(c, http.StatusBadRequest, "", err, "")
    return
  }
  err = database.DeleteDataSetsData(body.Name)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, fmt.Sprintf("delete dataSets name=%s success", body.Name), "")
}
