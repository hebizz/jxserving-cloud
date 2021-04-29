package database

import (
  "errors"
  "os"

  influx "github.com/influxdata/influxdb1-client/v2"
  "github.com/spf13/cast"
  "gitlab.jiangxingai.com/jxserving/components/interfaces"
  "gitlab.jiangxingai.com/jxserving/components/utils"
)

type InfluxDB struct {
  Host        string
  Username    string
  Password    string
  Database    string
  TableName   string
  client      influx.Client
  batchPoints influx.BatchPoints
}

func newInfluxDB() *InfluxDB {
  if ret := utils.ReadMap(ConfigEntryKey); ret != nil {
    value := cast.ToStringMapString(ret[TyInflux])
    return &InfluxDB{
      Host:      value["host"],
      Database:  value["database"],
      TableName: value["tablename"],
      Username:  os.Getenv("JXS_ANA_INFLUX_USR"),
      Password:  os.Getenv("JXS_ANA_INFLUX_PWD"),
      client:    nil,
    }
  } else {
    panic("unable to find InfluxDB config")
  }
}

func (db *InfluxDB) NewConnection() error {
  conn, err := influx.NewHTTPClient(influx.HTTPConfig{
    Addr: db.Host,
  })
  if err != nil {
    return err
  }
  db.client = conn

  db.batchPoints, err = influx.NewBatchPoints(influx.BatchPointsConfig{
    Precision: "s",
    Database:  db.Database,
  })
  if err != nil {
    return err
  }
  return nil
}

func (db *InfluxDB) UpdateDatabase(database string) {
  db.Database = database
}

func (db *InfluxDB) UpdateTable(table string) {
  db.TableName = table
}

func (db *InfluxDB) Insert(i interface{}) error {
  in := i.(interfaces.AnalysisData)
  point, err := in.ToInfluxPoint(db.TableName)
  if err != nil {
    return err
  }
  db.batchPoints.AddPoint(point)
  return nil
}

func (db *InfluxDB) Commit() error {
  return db.client.Write(db.batchPoints)
}

func (db *InfluxDB) Query(filter interface{}) (interface{}, error) {
  f, ok := filter.(map[string]string)
  if !ok {
    return nil, errors.New("failed to accept filter")
  }

  query := influx.Query{
    Command:  f["command"],
    Database: db.Database,
  }

  if result, err := db.client.Query(query); err == nil {
    if result.Error() == nil {
      return transInfluxDBResultToMap(result.Results), nil
    } else {
      return nil, result.Error()
    }
  } else {
    return nil, err
  }
}

func transInfluxDBResultToMap(result []influx.Result) map[string]interface{} {
  return map[string]interface{}{}
}
