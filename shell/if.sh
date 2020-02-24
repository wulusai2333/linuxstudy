#!/bin/bash

if [ -d $1 ];then
  echo "$1 是一个目录"
  elif [ -s $1 ]; then
  echo "$1 存在且大小为0"
  else
    echo "$1 可能不存在,也可能大小为0"
fi