package v1Branch

import (
  "encoding/json"
  "fmt"
  "os"
  "os/exec"
  "strings"

  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/model/pkg/network"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

func Zipmodel(modelpath string, value map[string]string, rawname string) error {
  err := os.Chdir(value["modelsavepath"])
  if err != nil {
    return err
  }
  cmd := exec.Command("unzip", "-o", modelpath, "-d", modelpath[:len(modelpath)-4])
  _, err = utils.CommandExecuteLogs(cmd)
  if err != nil {
    return err
  }
  if _, err = os.Stat("/data/go/model"); err == nil {
    err = os.Remove("/data/go/model")
    if err != nil {
      return err
    }
  }
  fr, err := os.Create("/data/go/model")
  if err != nil {
    return err
  }
  defer fr.Close()
  _, err = fr.WriteString(modelpath[:len(modelpath)-4] + "/" + rawname[:len(rawname)-4])
  if err != nil {
    return err
  }
  return nil
}

func Getdateset(body interfaces.Modeldate, value map[string]string) (string, error) {
  bodyJson, _ := json.Marshal(body)
  res, err := network.DatasetHandler(bodyJson, value["downloadapi"])
  if err != nil {
    return "", err
  }
  log.Info("the res::", res)
  dat := map[string]map[string]string{}
  err = json.Unmarshal([]byte(res), &dat)
  log.Info("path::", dat["data"]["path"])
  path := value["downloadpath"] + dat["data"]["path"]
  cmd := exec.Command("wget", "-P", "/data/go", path)
  _, err = utils.CommandExecuteLogs(cmd)
  if err != nil {
    return "", err
  }
  return path, nil
}

func UnzipDateset(path string, value map[string]string) error {
  s := strings.Split(path, "/")
  zipname := s[len(s)-1]

  relpath := value["handermodelpath"]
  err := os.Chdir(relpath)
  if err != nil {
    return err
  }
  cmd := exec.Command("unzip", "-o", zipname)
  _, err = utils.CommandExecuteLogs(cmd)
  if err != nil {
    return err
  }

  relname := "/data/go/" + zipname[:len(zipname)-4]
  fmt.Println("datesetpath::", relname)

  if _, err = os.Stat("/data/go/dateset"); err == nil {
    err = os.Remove("/data/go/dateset")
    if err != nil {
      return err
    }
  }
  f, err := os.Create("/data/go/dateset")
  if err != nil {
    return err
  }
  defer f.Close()
  _, err = f.WriteString(relname)
  if err != nil {
    return err
  }
  return nil
}

func Handler(value map[string]string) (string, interfaces.Performance, error) {
  err := os.Chdir(value["detect_rate"])
  if err != nil {
    return "", interfaces.Performance{}, err
  }
  log.Info("path::", value["detect_rate"])
  cmd := exec.Command("python3", "test.py")
  _, err = utils.CommandExecuteLogs(cmd)
  if err != nil {
    log.Error("execute evalutate fial::", err)
    return "", interfaces.Performance{}, err
  }
  configFile, err := os.Open("results.json")
  if err != nil {
    log.Error(err)
    return "", interfaces.Performance{}, err
  }

  jsonParser := json.NewDecoder(configFile)
  var test interfaces.Performance
  if err = jsonParser.Decode(&test); err != nil {
    return "", interfaces.Performance{}, err
  }
  date, err := json.Marshal(test)
  if err != nil {
    log.Error(err)
    return "", interfaces.Performance{}, err
  }
  return string(date), test, nil
}
