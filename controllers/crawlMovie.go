package controllers

import (
	"crawl_movie/models"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
	"net/url"
	"time"
	"fmt"
)

type CrawlMovieController struct {
	beego.Controller
}

/**
目前这个爬虫只能爬取静态数据 对于像京东的部分动态数据 无法爬取
对于动态数据 可以采用 一个组件 phantomjs
*/
func (c *CrawlMovieController) CrawlMovie() {

	//连接到redis
	models.ConnectRedis("212.64.16.41:6379")

	//爬虫入口url
	sUrl := "https://movie.douban.com/subject/27133303/"
	models.PutinQueue(sUrl)

	//爬ip
	 //models.CrawIP()

	for {
		len := models.GetQueueIPLength()
		if len == 0 {
			time.Sleep(100)
			fmt.Print("ip为空")
			continue
		}

		ip := models.PopfromQueueIP()
		if ip == ""{
			fmt.Println("ip拿不到")
			continue
		}
		fmt.Println("拿到ip"+ip)
		for i:=0;i<5;i++{
			length := models.GetQueueLength()
			if length == 0 {
				fmt.Println("url为空")
				models.PutinQueue(sUrl)
				continue //如果url队列为空 则退出当前循环
			}
			sUrl = models.PopfromQueue()
			//我们应当判断sUrl是否应该被访问过
			//if models.IsVisit(sUrl) {
			//	continue
			//}
			//fmt.Println(sUrl)
			logs.SetLogger(sUrl)
			var ipurl *url.URL
			urlproxy,err := ipurl.Parse("https://"+ip)
			if err!=nil{
				//panic(err)
				fmt.Print("解析出错"+ip)
				break
			}
			 models.Run(sUrl,urlproxy)
		}


	}
}


