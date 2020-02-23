package main

import (
	"flag"
	"fmt"
	"github.com/wulusai2333/gostudy/reflect/ini"
	"io"
	"os"
	"sync"
	"time"
)

/*
	go build -o 遍历文件名.exe //生成可执行文件
*/
type FilePath struct {
	Path string `ini:"path"`
}

//创建channel
var (
	fileChan     = make(chan *string, 10000)
	filePathChan = make(chan *string, 10000)
	wg           sync.WaitGroup
)

func main() {
	//从外部加载路径
	var path string
	flag.StringVar(&path, "path", "", "查找文件的路径")
	flag.Parse()
	//没有参数就打开配置文件
	if path == "" {
		file := FilePath{}
		ini.LoadIni(&file, "ini")
		path = file.Path
	}
	//打开程序时没指定就用标准输入指定
	if path == "" {
		fmt.Println("请输入要查找的文件夹路径:")
		_, err := fmt.Fscan(os.Stdin, &path)
		if err != nil {
			fmt.Println("read stdin failed,err:", err)
			return
		}
	}
	//path = filepath.Dir(path)
	//查看路径状态
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("can not find path: %s err:%v\n", path, err)
		return
	}
	//判断路径是不是文件夹

	if !fileInfo.IsDir() {
		fmt.Printf(" '%s' is not dir", path)
		return
	}
	//打开写入文件路径和文件名的文件
	filePathTxt, err := os.OpenFile("./filepath.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print("create/open filepath.txt failed,err:", err)
		return
	}
	fileNameTxt, err := os.OpenFile("./filename.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print("create/open filename.txt failed,err:", err)
		return
	}
	wg.Add(2)
	go writeFileNameTxt(fileNameTxt)
	go writeFileAbsNameTxt(filePathTxt)
	//遍历路径函数
	walk(path)

	close(fileChan)
	close(filePathChan)
	wg.Wait()
}

func walk(absPath string) {
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Println("get path stat failed,err:", err)
		return
	}
	//是文件夹吗
	rangeDir(info, absPath)
}

func rangeDir(info os.FileInfo, absPath string) {

	if info.IsDir() {
		filedir, err := os.Open(absPath)
		if err != nil {
			fmt.Println("open filedir failed,err:", err)
			return
		}
		for {
			infos, err := filedir.Readdir(10)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("read dir failed,err:", err)
				return
			}
			for _, v := range infos {
				path := absPath + "\\" + v.Name()
				rangeDir(v, path)
			}
		}
	}

	if !info.IsDir() {
		//不是文件夹
		//fmt.Println(info.Name())
		name := info.Name()
		fileChan <- &name
		realPath := absPath
		filePathChan <- &realPath
		return
	}
}

func writeFileNameTxt(fileNameTxt *os.File) {
	defer wg.Done()
	for {
		select {
		case fi, ok := <-fileChan:
			if !ok {
				fileNameTxt.Close()
				return
			}
			var _, err = fileNameTxt.Write([]byte(*fi + "\n"))
			if err != nil {
				fmt.Println("write filename failed,err:", err)
				return
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}

}

func writeFileAbsNameTxt(filePathTxt *os.File) {
	defer wg.Done()
	for {
		select {
		case fp, ok := <-filePathChan:
			if !ok {
				filePathTxt.Close()
				return
			}
			_, err := filePathTxt.Write([]byte(*fp + "\n"))
			if err != nil {
				fmt.Println("write filename failed,err:", err)
				return
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
