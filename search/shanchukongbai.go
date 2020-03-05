package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./xiuxian.txt")
	if err != nil {
		fmt.Println("read file failed")
		return
	}
	wfile, err := os.OpenFile("./修真聊天群.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("read file failed", err)
		return
	}
	readFile := bufio.NewReader(file)
	for {
		line, isPrefix, err := readFile.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read line failed", err)
			return
		}
		if isPrefix {
			fmt.Println("to long")
			str := string(line)
			if str != "" {
				_, err := wfile.Write([]byte(str))
				if err != nil {
					fmt.Println("write err", err)
				}
			}
		}
		str := string(line)
		str = strings.Trim(str, " ")
		str = strings.Trim(str, "\n")
		str = strings.Trim(str, "\r\n")
		//strings.Split()
		if str != "" {
			if strings.HasPrefix(str, "第") {
				if strings.HasSuffix(strings.Split(str, " ")[0], "章") {
					str = "\n" + str + "\n"
				}
			}

			_, err := wfile.Write([]byte(str))
			if err != nil {
				fmt.Println("write err", err)
			}
		}

	}

}
