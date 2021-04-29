package v1

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/montanaflynn/stats"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"

  "gitlab.jiangxingai.com/jxserving/components/portal/database"
  "gitlab.jiangxingai.com/jxserving/components/portal/response"
)

type Condition struct {
  Key   int8   `json:"key"`
  Logic string `json:"logic"`
  Field string `json:"field"`
  Value string `json:"value"`
}

type Filter struct {
  Timestamp []int64     `json:"timestamp" binding:"required"`
  Value     string      `json:"value" binding:"required"`
  Cond      []Condition `json:"cond" binding:"required"`
}

func (f *Filter) ToBSON() interface{} {
  cond := bson.M{
    "alerttime": bson.M{
      "$gte": f.Timestamp[0],
      "$lte": f.Timestamp[1],
    },
  }
  for _, c := range f.Cond {
    cond[c.Field] = bson.M{"$in": []string{c.Value}}
  }
  return cond
}

func Query(c *gin.Context) {
  var filter Filter

  if err := c.ShouldBindJSON(&filter); err != nil {
    response.ServerError(c, http.StatusInternalServerError, "failed to parse json", nil, "")
    return
  }
  mongoDB, err := database.NewDatabase(database.TyMongo)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  mongoDB.UpdateTable("reportData")
  if err := mongoDB.NewConnection(); err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  rawData, err := mongoDB.Query(filter.ToBSON())
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  rawBSON, ok := rawData.([]bson.M)
  if !ok {
    response.ServerError(c, http.StatusInternalServerError, "failed to parse data", nil, "")
    return
  }
  retData := make(map[string]interface{})
  var calcData []float64
  for _, d := range rawBSON {
    id := d["_id"].(primitive.ObjectID).Hex()
    reliability := d["reliability"].(float64)
    retData[id] = reliability
    calcData = append(calcData, reliability)
  }
  stat, err := statistic(calcData)
  if err != nil {
    response.ServerError(c, http.StatusInternalServerError, "", err, "")
    return
  }
  response.Success(c, "query analyst data success", map[string]interface{}{
    "data": retData,
    "stat": stat,
  })
}

func statistic(data []float64) (map[string]float64, error) {
  result := map[string]float64{
    "min": 0.0,
    "max": 0.0,
    "ave": 0.0,
    "std": 0.0,
    "cov": 0.0,
    "p95": 0.0,
    "p99": 0.0,
  }

  if min, err := stats.Min(data); err != nil {
    return nil, err
  } else {
    result["min"] = min
  }

  if max, err := stats.Max(data); err != nil {
    return nil, err
  } else {
    result["max"] = max
  }

  if ave, err := stats.Mean(data); err != nil {
    return nil, err
  } else {
    result["ave"] = ave
  }

  if std, err := stats.StandardDeviation(data); err != nil {
    return nil, err
  } else {
    result["std"] = std
  }

  result["cov"] = result["std"] / result["ave"]

  if p95, err := stats.Percentile(data, 95); err != nil {
    return nil, err
  } else {
    result["p95"] = p95
  }

  if p99, err := stats.Percentile(data, 99); err != nil {
    return nil, err
  } else {
    result["p99"] = p99
  }

  return result, nil
}
