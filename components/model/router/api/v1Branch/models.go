package v1Branch

import (
  "archive/tar"
  "bufio"
  "compress/gzip"
  "crypto/md5"
  "encoding/hex"
  "io"
  "os"
  "path/filepath"
  "strings"

  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

/**
  计算model的MD5值
*/
func CalculateHash(filePath string) (string, error) {
  file, err := os.Open(filePath)
  if err != nil {
    log.Error(err)
    return "", err
  }
  defer file.Close()
  hash := md5.New()
  if _, err = io.Copy(hash, file); err != nil {
    return "", err
  }
  hashInBytes := hash.Sum(nil)[:16]
  md5Str := hex.EncodeToString(hashInBytes)
  return md5Str, nil
}

/**
  解压 tar.gz
*/
func DeCompressTar(tarFilePath, dest string) (string, error) {
  srcFile, err := os.Open(tarFilePath)
  if err != nil {
    return "", err
  }
  defer srcFile.Close()
  gr, err := gzip.NewReader(srcFile)
  if err != nil {
    return "", err
  }
  defer gr.Close()
  tr := tar.NewReader(gr)
  var filePath string
  for {
    hdr, err := tr.Next()
    if err != nil {
      if err == io.EOF {
        break
      } else {
        return "", err
      }
    }
    filename := dest + hdr.Name
    filePath = filename
    file, _ := createFile(filename)
    _, _ = io.Copy(file, tr)
  }
  return filePath, nil
}

func createFile(name string) (*os.File, error) {
  err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
  if err != nil {
    return nil, err
  }
  return os.Create(name)
}

/**
  压缩 使用gzip压缩成tar.gz
*/

func TarFilesDirs(path string, tarFilePath string, modelName string) error {
  if !utils.PathOrFileExists(tarFilePath) {
    err := os.MkdirAll(tarFilePath, 0755)
    if err != nil {
      return err
    }
  }
  tarFilePath = tarFilePath + modelName
  file, err := os.Create(tarFilePath)
  if err != nil {
    return err
  }
  defer file.Close()
  gz := gzip.NewWriter(file)
  defer gz.Close()
  tw := tar.NewWriter(gz)
  defer tw.Close()
  if err := tarit(path, tw); err != nil {
    return err
  }
  return nil
}

func tarit(source string, tw *tar.Writer) error {
  info, err := os.Stat(source)
  if err != nil {
    return nil
  }
  var baseDir string
  if info.IsDir() {
    baseDir = filepath.Base(source)
  }
  return filepath.Walk(source,
    func(path string, info os.FileInfo, err error) error {
      if err != nil {
        return err
      }
      var link string
      if info.Mode()&os.ModeSymlink != 0 {
        if link, err = os.Readlink(path); err != nil {
          return err
        }
      }
      header, err := tar.FileInfoHeader(info, link)
      if err != nil {
        return err
      }
      if baseDir != "" {
        header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
      }
      if !info.Mode().IsRegular() { //nothing more to do for non-regular
        return nil
      }
      if err := tw.WriteHeader(header); err != nil {
        return err
      }
      if info.IsDir() {
        return nil
      }
      file, err := os.Open(path)
      if err != nil {
        return err
      }
      defer file.Close()
      buf := make([]byte, 16)
      if _, err = io.CopyBuffer(tw, file, buf); err != nil {
        return err
      }
      return nil
    })
}

/**
  read json
*/
func ReadJsonToStr(configJson string) (string, error) {
  file, err := os.Open(configJson)
  if err != nil {
    return "", err
  }
  defer file.Close()
  var line string
  inputReader := bufio.NewReader(file)
  for {
    str, err := inputReader.ReadString('\n')
    if err == io.EOF {
      break
    }
    line = line + str
  }
  return line, nil
}

/**
  将数据写入到制定json文件
*/
func WriteStrToJson(targetJsonPath string, data []byte) error {
  file, err := os.OpenFile(targetJsonPath, os.O_RDWR|os.O_CREATE, 0755)
  if err != nil {
    return err
  }
  defer file.Close()
  _, err = file.Write(data)
  return err
}
