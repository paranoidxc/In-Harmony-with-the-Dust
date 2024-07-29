<?php
$cnt = 1;
$max = 5;
function fork()
{
    $pid = pcntl_fork();
    if ($pid == 0) {
        echo posix_getpid() . PHP_EOL;

        pcntl_sigprocmask(SIG_UNBLOCK, array(SIGTERM));


        while (true) {
            pcntl_signal_dispatch();
            usleep(100000); // 等待 0.1 秒
        }
    }
}

echo posix_getpid() . PHP_EOL;
//pcntl_async_signals(true);
// 给进程安装信号...

pcntl_signal(SIGINT, function () {
    echo "INT" . PHP_EOL;
});

pcntl_signal(SIGCHLD, function () {
    echo "SIGCHLD" . PHP_EOL;
    pcntl_wait($status);
    fork();
});

for ($i = 0; $i < $cnt; $i++) {
    fork();
}
// while保持进程不要退出..
while (true) {
    pcntl_signal_dispatch();
    usleep(100000); // 等待 0.1 秒
}
