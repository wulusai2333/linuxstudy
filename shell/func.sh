#!/bin/bash

is_directory(){
  name=$1
  if [ -d $name ];then
    echo "$name 是一个目录"
   else
     mkdir $name
     if [ 0 -ne $? ] ;then
       echo "目录创建失败"
        exit
     fi
     echo "目录创建成功"
  fi
}
#调用函数
is_directory $1