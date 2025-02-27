#!/bin/bash

# 删除/home目录下所有文件
rm -rf /home/*

# 查找/opt/temp目录下以upgrade开头的.zip文件
upgrade_zip=$(find /opt/temp -name 'upgrade*.zip' -print -quit)

# 检查是否找到了upgrade开头的.zip文件
if [ -z "$upgrade_zip" ]; then
    echo "未找到以upgrade开头的.zip文件"
    exit 1
fi

# 将找到的upgrade开头的.zip文件解压缩到/home目录下
unzip "$upgrade_zip" -d /home

# 删除/opt/temp目录
rm -rf /opt/temp

# 对/home目录下的所有文件和目录赋予执行权限
chmod +x -R /home

# 重启系统
reboot