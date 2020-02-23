#!/bin/bash
#第一行如果忘了写会使用默认的命令解析器 /bin/sh
#指定shell脚本的命令解析器 linux终端打开的时候后台就有一个默认的bash进程
echo $PATH #打印path
which go #看看 go 这个命令在哪
pwd #进程工作目录
chmod u+x test.sh # u是文件所有者 +x 是加可执行权限
sh test.sh #作用跟 ./test.sh 一样
./test.sh