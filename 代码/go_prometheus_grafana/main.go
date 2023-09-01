package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//初始化日志服务
	logger := log.New(os.Stdout, "[Memory]", log.Lshortfile|log.Ldate|log.Ltime)

	//初始一个http handler
	http.Handle("/metrics", promhttp.Handler())

	//初始化一个容器
	diskPercent := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "memeory_percent",
		Help: "memeory use percent",
	},
		[]string{"percent"},
	)
	prometheus.MustRegister(diskPercent)

	// 启动web服务，监听1010端口
	go func() {
		logger.Println("ListenAndServe at:localhost:8080")
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			logger.Fatal("ListenAndServe: ", err)
		}
	}()

	//收集内存使用的百分比
	for {
		logger.Println("start collect memory used percent!")
		v := rand.Intn(1000)

		usedPercent := v
		logger.Println("get memeory use percent:", usedPercent)
		diskPercent.WithLabelValues("usedMemory").Set(float64(usedPercent))
		time.Sleep(time.Second * 2)
	}
}
