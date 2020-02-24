#!/bin/bash
list=`ls`
for var in $list ; do
    echo "当前文件: $var"
done