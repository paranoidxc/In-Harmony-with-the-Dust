#!/bin/bash

# 监控的文件夹路径
WATCH_DIR="./htdocs"

# 目标进程的进程ID（PID）

# 监控文件修改的事件
EVENTS="modify"

# 使用 inotifywait 监控文件夹
inotifywait -m -e "$EVENTS" "$WATCH_DIR" |
while read -r directory event file
do
    TARGET_PID=$(cat /tmp/master_pid)
    # 只处理文件修改事件
    if [[ "$event" == "MODIFY" ]]; then
        # 向目标进程发送信号
        kill -SIGUSR1 "$TARGET_PID"
        echo "文件 $file 被修改，发送信号给进程 $TARGET_PID"
    fi
done
