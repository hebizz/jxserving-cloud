package v1Branch

import (
  "archive/zip"
  "encoding/xml"
  "errors"
  "image"
  _ "image/jpeg" //  image 格式，如果没有导入该依赖包，会提示 unknown format
  "io"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"

  "gitlab.jiangxingai.com/jxserving/components/datasets/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

/**
  generate xml file
*/
func GenerateXML(labelName string, imageName string, timeStampPath string, imagePath string, labels []interfaces.Label) (string, error) {
  var root Annotation
  root.Folder = labelName
  root.Source = &Source{"dataSets"}
  fileName := strings.Split(imageName, ".")[0]
  root.Path = imagePath
  // 获取图片的宽高
  width, height := getImageInfo(imagePath)

  root.Size = &Size{Width: width, Height: height, Depth: 3} // 彩图 depth=3
  root.Segmented = 0
  for _, label := range labels {
    if len(label.D) == 0 {
      continue
    }
    xmin := label.D[0]
    ymin := label.D[1]
    xmax := label.D[2]
    ymax := label.D[3]
    root.Object = append(root.Object, &Object{Name: label.N, Pose: "Unspecified",
      Difficult: 0, Truncated: 0, BndBox: &BndBox{XMin: xmin, YMin: ymin, XMax: xmax, YMax: ymax}})
  }
  xmlOutput, err := xml.MarshalIndent(root, "", "    ")
  if err != nil {
    return "", err
  }
  headerBytes := []byte(xml.Header)
  xmlData := append(headerBytes, xmlOutput...)
  path := config.VOCDataSavePath() + timeStampPath + "/Annotations/"
  if !utils.PathOrFileExists(path) {
    os.MkdirAll(path, os.ModePerm)
    os.Chmod(path, os.ModePerm)
  }
  XMLFilePathName := path + fileName + ".xml"
  err = ioutil.WriteFile(XMLFilePathName, xmlData, os.ModePerm)
  XMLFileName := fileName + ".xml"
  return XMLFileName, err
}

/**
  获取图片的宽高
*/
func getImageInfo(path string) (int, int) {
  if path == "" {
    return 0, 0
  }
  file, err := os.Open(path)
  defer file.Close()
  c, _, err := image.DecodeConfig(file)
  if err != nil {
    log.Error("format image error:", err)
    return 0, 0
  }
  width := c.Width
  height := c.Height
  return width, height
}

/**
  保存文件名到val.txt
*/
func SaveXmlFileNameToTxt(timeStampPath string, xmFileNameList []string) error {
  path := config.VOCDataSavePath() + timeStampPath + "/ImageSets/Main/"
  if !utils.PathOrFileExists(path) {
    os.Chmod(path, os.ModePerm)
    os.MkdirAll(path, os.ModePerm)
  }
  xmlValPath := path + "val.txt"
  var f *os.File
  var err error
  if !utils.PathOrFileExists(xmlValPath) {
    f, err = os.Create(xmlValPath)
  } else {
    f, err = os.OpenFile(xmlValPath, os.O_APPEND, 0666)
  }
  for _, fileName := range xmFileNameList {
    kv := strings.Split(fileName, ".")
    _, err = io.WriteString(f, kv[0]+"\n")
    if err != nil {
      return err
    }
  }
  return err
}

/**
  copy image 到JPEGImages
*/
func CopyImageToTarget(targetPath string, srcAbName string, srcFileName string) error {
  src, err := os.Open(srcAbName)
  if err != nil {
    return err
  }
  if !utils.PathOrFileExists(targetPath) {
    os.Chmod(targetPath, os.ModePerm)
    os.MkdirAll(targetPath, os.ModePerm)
  }
  targetName := targetPath + srcFileName
  f, err := os.OpenFile(targetName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
  if err != nil {
    return err
  }
  defer f.Close()
  _, err = io.Copy(f, src)
  return err
}

/*
  zip file
*/
func Zip(timeStampPath string) (string, error) {

  path := config.VOCDataSavePath() + timeStampPath
  if !utils.PathOrFileExists(path) {
    return "", errors.New("this file path not exist")
  }
  os.Chdir(config.VOCDataSavePath())
  destZip := timeStampPath + ".zip"
  srcFile := timeStampPath
  filePath := config.VOCDataSavePath() + destZip
  zipfile, err := os.Create(destZip)
  if err != nil {
    return "", err
  }
  defer zipfile.Close()

  archive := zip.NewWriter(zipfile)
  defer archive.Close()

  walkErr := filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    header, err := zip.FileInfoHeader(info)
    if err != nil {
      return err
    }
    header.Name = path
    if info.IsDir() {
      header.Name += "/"
    } else {
      header.Method = zip.Deflate
    }
    writer, err := archive.CreateHeader(header)
    if err != nil {
      return err
    }
    if !info.IsDir() {
      file, err := os.Open(path)
      if err != nil {
        return err
      }
      defer file.Close()
      _, err = io.Copy(writer, file)
    }
    return err
  })
  return filePath, walkErr
}
