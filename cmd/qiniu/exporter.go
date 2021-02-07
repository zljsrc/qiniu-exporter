package qiniu

import (
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type qiniuMetricsCollector struct {
	qiniuAccessKey string
	qiniuSecretKey string
	spaceMetric    *prometheus.Desc
	countMetric    *prometheus.Desc
}

func QiniuMetricsController(qiniuAccessKey string, qiniuSecrteKey string) *qiniuMetricsCollector {
	return &qiniuMetricsCollector{
		qiniuAccessKey: qiniuAccessKey,
		qiniuSecretKey: qiniuSecrteKey,
		spaceMetric:    prometheus.NewDesc("qiniu_bucket_space", "存储的存储量统计 space:标准存储 space_line:低频存储 space_archive:归档存储", []string{"bucket", "storage_type_key"}, nil),
		countMetric:    prometheus.NewDesc("qiniu_bucket_count", "存储的文件数量统计 count:标准存储 count_line:低频存储 count_archive:归档存储", []string{"bucket", "storage_type_key"}, nil),
	}
}

func (collector *qiniuMetricsCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.spaceMetric
	ch <- collector.countMetric
}

func (collector *qiniuMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	//获取统计数据
	mac := qbox.NewMac(collector.qiniuAccessKey, collector.qiniuSecretKey)
	cfg := storage.Config{UseHTTPS: true}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	bucketsInfo, err := bucketManager.Buckets(true)
	if err != nil {
		log.Fatal("获取bucket列表失败")
		log.Fatal(err)
		return
	}

	var wg = sync.WaitGroup{}
	for _, bucketName := range bucketsInfo {
		wg.Add(6)
		go collector.collectSpaceInfo("count", bucketName, bucketManager, ch, &wg)
		go collector.collectCountInfo("space", bucketName, bucketManager, ch, &wg)
		go collector.collectSpaceInfo("count_line", bucketName, bucketManager, ch, &wg)
		go collector.collectCountInfo("space_line", bucketName, bucketManager, ch, &wg)
		go collector.collectSpaceInfo("count_archive", bucketName, bucketManager, ch, &wg)
		go collector.collectCountInfo("space_archive", bucketName, bucketManager, ch, &wg)
	}

	time.Sleep(time.Second * 10)
	wg.Wait()

	return
}

func (collector *qiniuMetricsCollector) collectSpaceInfo(storageTypeKey string, bucketName string, bm *storage.BucketManager, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	value, err := QiniuSampleStatistic(storageTypeKey, bucketName, bm)
	if err != nil {
		log.Fatal(bucketName, storageTypeKey, " get space error")
		log.Fatal(err)
	}
	ch <- prometheus.MustNewConstMetric(collector.spaceMetric, prometheus.CounterValue, float64(value), []string{bucketName, storageTypeKey}...)
	log.Println(bucketName, storageTypeKey, "space:", value)

	wg.Done()
	return
}

func (collector *qiniuMetricsCollector) collectCountInfo(storageTypeKey string, bucketName string, bm *storage.BucketManager, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	value, err := QiniuSampleStatistic(storageTypeKey, bucketName, bm)
	if err != nil {
		log.Fatal(bucketName, storageTypeKey, " get count error")
		log.Fatal(err)
	}
	ch <- prometheus.MustNewConstMetric(collector.countMetric, prometheus.CounterValue, float64(value), []string{bucketName, storageTypeKey}...)
	log.Println(bucketName, storageTypeKey, "count:", value)

	wg.Done()
	return
}
