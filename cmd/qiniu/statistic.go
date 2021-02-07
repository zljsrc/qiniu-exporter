package qiniu

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type SimpleInfo struct {
	Times []int64 `json:"times"`
	Datas []int64 `json:"datas"`
}

func QiniuSampleStatistic(statisticName string, bucket string, bucketManager *storage.BucketManager) (r int64, err error) {
	currentTime := time.Now()
	reqHost := bucketManager.Cfg.ApiReqHost()
	var result SimpleInfo
	reqURL := fmt.Sprintf("%s/v6/%s?bucket=%s&begin=%s&end=%s&g=5min", reqHost, statisticName, bucket, currentTime.Format("20060102150400"), currentTime.Add(time.Minute*5).Format("20060102150400"))
	err = bucketManager.Client.CredentialedCall(context.Background(), bucketManager.Mac, auth.TokenQiniu, &result, "POST", reqURL, nil)
	r = result.Datas[0]
	return
}
