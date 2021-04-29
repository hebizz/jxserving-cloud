package database

import (
  "os"
  "strings"

  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/portal/database"
  log "k8s.io/klog"
)

var db database.MongoDatabase

func init() {
  rdb, _ := database.NewDatabase(database.TyMongo)
  dbClient := rdb.(database.MongoDatabase)
  err := dbClient.NewConnection()
  dbClient.UpdateDatabase("jxsCloud")
  if err != nil {
    panic(err)
  }
  db = dbClient
  InitAssistantBackendPortal()
}

/**
  初始化远程backend的信息
*/
func InitAssistantBackendPortal() {
  port := os.Getenv("PORTS")
  host := os.Getenv("HOSTS")
  api := os.Getenv("APIS")
  project := os.Getenv("PROJECTS")
  if port == "" {
    log.Info("Please set env args in formal environment")
    return
  }
  log.Info("assistant env host & port & apis & projects:", host, port, api, project)
  // 存储到数据库
  ports := strings.Split(port, ",")
  hosts := strings.Split(host, ",")
  apis := strings.Split(api, ",")
  projects := strings.Split(project, ",")

  for i := 0; i < len(ports); i++ {
    err, proJ := QueryProjectInfoByName(projects[i])
    if err != nil {
      log.Info("query project info error:", err)
    }
    if proJ.Name == "" {
      err := InsertProjectInfo(interfaces.Project{Name: projects[i], Host: hosts[i], Port: ports[i], Api: apis[i]})
      if err != nil {
        panic(err)
      }
    }
  }
}
