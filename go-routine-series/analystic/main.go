package main 


import(
	"fmt"
	"flag"
	"github.com/sirupsen/logrus"
	// "golang.org/x/net/html/atom"
	// "crypto/dsa"
	// "net"
	"github.com/mgutz/str"
)
// 1 消费日志 -> chan 
// 2 解析日志  <- chan  可以定义参数
// 3 UV/PV 统计  <- chan 
// 4 数据存储    chan


const HANDLE_DIG = " /dig?"

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
}

type urlNode struct {

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

	//日志消费
	go readFileByLine(params , logChannel)
	//创建一组日志处理
	for i:=0;  i < params.routineNum; i++ {
		go logConsumer (logChannel, pvChannel, uvChannel)
	}
	//创建pv uv 统计器
	go pvCounter( pvChannel, storageChannel)
	go uvCounter( uvChannel, storageChannel)
	//可扩展的xx Counter
	// 创建存储器
	go dataStorage( storageChannel)

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
			uData  :=  urlData{ data ,  uid }

			pvChannel <- uData
			uvChannel <- uData
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
			url : data.Get("url"),
			refer : data.Get("refer"),
			ua : data.Get("ua"),
		}
}

func pvCounter ( pvChannel chan urlData, storageChannel chan storageBlock	){

}

func uvCounter ( uvChannel chan urlData , storageChannel chan storageBlock){

}

func dataStorage (storageChannel chan storageBlock){

}
