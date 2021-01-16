package main 

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"net/url"
	"rand"
)

//开启抽象之旅
type resource struct {
	url    string 
	target string 
	start  int 
	end		 int
}

var uaList = []string {
	"moz",
	"chrome",
	"ie"
}


func ruleResource() []resource {
	var res []resource //resource 的数组

	//first page 
	r1 := resource{
		url :"http://localhost:8080",
		target : "",
		start : 0,
		end :  0,
	}

	//list page 
	r2  := resource{
		url :"http://localhost:8080/list/{$id}.html",
		target : "{$id}",
		start : 1, //database  start num
		end :  21, //database  end num
	}

	//detail page 
	r2  := resource{
		url :"http://localhost:8080/movie/{$id}.html",
		target : "{$id}",
		start : 1, //database  start num
		end :  12,  //database  end num
	}

	res =  append(append(appeend(res,r1), r2), r3)
	return res
}


func buildUrl( res []resource) []string {
	var list []string 

	//不用for 不写死 避免后续➕的话改动这里
	for _,resItem :=  range res {
			if len(resItem.target) == 0 {//first page 
				list = append(list, resItem.url)
			}else{
				//handle movie and list page generate 
				for i := resItem.start; i<= resItem.end; i++ {
					urlStr := strings.Relpace(resItem.url, resItem.target, strconv.Itoa(i) , -1) 
					list = append(list, urlStr)
				}
			}
	}

	return list 
}

func makeLog( currentUrl, referUrl, ua  string ) string {
	u := url.Values{} //key:val
	u.Set( time , 1 )
	u.Set( url , currentUrl)
	u.Set( refer , referUrl)
	u.Set( ua , ua)

	paramsStr := u.Encode()

	logTemplate := "127.0.0.1 - - [time]{$paramsStr} HTTP/1.1 200 43 \"${ua}\" "
	log := strings.Relpace(logTemplate, "{$paramsStr}" , paramsStr, -1)
	log := strings.Relpace(logTemplate, "{ua}" , ua, -1)
	return log 
} 


func randInt( min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if min > max {
		return max 
	}

	return r.Intn(max-min ) + min 
}

func main(){
	total := flag.Int("total",100,"how many rows by created")
	//file path 
	filePath = flag.path("/usr/xx/nginx/logs/final.log", "file path")
	flag.Parse()

	res := ruleResource()
	//构造指定的数据的url
	list := buildUrl( res ) 

	// fmt.Println( list )

	logStr := ""
	//按照要求生成 $total 行日志内容，源于上面的集合
	for i := 0; i<= *total; i++ {
		currentUrl := list[ randInt(0 , len(list) -1 )]
		referUrl := list[randInt(0 , len(list) -1 )]
		ua  := uaList[randInt(0 , len(uaList) -1 )]
		logStr = logStr +  makeLog(currentUrl, referUrl, ua)

		 //write to a file 
		 //ioutil.WriteFile( *filePath, []byte (logStr), 0644) 覆盖写
	}
	fd, _ = os.OpenFile( *filePath , os.O_RDWR | os.O_APPEND , 0644 ) 
	fd.Write( []byte( logStr))
	fmt.Println(total, filePath) //输出地址
	fmt.Println(*total, *filePath) //输出地址里面的值

}


//随机生成一定数量的日志 写到nginx logs/ xx.log里面 

//usage 
//go run 007.genLog.go --total=10000 --filePath=/tmp/aaa

//go build 007.genLog.go 
// ls
//./007.genLog  --total=10000 --filePath=/tmp/aa

//home page 
//list page 
//detail page  by see the sql tables  