package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func Lookup(host string) (string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}
	if len(addrs) < 1 {
		return "", errors.New("unknown host")
	}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return addrs[rd.Intn(len(addrs))], nil
}

var Data = []byte("abcdefghijklmnopqrstuvwabcdefghi")

type Reply struct {
	Time  int64
	TTL   uint8
	Error error
}
var (
	Pingerr int
	Pingok int
	debug_req int
	retry int
	TimeOut int
)
//多线程
func IPpingStart(line string,req chan []string,status chan int) {
	ping, err := Run(line, 8, Data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ping.Close()
	respones:=ping.Ping(1)
	if respones!="is alive" {
		for r:=0;r<=retry;r++ {
			respones=ping.Ping(1)
			}
	}

	if debug_req==1 {
		fmt.Println(line+" "+respones)
	}
	<-status
	req<- []string{line,respones}
}
//读文件
func main()  {
	Pingerr=0
	Pingok=0
	//命令行参数
	var ipfile string
	flag.StringVar(&ipfile, "i", "ip.txt", "IP文件所在路径,,默认读取当前目录ip.txt")
	var output string
	flag.StringVar(&output, "o", "out.csv", "结果输出文件所在路径,默认放在当前目录out.csv")
	var num int
	flag.IntVar(&num, "n", 20, "开启的并发协程数量,默认20")
	flag.IntVar(&debug_req, "d", 0, "是否打开debug模式,即开启显示每条IP结果记录")
	flag.IntVar(&retry, "r", 2, "IP检测失败重试次数,默认2次")
	flag.IntVar(&TimeOut, "t", 1, "检测超时时间（单位秒），默认1秒")
	var help bool
	flag.BoolVar(&help,"h", false, "显示此帮助页\nbuild by zhangjikun@haima.me\nversion 1.0")
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	t:=time.Now() //取当前时间戳
	status := make(chan int, num)
	req:=make(chan []string)
	fmt.Println("正在加载IP资源文件....")
	f, err := os.Open(ipfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	LinsNum:=CountFileLine(ipfile)
	ips:=[]string{}

	for lines:=0;lines<=LinsNum;lines++ {
		line, _ := rd.ReadString('\n') //以'\n'为结束符读入一行
		line = strings.TrimSpace(line)
		if len(line)>5 && len(strings.Split(line, "."))==4 {
			ips=append(ips,line)
		}
	}
	LinsNum=len(ips)
	fmt.Println("IP资源加载完成,正在启动协程...")
	for _,ip:=range ips{
		status <- 1
		go IPpingStart(ip,req,status)
	}
	strings_req:=[]string{"ip","status"}
	fd,_:=os.OpenFile(output,os.O_RDWR|os.O_CREATE,0644)
	title:=[]byte(strings_req[0]+","+strings_req[1]+"\n")
	fd.Write(title)
	for i:=0;i<LinsNum;i++ {
		strings_req=<-req
		buf:=[]byte(strings_req[0]+","+strings_req[1]+"\n")
		//判断成功失败数量
		if strings_req[1]=="is alive" {
			Pingok++
		} else {
			Pingerr++
		}
		fd.Write(buf)
	}
	fd.Close()
	fmt.Println("批量ping已经完成!\n运行耗时: ",time.Since(t),"\nping成功主机数: ",Pingok,"\nping失败主机数: ",Pingerr,"\n主机总数: ",LinsNum)
}

//获取文件行数
func CountFileLine(name string) (count int) {
	data, _ := ioutil.ReadFile(name)
	count = 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			count++
		}
	}
	return count
}


func MarshalMsg(req int, data []byte) ([]byte, error) {
	xid, xseq := os.Getpid()&0xffff, req
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: xid, Seq: xseq,
			Data: data,
		},
	}
	return wm.Marshal(nil)
}

type ping struct {
	Addr string
	Conn net.Conn
	Data []byte
}

func (self *ping) Dail() (err error) {
	self.Conn, err = net.Dial("ip4:icmp", self.Addr)
	if err != nil {
		return err
	}
	return nil
}

func (self *ping) SetDeadline(timeout int) error {
	return self.Conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (self *ping) Close() error {
	return self.Conn.Close()
}

func (self *ping) Ping(count int)(string) {
	if err := self.Dail(); err != nil {
		//返回状态0 缺少参数,没有主机IP
		return "is unreachable"
	}
	//fmt.Println("Start ping from ", self.Conn.LocalAddr())
	//ping的回复超时时间,默认1秒
	self.SetDeadline(TimeOut)
	//多次ping
	//for i := 0; i < count; i++ {
	//	r := sendPingMsg(self.Conn, self.Data)
	//	if r.Error != nil {
	//		//失败
	//		return "is unreachable"
	//	} else {
	//		//成功
	//		//fmt.Printf("From %s reply: time=%d ttl=%d\n", self.Addr, r.Time, r.TTL)
	//		return "is alive"
	//	}
	//}
	//只ping一次
	r := sendPingMsg(self.Conn, self.Data)
	if r.Error != nil {
		//失败
		return "is unreachable"
	} else {
		//成功
		//fmt.Printf("From %s reply: time=%d ttl=%d\n", self.Addr, r.Time, r.TTL)
		return "is alive"
	}
	
}

func (self *ping) PingCount(count int) (reply []Reply) {
	if err := self.Dail(); err != nil {
		//fmt.Println("Not found remote host")
		return
	}
	self.SetDeadline(10)
	for i := 0; i < count; i++ {
		r := sendPingMsg(self.Conn, self.Data)
		reply = append(reply, r)
		time.Sleep(1e9)
	}
	return
}

func Run(addr string, req int, data []byte) (*ping, error) {
	wb, err := MarshalMsg(req, data)
	if err != nil {
		return nil, err
	}
	addr, err = Lookup(addr)
	if err != nil {
		return nil, err
	}
	return &ping{Data: wb, Addr: addr}, nil
}

func sendPingMsg(c net.Conn, wb []byte) (reply Reply) {
	start := time.Now()

	if _, reply.Error = c.Write(wb); reply.Error != nil {
		return
	}

	rb := make([]byte, 1500)
	var n int
	n, reply.Error = c.Read(rb)
	if reply.Error != nil {
		return
	}

	duration := time.Now().Sub(start)
	ttl := uint8(rb[8])
	rb = func(b []byte) []byte {
		if len(b) < 20 {
			return b
		}
		hdrlen := int(b[0]&0x0f) << 2
		return b[hdrlen:]
	}(rb)
	var rm *icmp.Message
	rm, reply.Error = icmp.ParseMessage(1, rb[:n])
	if reply.Error != nil {
		return
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		t := int64(duration / time.Millisecond)
		reply = Reply{t, ttl, nil}
	case ipv4.ICMPTypeDestinationUnreachable:
		reply.Error = errors.New("Destination Unreachable")
	default:
		reply.Error = fmt.Errorf("Not ICMPTypeEchoReply %v", rm)
	}
	return
}