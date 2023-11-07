#!/bin/bash
sudo /usr/bin/timedatectl set-time "$1"
echo "設定日期成功"
echo "$1"
#date +'%Y-%m-%d %H:%M:%S'
