#!/bin/bash

# 监控的文件夹路径
WATCH_DIR="./htdocs"

# 读取目标进程的 PID
#TARGET_PID="124"

# 使用 fswatch 监控文件夹
fswatch -0 "$WATCH_DIR" |
while read -d "" event
do
    TARGET_PID=$(cat /tmp/master_pid)
    # 向目标进程发送信号
    kill -SIGUSR1 "$TARGET_PID"
    echo "文件 $event 被修改，发送信号给进程 $TARGET_PID"
done
