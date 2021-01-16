# go-basic

## editor

- goLand
- vscode

## goroutine

> 如何开启？

- 1. go + 匿名函数

```go
message := make(chan string)
```

- 2. go + 函数名

```go
var  message = make(chan string)
func test(){
	message <- "hello goroutine"
}
func test2(){
	time.Sleep(2* time.Second)
	anotherStr := <-message
	anotherStr = anotherStr + " another"
	message<-anotherStr
}

func main() {
	go test()
	go test2()
	time.Sleep(3*time.Second)
	fmt.Println(<-message)
	fmt.Println("test over")
}s
```

> 如何输入输出？

- 声明 chan
- chan <- 往 chan 里写
- <- chan 从 chan 里面输出
- 多个的话 指定 chan 长度

## select

- 5
- 6 需要再看一下

## wirte a js function to upload some user information

- send by a ajax get method
- url
- referer
- time
- ip can get easy from backend or nginx
- cookie

## NGINX

- EMPTY_GIF MODULE
- ACCESS_LOG [nginx access_log 日志](https://lanjingling.github.io/2016/03/14/nginx-access-log/)

```
//conf/nginx.conf

log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
								'"$status" $body_bytes_sent "$http_referer" '
								'"$http_user_agent" "$http_x_forwarded_for" '
								'"$gzip_ratio" $request_time $bytes_sent $request_length';
需要用access_log指令指定日志文件的存放路径；
access_log : logs/logFileName.log main;


设置刷盘策略：
access_log /data/logs/nginx-access.log buffer=32k flush=5s;

buffer 满 32k 才刷盘；假如 buffer 不满 5s 钟强制刷盘。

Nginx 的日志都是写在一个文件当中的，不会自动地进行切割，如果访问量很大的话，将导致日志文件容量非常大，不便于管理和造成Nginx 日志写入效率低下等问题。所以，往往需要要对access_log、error_log日志进行切割。

切割日志一般利用USR1信号让nginx产生新的日志。实例：


#!/bin/bash

logdir="/data/logs/nginx"
pid=`cat $logdir/nginx.pid`
DATE=`date -d "1 hours ago" +%Y%m%d%H`
DATE_OLD=`date -d "7 days ago" +%Y%m%d`
for i in `ls $logdir/*access.log`; do
        mv $i $i.$DATE
done
for i in `ls $logdir/*error.log`; do
        mv $i $i.$DATE
done
kill -s USR1 $pid
rm -v $logdir"/access.log."$DATE_OLD*

rm -v $logdir"/error.log."$DATE_OLD*


将上面的脚本放到crontab中，每小时执行一次（0 ），这样每小时会把当前日志重命名成一个新文件；然后发送USR1这个信号让Nginx 重新生成一个新的日志。（相当于备份日志）
将前7天的日志删除；
说明：
在没有执行kill -USR1 $pid之前，即便已经对文件执行了mv命令而改变了文件名称，nginx还是会向新命名的文件”*access.log.2016032623”照常写入日志数据的。原因在于：linux系统中，内核是根据文件描述符来找文件的。

使用系统自带的logrotates，也可以实现nginx的日志分割，查看其bash源码，发现也是发送USR1这个信号。
```

- reload nginx `./sbin.nginx -s reload`

- 查看 logs `tail -f logs/logFileName.log` 滚动输出。 前端触发一下
- `wc -l logFileName.log` 查看文件行数
- `gizp on` 几百 kb 的时候 效果非常好 200kb -> 几十 kb

## 解耦 日志与打点服务器

- user -> open website -> website embed js upload sdk -> every time user open a page -> js will invoked the information to the log server (nginx) -> save the logs to a named file -> analystic -> show data in frontend

- ls -a 显示.文件
- .htaccess 是 rewrite 的规则

- 清空 log echo '' > /tmp/log
- 查看 tail -f /tmp/log

## pv uv vv

- UV unique visitor : 访问您电脑的一台客户端为一个访客，00:00 - 24:00 相同的客户端只被计算一次
- IP Internet Protocol 访问过某站点的 IP 总数， 以用户的 IP 地址作为统计依据 00:00 - 24:00 相同的客户端只被计算一次
  > UV 和 IP 的区别
- a b 各自账号在同一台电脑上登录 weibo, IP 数 +1 ， UV +2 ,因为使用的是一台电脑

- PV page view 页面浏览量或者用户点击量 用户每一次对网站中的每个网页访问均被记录为一个 PV,用户对同一个网页多次访问，PV 累计增加

- VV visit view 统计访客 一天内访问网站的次数 完成浏览并且关闭网页就完成一次 一天内可能有多次访问行为

> 你今天 10 点钟打开 google.com 访问三个页面， 12 点又打开 google.com 访问两个页面
> PV +5 VV+2 pv 指的是页面的浏览次数 vv 指的是你访问网站的次数

## 并发模型里面 routine 里面使用外部的变量 比如 redis 应该是连接池

## go

- 完整项目 不到 300 行代码
- 适合做流水线，扛住压力的情况

## 可视化 TODO

- `yarn create umi showData`
- cd `showData`
- `npm i`
- `npm start`

## 提供两个接口 TODO

- 1. getPvUv
- 2. getRank top=200

## 回顾

- mamp 搭建网站 -（js 上报数据）-> nginx(access_log 收集并且存储 log) -> go 处理监听新增数据 协程处理 -(存到 redis) -> (api 从 redis 取数据提供接口) -> (ant design pro 脚手架搭建项目， 前端 ajax 请求 pv uv top200 show in frontend)

## 企业级收集看重：

> 数据收集侧

- 1. user information can be as detail as possible
- 2. app 的数据收集
- 3. 尽可能避免丢失 数据上报失败或者关键环节忘记加上
     > 数据处理侧
- 1. 准确性
- 2. 稳定性加强 多维度 不仅仅依赖程序自动启动
- 3. 存储方案的迭代 不同发展阶段对应数据存储压力不同，提前做好储备

## final

> go

- 1. go routine 以及控制
- 2. go routine 之间的通信 channel
- 3. go rountne 组成的流水线

> project

- 1. nginx 打点方案
- 2. 流量统计服务
- 3. 数据统计服务
- 4. ant Desgin Pro
