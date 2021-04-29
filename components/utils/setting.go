package utils

import (
  "bufio"
  "bytes"
  "container/list"
  "crypto/md5"
  "encoding/base64"
  "encoding/hex"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "os/exec"
  "os/user"
  "regexp"
  "strconv"
  "time"

  "github.com/spf13/cast"
  "github.com/spf13/viper"
  "k8s.io/klog"
)

var (
  viper_ = viper.New()
)

func init() {
  viper_.SetConfigName("config")
  configDir := os.Getenv("JXS_CONFIG")
  if configDir == "" {
    curUser, err := user.Current()
    if err != nil {
      panic("unable to access current user's home directory")
    }
    configDir = curUser.HomeDir
  }
  configPath := fmt.Sprintf("%s%s", configDir, "/config.yaml")
  if !PathOrFileExists(configPath) {
    _, _ = os.Create(configPath)
  }
  viper_.AddConfigPath(configDir)
  viper_.SetConfigType("yaml")

  viper_.AutomaticEnv()
  if err := viper_.ReadInConfig(); err != nil {
    panic(err)
  }
}

func ReadString(key string) string {
  if viper_.IsSet(key) {
    if ret := viper_.Get(key); ret != nil {
      return cast.ToString(ret)
    } else {
      return ""
    }
  }
  return ""
}

func ReadMap(key string) map[string]interface{} {
  if viper_.IsSet(key) {
    if ret := viper_.Get(key); ret != nil {
      return cast.ToStringMap(ret)
    } else {
      return map[string]interface{}{}
    }
  }
  return map[string]interface{}{}
}

func Base64ToImage(bStr string, filePath string, fileName string) (string, string, error) {
  dd, _ := base64.StdEncoding.DecodeString(bStr)

  if !PathOrFileExists(filePath) {
    _ = os.Chmod(filePath, os.ModePerm)
    _ = os.MkdirAll(filePath, os.ModePerm)
  }
  imagePath := fmt.Sprintf("%s%s", filePath, fileName)
  err := ioutil.WriteFile(imagePath, dd, 0667)
  return imagePath, fileName, err
}

func GetMD5Str(timestamp int64) string {
  str := strconv.FormatInt(timestamp, 10)
  h := md5.New()
  h.Write([]byte(str))
  return hex.EncodeToString(h.Sum(nil))
}

func GetMD5ByStr(str string) string {
  h := md5.New()
  h.Write([]byte(str))
  return hex.EncodeToString(h.Sum(nil))
}

func ImageToBase64Str(path string) (string, error) {
  imgFile, err := os.Open(path)
  if err != nil {
    return "", err
  }
  defer imgFile.Close()
  f, _ := imgFile.Stat()
  size := f.Size()
  buf := make([]byte, size)
  fReader := bufio.NewReader(imgFile)
  _, err = fReader.Read(buf)
  if err != nil {
    return "", err
  }
  return base64.StdEncoding.EncodeToString(buf), nil
}

func PathOrFileExists(path string) bool {
  _, err := os.Stat(path)
  if err == nil {
    return true
  }
  if os.IsNotExist(err) {
    return false
  }
  return false
}

/**
  date to timestamp
*/
func DateTimeStr2TimeStamp(dateStr string) int64 {
  dateStr = dateStr + " 23:59:59"
  timeLayout := "2006-01-02 15:04:05"
  loc, _ := time.LoadLocation("Local")
  tmp, _ := time.ParseInLocation(timeLayout, dateStr, loc)
  return tmp.Unix()
}

/**
  获取日期
*/
func GetDate() string {
  year := time.Now().Year()
  month := time.Now().Month()
  day := time.Now().Day()
  return fmt.Sprintf("%d-%d-%d/", year, month, day)
}

/**
  获取终端日志输出
*/
func CommandExecuteLogs(cmd *exec.Cmd) (string, error) {
  var out bytes.Buffer
  var stderr bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &stderr
  err := cmd.Run()
  if err != nil {
    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    return stderr.String(), errors.New(stderr.String())
  }
  klog.Info("cmd execute result is :\n" + out.String())
  return out.String(), nil
}

/**
  执行系统指令并逐行输出打印信息
*/
func ExecCommandAndPrintLines(commandName string, args ...string) (*list.List, error) {
  cmd := exec.Command(commandName, args...)
  //显示运行的命令
  klog.Info("exec command args:", cmd.Args)
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return list.New(), err
  }
  _ = cmd.Start()
  reader := bufio.NewReader(stdout)
  //实时循环读取输出流中的一行内容
  contents := list.New()
  for {
    line, err2 := reader.ReadString('\n')
    if err2 != nil || io.EOF == err2 {
      break
    }
    contents.PushBack(line)
  }
  _ = cmd.Wait()
  return contents, nil
}

/**
  正则简单校验cidr
*/
func RegexpExecute(regexpStr string) bool {
  const reg = `\d{0,3}\.\d{0,3}\.\d{0,3}\.\d{0,3}\/\d{0,3}`
  if match, _ := regexp.MatchString(reg, regexpStr); match {
    return true
  }
  return false
}
