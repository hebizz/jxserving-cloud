package v1Branch

import (
  "archive/zip"
  "encoding/xml"
  "io"
  "io/ioutil"
  "mime/multipart"
  "os"
  "path/filepath"
  "strings"

  "gitlab.jiangxingai.com/jxserving/components/datasets/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

func UploadFile(val int, fileName string, inputFile multipart.File) (string, error) {
  var path string
  if val == 1 {
    // voc 格式zip包
    path = config.VOCDataUploadZipPath()
  } else {
    // coco
    path = config.COCODataUploadPath()
  }
  if !utils.PathOrFileExists(path) {
    err := os.MkdirAll(path, os.ModePerm)
    //os.Chmod(path, os.ModePerm)
    if err != nil {
      return "", err
    }
  }
  fileName = path + fileName
  out, err := os.Create(fileName)
  defer out.Close()
  defer inputFile.Close()
  _, err = io.Copy(out, inputFile)
  return fileName, err
}

/**
  解压文件
*/
func Unzip(zipFile string, destDir string) error {
  zipReader, err := zip.OpenReader(zipFile)
  if err != nil {
    return err
  }
  defer zipReader.Close()

  for _, f := range zipReader.File {
    fPath := filepath.Join(destDir, f.Name)
    if f.FileInfo().IsDir() {
      os.MkdirAll(fPath, os.ModePerm)
    } else {
      if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
        return err
      }
      inFile, err := f.Open()
      if err != nil {
        return err
      }
      defer inFile.Close()
      outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
      if err != nil {
        return err
      }
      defer outFile.Close()
      _, err = io.Copy(outFile, inFile)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

/**
  解析xml文件
*/
func FormatXMLFiles(filePath string) ([]Annotation, error) {
  var annoList []Annotation
  filePath = filePath + "/Annotations/"
  err := os.Chdir(filePath)
  if err != nil {
    return annoList, err
  }
  files, err := ioutil.ReadDir(filePath)
  if err != nil {
    return annoList, err
  }
  for _, fi := range files {
    if !strings.HasSuffix(fi.Name(), "xml") {
      continue
    }
    // 解析xml
    annotation, err := formatXml(fi.Name())
    if err != nil {
      log.Error("error: ", err, " fileName:", fi.Name())
      continue
    }
    annoList = append(annoList, annotation)
  }
  return annoList, err
}

func formatXml(fileName string) (Annotation, error) {
  obj := Annotation{}
  file, err := os.Open(fileName) // For read access.
  if err != nil {
    return obj, err
  }
  defer file.Close()
  data, err := ioutil.ReadAll(file)
  if err != nil {
    return obj, err
  }
  err = xml.Unmarshal(data, &obj)
  // 由于标注数据xml中对应的fileName与文件本身的fileName不一致
  // 以文件本身的fileName为准
  imageSuffixArray := strings.Split(obj.FileName, ".")
  var imageSuffix = "jpg"
  if len(imageSuffixArray) == 2 {
    imageSuffix = imageSuffixArray[1]
  }
  reImageName := strings.Replace(fileName, "xml", imageSuffix, 1)
  obj.FileName = reImageName
  return obj, err
}
