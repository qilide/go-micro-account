package micro

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

func PrometheusBoot(host string,port int){
	http.Handle("/metrics",promhttp.Handler())
	//启动web服务
	go func() {
		err := http.ListenAndServe(host+":"+strconv.Itoa(port),nil)
		if err!= nil{
			log.Fatal(("监控启动失败"))
		}
		log.Fatal("监控启动,端口为: "+strconv.Itoa(port))
	}()
}
