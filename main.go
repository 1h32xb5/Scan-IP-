package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	ip string
)

func init() {
	flag.StringVar(&ip, "ip", "", "IP/Host/Domain")
}

func CmdPing(b string, cha chan int, wg *sync.WaitGroup) {

	sysType := runtime.GOOS

	if sysType == "windows" {
		for p := range cha { 
			//p 1-255 数字
			address := fmt.Sprintf("%s%d", b, p)

			//fmt.Sprintf("%s:%d",b, p)
			cmd := exec.Command("cmd", "/c", "ping -n 1 "+address)
			//						  "bin/bash"
			var out bytes.Buffer 
			cmd.Stdout = &out
			cmd.Run() //运行命令！

			if strings.Contains(out.String(), "TTL=") {
				//fmt.Println("ISOK")
				fmt.Printf("%s主机存活\n", address)
			}
			wg.Done()
		}
	} else if sysType == "linux" {
		for p := range cha {
			//p 1-255 数字
			address := fmt.Sprintf("%s%d", b, p)

			//fmt.Sprintf("%s:%d",b, p)
			cmd := exec.Command("/bin/bash", "/c", "ping -n 1 "+address)
			//						  "bin/bash"
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Run() //运行命令！

			if strings.Contains(out.String(), "TTL=") {
				//fmt.Println("ISOK")
				fmt.Printf("%s主机存活\n", address)
			}
			wg.Done()
		}
	}

}

func main() {
	flag.Parse()
	defer fmt.Println("c段扫描完成！")
	cha := make(chan int, 100000) 
	var wg sync.WaitGroup

	if ip == "" {
		fmt.Println("Please " + os.Args[2] + " -h")
		os.Exit(0)
	}

	a := os.Args[2]
	s := strings.Split(a, ".")
	b := s[0] + "." + s[1] + "." + s[2] + "."
	for i := 0; i < cap(cha); i++ {
		go CmdPing(b, cha, &wg) 
	}
	for i := 1; i < 255; i++ { 
		wg.Add(1)
		cha <- i 
	}
	wg.Wait()
	close(cha)
}
