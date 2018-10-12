package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strings"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	_ "net/url"
	"net/url"
	"fmt"
	"time"
	"math/rand"
)

var (
	db orm.Ormer
)

type MovieInfo struct{ 
  Id int64
  Movie_id int64
  Movie_name string
  Movie_pic string
  Movie_director string
  Movie_writer string
  Movie_country string
  Movie_language string
  Movie_main_character string
  Movie_type string
  Movie_on_time string
  Movie_span string
  Movie_grade string
}

func init() {
	//orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	//orm.RegisterDataBase("default", "mysql", "root:ww0819@tcp(127.0.0.1:3306)/data?charset=utf8", 30)
	//orm.RegisterModel(new(MovieInfo))
	//db = orm.NewOrm()
}

func AddMovie(movie_info *MovieInfo)(int64,error){
	movie_info.Id = 0
	id,err := db.Insert(movie_info)
	logs.Error(err)
	return id,err
}

//电影导演
func GetMovieDirector(movieHtml string) string{
	if movieHtml == ""{
		return ""
	}
	reg := regexp.MustCompile(`<a.*?rel="v:directedBy">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0{
		return ""
	}
	logs.SetLogger(result[0][1])
	return string(result[0][1])
}
//电影名字
func GetMovieName(movieHtml string)string{
	if movieHtml == ""{
		return ""
	}

	reg := regexp.MustCompile(`<span\s*property="v:itemreviewed">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0{
		return ""
	}

	return string(result[0][1])
}
//电影主演
func GetMovieMainCharacters(movieHtml string)string{
	reg := regexp.MustCompile(`<a.*?rel="v:starring">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0{
		return ""
	}

	mainCharacters := ""
	for _,v := range result{
		mainCharacters += v[1] + "/"
	}

	return strings.Trim(mainCharacters, "/")
}
//电影评分
func GetMovieGrade(movieHtml string)string{
	reg := regexp.MustCompile(`<strong.*?property="v:average">(.*?)</strong>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0{
		return ""
	}
	return string(result[0][1])
}
//电影类型
func GetMovieGenre(movieHtml string)string{
	reg := regexp.MustCompile(`<span.*?property="v:genre">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0{
		return ""
	}

	movieGenre := ""
	for _,v := range result{
		movieGenre += v[1] + "/"
	}
	return strings.Trim(movieGenre, "/")
}
//上映时间
func GetMovieOnTime(movieHtml string) string{
	reg := regexp.MustCompile(`<span.*?property="v:initialReleaseDate".*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0{
		return ""
	}

	return string(result[0][1])
}
//电影时长
func GetMovieRunningTime(movieHtml string) string{
	reg := regexp.MustCompile(`<span.*?property="v:runtime".*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0{
		return ""
	}

	return string(result[0][1])
}


func GetMovieUrls(movieHtml string)[]string{
	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/.*?)"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	var movieSets []string
	for _,v := range result{
		movieSets = append(movieSets, v[1])
	}

	return movieSets
}

func Run(sUrl string,urlproxy *url.URL)  {
	fmt.Print("开始爬页面")


	request, _ := http.NewRequest("GET", sUrl, nil)
	//随机返回User-Agent 信息
	request.Header.Set("User-Agent", getAgent())
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	//设置超时时间
	timeout := time.Duration(10* time.Second)
	client := &http.Client{}

		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
			Timeout: timeout,
		}

	rsp, err := client.Do(request)
	if err!=nil{
		fmt.Println("获取页面内容出错"+sUrl)
		//PutinQueue(sUrl)
		return
	}
	if rsp!= nil{
		fmt.Println("记录电影信息111111111111111111111111111111111111111111111111")
		var movieInfo MovieInfo
		defer rsp.Body.Close()
		body, _ := ioutil.ReadAll(rsp.Body)
		//fmt.Println(body)
		sMovieHtml := string(body)
		movieInfo.Movie_name            = GetMovieName(sMovieHtml)
		fmt.Print(sMovieHtml)
		//记录电影信息
		if movieInfo.Movie_name != ""{
			movieInfo.Movie_director        = GetMovieDirector(sMovieHtml)
			movieInfo.Movie_main_character  = GetMovieMainCharacters(sMovieHtml)
			movieInfo.Movie_type            = GetMovieGenre(sMovieHtml)
			movieInfo.Movie_on_time         = GetMovieOnTime(sMovieHtml)
			movieInfo.Movie_grade           = GetMovieGrade(sMovieHtml)
			movieInfo.Movie_span            = GetMovieRunningTime(sMovieHtml)

			//AddMovie(&movieInfo)
			fmt.Println(movieInfo)
		}

		//提取该页面的所有连接
		urls := GetMovieUrls(sMovieHtml)
		//urls := GetMovieUrls1(sMovieHtml)

		for _,url := range urls{
			PutinQueue(url)
			//c.Ctx.WriteString("<br>" + url + "</br>")
		}

		//sUrl 应当记录到 访问set中
		AddToSet(sUrl)

		//time.Sleep(time.Second)
	}




}

func GetMovieUrls1(movieHtml string)[]string{
	reg := regexp.MustCompile(`<a href="(.*?)"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	var movieSets []string
	for _,v := range result{
		movieSets = append(movieSets, v[1])
		fmt.Print(v[1])
	}

	return movieSets
}


func GetRep(urls string) *http.Response {
	ip := string(GETIP())
	fmt.Print(ip)
	request, _ := http.NewRequest("GET", urls, nil)
	//随机返回User-Agent 信息
	request.Header.Set("User-Agent", getAgent())
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	proxy, err := url.Parse(ip)
	//设置超时时间
	timeout := time.Duration(20* time.Second)
	fmt.Printf("使用代理:%s\n",proxy)
	client := &http.Client{}
	if ip != "local"{
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
			Timeout: timeout,
		}
	}

	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200{

		fmt.Printf("line-99:遇到了错误-并切换ip %s\n",err)
		GETIP()
		GetRep(urls)
		//getIp(returnIP())

	}

	return response
}


/**
* 随机返回一个User-Agent
*/
func getAgent() string {
	agent  := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}