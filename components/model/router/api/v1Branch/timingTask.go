package v1Branch

import (
  "os"

  "github.com/robfig/cron"
  "gitlab.jiangxingai.com/jxserving/components/model/pkg/config"
  "gitlab.jiangxingai.com/jxserving/components/utils"
  log "k8s.io/klog"
)

func init() {
  go TimingTask()
}

/**
  定时任务
*/
func TimingTask() {
  cn := cron.New()
  // 每天凌晨2点执行一次
  spec := "0 2 * * *"
  _, err := cn.AddFunc(spec, func() {
    //delete 临时文件
    err := DeleteModelTmpFile()
    if err != nil {
      log.Error("delete model tmp file error:", err)
    }
  })
  if err != nil {
    log.Error("timing task error:", err)
  }
  cn.Start()
  select {}
}

/**
  删除下载模型生成的临时文件
*/
func DeleteModelTmpFile() error {
  dir := config.ModelSavePath()
  tar := dir + "tar"
  var err error
  if utils.PathOrFileExists(tar) {
    err = os.RemoveAll(tar)
  }
  unTar := dir + "unTar"
  if utils.PathOrFileExists(unTar) {
    err = os.RemoveAll(unTar)
  }
  return err
}
