package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)
var (
	host string
)

func init() {
	flag.StringVar(&host, "host", "", "IP/Host/Domain")
}
func CmdPing(b string,ip chan int, wg *sync.WaitGroup)  {

	for p := range ip{ 				//channel ports 里面有缓冲 就遍历 出来！ 然后 执行完 之后-1
	//p 1-255 数字
		address :=  fmt.Sprintf("%s%d",b,p)

			//fmt.Sprintf("%s:%d",b, p)
		cmd := exec.Command("cmd","/c","ping -a -n 1 "+address)
		//						  "bin/bash"
		var out bytes.Buffer   //bytes.buffer是一个缓冲byte类型的缓冲器存放着都是byte
		cmd.Stdout = &out
		cmd.Run()          //运行命令！

		if strings.Contains(out.String(), "TTL=") {
			//fmt.Println("ISOK")
			fmt.Printf("%s主机存活\n", address)
		}
		wg.Done()
	}
}

func main() {
	flag.Parse()
	defer fmt.Println("c段扫描完成！")
	ip := make(chan int , 100000)      //开启1000个缓冲通道   /也就是下面以为着开几个go程
	var wg sync.WaitGroup

	if host  == "" {
		fmt.Println("Please " + os.Args[2] + " -h")
		os.Exit(0)
	}
	a := os.Args[2]
	s := strings.Split(a, ".")
	b := s[0]+"."+s[1]+"."+s[2]+"."
	for i := 0;i<cap(ip);i++{
		go CmdPing(b,ip,&wg)         //传参数 为host ip(管道channel)  和wg计数器
	}
	for i := 1;i<255;i++{					//扫描一个段里的
		wg.Add(1)
		ip <- i      		//传入channel 缓冲管道 里面数据
	}
	wg.Wait()
	close(ip)
}
