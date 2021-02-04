package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)


var wg sync.WaitGroup


func testtime(start time.Time) {
	fmt.Printf("耗时: %v", time.Since(start))
}


func test01(file *os.File, url string, urls chan string) {
	m := bufio.NewScanner(file)
	for m.Scan() {
		var okurl = fmt.Sprintf("%s%s", url, m.Text())
		urls <- okurl
	}
	err := m.Err()
	if err != nil {
		fmt.Println("错误")
	}
	close(urls)
	fmt.Println("读取完毕")
}

func gourl(urls chan string)  {
	for {
		select {
		case url, ok := <-urls:
			if !ok {
				wg.Done()
				return
			}
			client := &http.Client{}
			request, err := http.NewRequest("HEAD", url, nil)
			if err != nil {
				continue
			}
			request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:84.0) Gecko/20100101 Firefox/84.0")
			resp, err := client.Do(request)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if resp.StatusCode == 200 {
				fmt.Printf("%s 状态码: %v\n", url, resp.StatusCode)
			} else {
				fmt.Printf("%s 状态码:%v\n", url, resp.StatusCode)
			}
			resp.Body.Close()
			case <- time.After(time.Duration(2) * time.Second):
				wg.Done()
				return
		}
	}
}

func main()  {
	var url string
	flag.StringVar(&url, "u", "", "输入目标地址")
	flag.Parse()
	if url == "" {
		fmt.Println(`
 __          __       _                   
 \ \        / /      | |                  
  \ \  /\  / /__  ___| |_ ___  _ __   ___ 
   \ \/  \/ / _ \/ __| __/ _ \| '_ \ / _ \
    \  /\  /  __/\__ \ || (_) | | | |  __/
     \/  \/ \___||___/\__\___/|_| |_|\___|

		磐石目录扫描工具
		-h        查看帮助
		版本       1.0
		`)
		os.Exit(1)
	}
	times := time.Now()
	urls := make(chan string)
	file, err := os.Open("dir.txt")
	if err != nil {
		fmt.Println("字典打开失败: ", err)
		return
	}
	defer testtime(times)
	defer file.Close()
	go test01(file, url, urls)
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go gourl(urls)
	}
	wg.Wait()
}