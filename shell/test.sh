#!/bin/bash
#第一行如果忘了写会使用默认的命令解析器 /bin/sh
#指定shell脚本的命令解析器 linux终端打开的时候后台就有一个默认的bash进程
# shell脚本的基本命令 set设置系统变量 export设置变量 $0 $1 位置变量 变量取值 if判断 for循环 函数
#写了一个小功能，遍历指定文件夹的所有文件，输出文件名及文件全路径到两个文件中，http网页搜索关键词，程序加载这两个文件从中找到匹配项返回到浏览器
echo $PATH #打印path
which go #看看 go 这个命令在哪
pwd #进程工作目录
chmod u+x test.sh # u是文件所有者 +x 是加可执行权限
sh test.sh #作用跟 ./test.sh 一样
./test.sh
#shell中定义变量 shell的变量等号前后不能空格 " = " 的等号等价于 "=="
temp=123  #变量只能在当前进程使用
#环境变量 全局变量 在当前系统内可以访问
env #列出当前系统所有环境变量
#系统自带
PWD
SHELL
HOME
PATH
#用户自定义的
GOPATH=/user/go/src #目前是普通变量
set GOPATH=/user/go/src #提升为系统变量
export GOPATH=/user/go/src #导入
~/.bashrc #上面的导入到这里

#位置变量
./test aa bb cc
#如何获取参数
$0 #脚本本身的名字
$1 #第一个参数 aa

#特殊变量
$# #脚本传递参数的个数
$@ #脚本传递的所有参数
$? #脚本执行的情况 成功=0 失败>0
$$ #脚本进程执行之后对应的进程ID

#给变量赋值默认以字符串 以下两个等价
value=123
value="123"
#变量取值 可以不加大括号
echo $value
echo ${value}
var=$(pwd) #小括号中是命令 命令的执行结果给var
var=`pwd` #效果同上

#if 语句 注意 命令 [ ] 参数 都有空格
if [ true ] ;then
  echo "逻辑语句"
fi
#---还可以
if [ 1 ]
then
  echo "逻辑语句"
fi
#else if
if [ true ] ;then
  echo "逻辑语句"
elif [ true ]; then
      echo "其他逻辑"
fi
#常用的判断条件
if [ -d filename ] #存在且是目录返回true
if [ -s filename ] #文件或目录存在且大小>0返回true
if [ -w pathname ] #文件或目录存在且可写时返回true
if [ -x pathname ] #文件或目录存在且可执行时返回true
#字符串测试
if [ -z string ] #字符串为空返回true
if [ -n string ] #字符串不为空返回true
if [ str1 > str2 ] #按字典顺序排序 str1在str2之后
if [ str1 < str2 ] #按字典顺序排序 str1在str2之前
#数值测试
if [ int1 -eq int2 ] #等于
if [ int1 -ne int2 ] #不等于
if [ int1 -lt int2 ] #小于
if [ int1 -le int2 ] #小于等于
if [ int1 -gt int2 ] #大于
if [ int1 -ge int2 ] #大于等于
#逻辑操作符
if [ a -a b ] #与
if [ a -o b ] #或
if [ ! b ] #非
touch t.txt #创建目录

#for循环

for i in {1..5} ; do
  echo "循环条件 $i"
done
#引号的使用
#单引号
echo '原样输出 $1'
#双引号
echo "会取值 $1"

#函数的书写规范 有小括号但是不传参 传参靠$ 无返回值
funcName(){
  echo $1
  arg2=$2
}
#函数调用
funcName 1 2 3
#返回参数 0成功 非0失败