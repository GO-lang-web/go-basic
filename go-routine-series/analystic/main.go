package main 


import(
	"fmt"
	"flag"
	"github.com/sirupsen/logrus"
	// "golang.org/x/net/html/atom"
	// "crypto/dsa"
	// "net"
	"github.com/mgutz/str"
	"github.com/mediocregopher/radix.v2/pool" //redis
	"strconv"
)
// 1 消费日志 -> chan 
// 2 解析日志  <- chan  可以定义参数
// 3 UV/PV 统计  <- chan 
// 4 数据存储    chan


//全局建立redis 链接


const HANDLE_DIG = " /dig?"
cosnt HANDLE_MOVIE = "/movie/"
cosnt HANDLE_LIST = "/list/"
cosnt HANDLE_HTML = ".html" //end

type cmdParams struct {
	logFilePath string 
	routineNum  int 
}

type digData struct {
	time    string 
	url     string
	refer   string 
	ua      string 
}

type urlData struct {
	data    digData
	uid     string 
	unode   urlNode
}


//存储方式对应不同的数据结构
type urlNode struct {
	//urlNodeType
	unType  string  //main  list detail 
	//urlNodeResourceId
	unRid   int
	unUrl   string  //page url 
	unTime  string  // visit page time 
}

type storageBlock struct {
	counterType     string //pv or uv 
	storageModel    string 
	unode           urlNode
}

var log = logrus.New()
func init(){
	log.Out = os.Stdout
	log.SetLevel( logrus.DebugLevel)
}

func main(){
	//获取参数
	logFilePath := flag.String("logFilePath", "/usr/xx.log")
	routineNum := flag.Int("routineNum" , 5 , "consumer number")
	l := flag.String("l", "/tmp/log","log target file path")
	flag.Parse()

	params := cmdParams{
		*logFilePath,
		*routineNum
	}
	//打印日志、

	logFd , err =  os.OpenFile( *l , os.O_CREATE | os.O_WRONLY, 0644)
	if err == nil {
		log.Out = logFd 
		defer logFd.Close()
	}

	log.Infof("start log ")
	log.Infof("Params: logFilePath=%s, routineNum=%d", params.logFilePath , routineNum)

	//初始化一些channel 用于数据传递
	var  logChannel = make(chan string ,3 * params.routineNum ) //传入的数据很重要 所以调整一下
	var  pvChannel  = make(chan urlData, params.routineNum)
	var  uvChannel  = make(chan urlData,params.routineNum)
	var  storageChannel  = make(chan storageBlock, params.routineNum)

	//连接池 redis pool 
	redisPool, err := pool.New("tcp", "localhost:6379",  2* params.routineNum)\
	if err != nil {
		log.Fotalln("redis pool created failed")
		panic(err)
	}else{
		//odle 
		go func() {
			//起一个 goroutine 一直ping
			for {
				redisPool.Cmd("PING")
				time.Sleep(3 * time.Second)
			}
		}()
	}

	//日志消费
	go readFileByLine(params , logChannel)
	//创建一组日志处理
	for i:=0;  i < params.routineNum; i++ {
		go logConsumer (logChannel, pvChannel, uvChannel)
	}
	//创建pv uv 统计器
	go pvCounter( pvChannel, storageChannel)
	go uvCounter( uvChannel, storageChannel, redisPool)
	//可扩展的xx Counter
	// 创建存储器
	go dataStorage( storageChannel, redisPool)

	time.Sleep( 1000 * time.Second)
}

func readFileByLine (params cmdParams , logChannel chan string ){
		fd , err = os.Open(params.logFilePath)
		if err != nil {
			log.Warningf("readFileByLine : can't open file %s", params.logFilePath)
			return err
		}

		defer fd.Close()

		bufferRead := buffer.NewReader( fd)
		
		count := 0 
		for {
			line, err := bufferRead.ReadString('\n') //'' 

			logChannel <- line
			count++

			if count %( 1000 * params.routineNum)  == 0 {
				log.Infof( "readFileByLine: line %d", count)
			}

			if err != nil {
				if err == io.EOF {
					time.Sleep( 3 * time.Second)
					log.Infof("readFileByLine: wait readline %d", count)
				}else{
					log.Warningf("readFileByLine read log error " )
				}
			}
		}

		return nil
}

func logConsumer( logChannel chan string , pvChannel, uvChannel chan urlData){
		for logStr := range logChannel {
			//切割字符串， 取出上报的数据
			data := cutLogUploadData( logStr)

			//uid  一般服务端生成 给cookie  类比百度的网站
			//自己生成  模拟 md5(refer +ua)
			hasher := m5.New()
			hasher.Write([]byte(data.refer + data.ua ))
			uid := hex.EncodeToString( hasher.Sum( nil))

			//many parse job  can added here 
			uData  :=  urlData{ data ,  uid， formatUrl(data.url , data.time ) }

			pvChannel <- uData
			uvChannel <- uData
		}
}

func formatUrl(url ,t string ) urlNode{
	 //大量的页面入手 detail页面 >= home list page 
	 pos1 := str.IndexOf(url , HANDLE_MOVIE , 0 )
	 if pos1 != -1 { //detail  page 
			pos1 += len(HANDLE_MOVIE)
			pos2 := str.IndexOf(url, HANDLE_HTML, 0)

			idStr := str.Substr(url, po1, pos2 -pos1 )
			id := strconv.Atoi(idStr) //str => number
			return urlNode{"movie", id, url ,t}
	 }else{ //list page 
		 pos1 =  str.IndexOf(url , HANDLE_LIST , 0 )
		 if pos1 != -1 {
			 pos1 += len( HANDLE_LIST)
			 pos2 := str.IndexOf(url, HANDLE_HTML, 0)

			 idStr := str.Substr(url, po1, pos2 -pos1 )
			 id := strconv.Atoi(idStr) //str => number
			 return urlNode{"list", id, url ,t}
		 }else{
			 //home page    redis 0 invalid so change to 1 
			 return urlNode("home", 1 , url , t )
		 }
	 }

}

func cutLogUploadData ( logStr string ) digData {
		logStr = strings.TrimSpace(logStr)
		//截取字符串 常规操作
		pos1 := str.IndexOf(logStr , HANDLE_DIG , 0 )
		if pos1 == -1 {
			return digData{}
		}
		//找到了
		pos1 += len(HANDLE_DIG)
		pos2 = str.IndexOf(logStr, "HTTP/" , pos1)

		d := str.Substr(logStr , pos1, pos2 - pos1)

		//key val 的形式 需要url 处理一下
		urlInfo , err := url.Parse("http://localhost/?" + d)
		if err != nil {
			return  digData{}
		}
		data := urlInfo.Query()

		return digData{
			time : data.Get("time"),
			url :  data.Get("url"),
			refer : data.Get("refer"),
			ua : data.Get("ua"),
		}
}

func pvCounter ( pvChannel chan urlData, storageChannel chan storageBlock	){
	 for data := range pvChannel {
		 sItem := storageBlock{ "pv" , "ZINCRBY" , data.unode}

		 storageChannel <- sItem
	 }
}

func getTime (logTime , timeType string ) string{
	 var item string 
	 switch timeType {
	 	case "day" : 
			item ="2006-01-02"
			break
	 	case "hour" :
			item ="2006-01-02 15"
			break 
		case "min" :
			item ="2006-01-02 15:04"
			break 
	 }
	 t , _ := time.Parse(item , time.Now().Format(item))
	 return strconv.FormatInt( t.Unix(), 10 )
}

func uvCounter ( uvChannel chan urlData , storageChannel chan storageBlock, redisPool *pool.Pool){
		//uv 要去重
		for data := range uvChannel {
			//hyperLoglog  redis 
			hyperLoglogKey := "uv_hpll" + getTime(data.data.time, "day")

			// EX 设置过期时间
			ret, err := redisPool.Cmd("PFADD", hyperLoglogKey, data.uid, "EX", 86400).Int()
			if err != nil{
				log.Warningln("uvCounter : check redis hyperloglog failed")
				//查失败 相当于没有 不return 
			}
			if ret != 1 {
				continue
			}

			sItem := storageBlock{ "uv" , "ZINCRBY" , data.unode}

			storageChannel <- sItem
		}
}


//企业级 使用HBase 劣势： 列簇需要声明清楚
func dataStorage (storageChannel chan storageBlock, redisPool *pool.Pool){
		for  block := range  storageChannel {
			redis_prefix := block.counterType + '_'

			//逐层增加， 
			//movie/1.html  
			//detail + pv1 
			//cate +  pv1
			//home + pv1 
			//网页- 大分类- 小分类- 详情页面
			//维度 天-小时-分钟
			//存储模型 redis SortedSet
			setKeys := [] string {
				//天-小时-分钟
				redis_prefix + "day_" + getTime(block.unode.unTime, "day")
				redis_prefix + "hour_" + getTime(block.unode.unTime, "hour")
				redis_prefix + "min_" + getTime(block.unode.unTime, "min")
				//网页- 大分类- 小分类- 详情页面
				redis_prefix + block.unode.unType + "_day_" + getTime(block.unode.unTime, "day")
				redis_prefix + block.unode.unType + "_hour_" + getTime(block.unode.unTime, "hour")
				redis_prefix + block.unode.unType + "_min_" + getTime(block.unode.unTime, "min")

			}

			rowId := block.unode.uid

			for _,key := range setKeys {
				//set to redis 
				ret, err := redisPool.Cmd(block.storageModel, key, 1, rowId)
				if err != nil{
					//提供报错信息尽可能完整 统计学里面少数错误可以忽略
					log.Errorln("dataStorage: resis storage error",block.storageModel, key, rowId)
				}else{
					//success 为什么不打日志 数量很大
					// 10000 * N* M uv /pv 后面还可能加别的

				}
			}
		}
}
