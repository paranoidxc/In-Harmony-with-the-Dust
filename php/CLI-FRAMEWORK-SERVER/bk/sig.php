<?php
function signalHandler($signal)
{
    echo "Received signal: $signal\n";
}

// 安装信号处理函数
pcntl_signal(SIGUSR1, "signalHandler");

$pid = pcntl_fork();
$childPid = 0;

if ($pid === -1) {
    die("Fork failed");
} elseif ($pid === 0) {
    // 子进程
    $childPid = pcntl_fork();

    if ($childPid === -1) {
        die("Fork failed");
    } elseif ($childPid === 0) {
        // 孙子进程
        sleep(10); // 假装在做一些工作
    } else {
        // 子进程
        // 这里可以对孙子进程的信号处理进行设置
        while (true) {
            // 子进程在这里等待信号
            pcntl_signal_dispatch();
            usleep(100000); // 等待 0.1 秒
        }
    }
} else {
    // 父进程
    sleep(2); // 等待一段时间，确保子进程和孙子进程都创建完毕
    posix_kill($pid, SIGUSR1); // 向子进程发送信号
    sleep(5); // 等待一段时间，让孙子进程有机会响应信号
    posix_kill($childPid, SIGUSR1); // 向子进程发送终止信号
    pcntl_waitpid($pid, $status); // 等待子进程结束
}
