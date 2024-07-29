<?php

define('CORE_PATH', ROOT_PATH . "../core/");
define('APP_PATH', ROOT_PATH . "../app/");

require_once CORE_PATH . "Http.php";
require_once CORE_PATH . "Log.php";
require_once CORE_PATH . "YaNanDream.php";
require_once CORE_PATH . "Cathedral.php";
require_once CORE_PATH . "Router.php";
require_once CORE_PATH . "Cntl.php";
require_once CORE_PATH . "Model.php";
require_once CORE_PATH . "View.php";
//require_once "./YaNan.php";

class TheOldHunter
{
    public static $OS_OSX = 1;
    public static $OS_LINUX = 2;

    private $os = 0;
    private $pid_file = '/tmp/master_pid';
    private $master_pid = NULL;
    private $cf = [];
    private $childs = [];

    private $key = NULL;
    private $share_fd = NULL;
    private $listen_socket;

    public static $SERVER_PREPARE = 0;
    public static $SERVER_RUNNING = 1;
    public static $SERVER_RELOAD = 2;
    public static $SERVER_SUSPEND = 3;

    public static $TYPE_SELECT = 0;
    public static $TYPE_EPOLL = 1;
    public static $TYPE_UNKNOW = 9;

    private static $CUR_SERVER_STATUS = 0;
    public static $CUR_SERVER_TYPE = 0;
    public static $CUR_YN_HUNT_CNT_MAX = 0;
    public static $CUR_CF = [];

    public $reload_idx = 1;
    public $reload_max = 9;
    public $share_init_data = "0000000000";
    public $sem_id;

    public $o_event_base;

    public function __construct($cf = [])
    {
        $this->cf['host']                 = $cf['host'] ?? "0.0.0.0";
        $this->cf['port']                 = $cf['port'] ?? 9000;
        $this->cf['master_process_title'] = $cf['master_process_title'] ?? "TheOldHunter Process";
        $this->cf['worker_process_title'] = $cf['worker_process_title'] ?? "TheOldHunter YaNan Process";
        $this->cf['yanan_cnt']            = $cf['worker_process_cnt'] ?? 2;
        $this->cf['yanan_hunt_cnt_max']   = $cf['worker_process_hunt_max'] ?? 10;

        self::$CUR_SERVER_TYPE            = $cf['server_type'] ?? self::$TYPE_EPOLL;
        self::$CUR_CF                     = $this->cf;
    }

    public function hunter()
    {
        Log::info(1);
        $this->parseCmd();
        Log::info(2);
        $this->daemonize();
        Log::info(3);

        $this->key = ftok(__FILE__, 'h');
        $this->sem_id = sem_get($this->key);
        $this->share_fd = shmop_open($this->key, 'c', 0664, 10);
        shmop_write($this->share_fd, $this->share_init_data, 0);
        //shmop_read($this->share_fd, 0, 10);
        //Log::info("share data {$c}");

        cli_set_process_title($this->cf['master_process_title']);
        $this->master_pid = posix_getpid();
        file_put_contents($this->pid_file, $this->master_pid);

        Cathedral::socket($this->listen_socket);
        //$this->serverSocket();

        $this->installSignal();
        $this->forkYnPool();
        $this->loop();
    }

    private function moniterYnPool()
    {
        $info = $this->getYaNan();
        if (count($info[1])) {
            foreach ($this->childs as $pid) {
                if (!in_array($pid, $info[1])) {
                    unset($this->childs[$pid]);
                }
            }
            $this->forkYnPool();
            Log::info("S moniterYnPool");
        }
    }

    private function forkYnPool()
    {
        while (count($this->childs) < $this->cf['yanan_cnt']) {
            $this->forkYaNan();
        }
    }

    private function forkYaNan()
    {
        $pid = pcntl_fork();

        switch ($pid) {
            case -1:
                break;
            case 0:
                $this->master_pid = 0;
                cli_set_process_title($this->cf['worker_process_title']);
                srand();
                mt_srand();
                YaNanDream::factory(
                    $this->share_fd,
                    $this->listen_socket,
                    $this->cf,
                    $this
                )->run();
                //$obj->run();
                /*
                if (self::$CUR_SERVER_TYPE == self::$TYPE_EPOLL) {
                    include ROOT_PATH . "YaNanEpoll.php";
                    $yn = new YaNanEpoll($this->share_fd, $this->key, $this->listen_socket, $this->cf, $this);
                    $yn->run();
                } else {
                    include ROOT_PATH . "YaNanSelect.php";
                    $yn = new YaNanSelect($this->share_fd, $this->key, $this->listen_socket, $this->cf, $this);
                    $yn->run();
                }
                */
                //Log::info("YN KILL SELF");
                //posix_kill(posix_getpid(), SIGKILL);
                exit(0);
            default:
                $this->childs[$pid] = $pid;
                break;
        }
    }

    private function serverSocket()
    {
        $host = '0.0.0.0';
        $port = $this->cf['port'];
        $this->listen_socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
        socket_set_option($this->listen_socket, SOL_SOCKET, SO_REUSEADDR, 1);
        socket_set_option($this->listen_socket, SOL_SOCKET, SO_REUSEPORT, 1);

        socket_bind($this->listen_socket, $host, $port);
        socket_listen($this->listen_socket);
        socket_set_nonblock($this->listen_socket);

        //socket_getsockname($this->listen_socket, $addr, $port);
        //$this->clients = array($this->listen_socket);
    }

    private function assert($exp, $msg)
    {
        if ($exp) {
            echo $msg . PHP_EOL;
            exit;
        }
    }

    private function isAlreadyRunning()
    {
        return file_exists($this->pid_file);
    }

    private function parseCmd()
    {
        switch (TRUE) {
            case stristr(PHP_OS, 'DAR'):
                $this->os = self::$OS_OSX;
                break;
            case stristr(PHP_OS, 'LINUX'):
                $this->os = self::$OS_LINUX;
                break;
            default:
                die("OS ERROR" . PHP_OS);
        }

        $cmd = $_SERVER['argv'][1] ?? '';
        $is_soft = ("--force" == strtolower($_SERVER['argv'][2] ?? '')) ? FALSE : TRUE;
        if (in_array(self::$CUR_SERVER_TYPE, [self::$TYPE_EPOLL, self::$TYPE_UNKNOW])) {
            $is_soft = FALSE;
        }

        switch ($cmd) {
            case 'start':
                $this->assert($this->isAlreadyRunning(), "SERVER ALREADY RUNNING");
                break;
            case 'stop':
                Log::info("S RECV STOP CMD");
                $this->assert(!$this->isAlreadyRunning(), "SERVER NOT RUNNING");
                $info = $this->getYaNan();
                if ($is_soft) {
                    posix_kill($info[0], SIGUSR2);
                    while (TRUE) {
                        $info = $this->getYaNan();
                        if (count($info[1])) {
                            usleep(10000);
                        } else {
                            break;
                        }
                    }
                    posix_kill($info[0], SIGKILL);
                } else {
                    posix_kill($info[0], SIGKILL);
                    $this->killYaNan($info, SIGKILL);
                }
                @unlink($this->pid_file);
                Log::info("S STOP SUC", ["pid" => $info[0]]);
                exit;
            case 'reload':
                Log::info("S RECV RELOAD CMD");
                $this->assert(!$this->isAlreadyRunning(), "SERVER NOT RUNNING");
                $info = $this->getYaNan();
                if ($is_soft) {
                    posix_kill($info[0], SIGUSR1);
                } else {
                    $this->killYaNan($info, SIGKILL);
                }
                Log::info("S RELOAD SUC");
                exit;
            default:
                echo "Enter the command" . PHP_EOL;
                echo "usage start|stop|reload" . PHP_EOL;
                exit;
        }
    }

    private function daemonize()
    {
        $pid = pcntl_fork();
        switch ($pid) {
            case -1:
                break;
            case 0:
                if (($sid = posix_setsid()) < 0) {
                    die("Err: posix_setsid");
                }

                if (chdir('/') === false) {
                    die("Err: chdir");
                }

                umask(0);

                fclose(STDIN);
                fclose(STDOUT);
                fclose(STDERR);

                break;
            default:
                exit;
        }
    }

    private function installSignal()
    {
        pcntl_async_signals(TRUE);
        pcntl_signal(SIGCHLD, [$this, "signalHandel"], false);
        pcntl_signal(SIGUSR1, [$this, "signalHandel"], false);
        pcntl_signal(SIGUSR2, [$this, "signalHandel"], false);
        pcntl_signal(SIGTERM, [$this, "signalHandel"], false);
        //pcntl_sigprocmask(SIG_UNBLOCK, array(SIGTERM,SIGINT));
    }

    private function setShareToReload()
    {
        sem_acquire($this->sem_id); // 获取信号量锁
        if ($this->reload_idx == $this->reload_max) {
            $this->reload_idx = 1;
            $d = str_repeat("0", 7) . self::$SERVER_RELOAD . self::$SERVER_RELOAD;
            $cnt = shmop_write($this->share_fd, $d, 1);
        } else {
            $cnt = shmop_write($this->share_fd, self::$SERVER_RELOAD . "0", $this->reload_idx);
            $this->reload_idx += 1;
        }
        sem_release($this->sem_id); // 释放信号量锁
        Log::info("S SHARE W CNT {$cnt}");
    }

    private function setMasterToSuspend()
    {
        self::$CUR_SERVER_STATUS = self::$SERVER_SUSPEND;
        shmop_write($this->share_fd, self::$CUR_SERVER_STATUS, 0);
    }

    private function setMasterToRunning()
    {
        self::$CUR_SERVER_STATUS = self::$SERVER_RUNNING;
        shmop_write($this->share_fd, self::$CUR_SERVER_STATUS, 0);
    }

    private function signalHandel($signo)
    {
        Log::info("SIG HANDEL");
        if ($this->master_pid == posix_getpid()) {
            switch ($signo) {
                case SIGCHLD:
                    $exit_pid = pcntl_wait($status);
                    unset($this->childs[$exit_pid]);
                    $_["exit_pid"] = $exit_pid;
                    $_["sno"] = SIGCHLD;
                    if ($exit_pid == -1) {
                        $_['msg'] = "pcntl_wait() WARN";
                    }
                    Log::info("S RECV SIG SIGCHID", $_);
                    break;
                case SIGINT:
                    Log::info("SIGINT WILL IGN");
                    break;
                case SIGTERM:
                    Log::info("S RECV TERM");
                    break;
                case SIGUSR1:
                    Log::info("S RECV SIG USR1", ['sno' => SIGUSR1]);
                    $this->setShareToReload();
                    $data = shmop_read($this->share_fd, 0, 10);
                    Log::info("S SHARE D", ['data' => $data]);
                    break;
                case SIGUSR2:
                    Log::info("S RECV SIG USR2", ['sno' => SIGUSR2]);
                    $this->setShareToReload();
                    $this->setMasterToSuspend();

                    $data = shmop_read($this->share_fd, 0, 10);
                    Log::info("S SHARE D", ['data' => $data]);
                    break;
            }
            if (self::isRunning()) {
                $this->forkYnPool();
            }
        } else {
            switch ($signo) {
                case SIGINT:
                    Log::info("SIGINT WILL IGN");
                    break;
                case SIGTERM:
                    Log::info("YN RECV TERM");
                    break;
            }
        }
    }

    public static function isRunning()
    {
        return in_array(self::$CUR_SERVER_STATUS, [self::$SERVER_RUNNING]);
    }

    private function getYaNan()
    {
        $master_pid = file_get_contents($this->pid_file);
        $_exec = "ps --ppid {$master_pid} | awk '/[0-9]/{print $1}' | xargs";
        if (self::$OS_OSX == $this->os) {
            $_exec = "pgrep -P {$master_pid} | xargs";
        }
        Log::info("S GET YN", ['cmd' => $_exec]);
        exec($_exec, $output, $status);
        $childs = [];
        $str = current($output);
        if ($status == 0 && strlen($str)) {
            $childs = explode(" ", $str);
        }
        return [$master_pid, $childs];
    }

    private function killYaNan($info, $signo)
    {
        foreach ($info[1] as $pid) {
            $msg = posix_kill($pid, $signo) ? "OK" : "ERR";
            Log::info("S KILL YN", ['sno' => $signo, 'pid' => $pid, 'msg' => $msg]);
            //usleep(10000);
        }
    }

    private function loop()
    {
        socket_getsockname($this->listen_socket, $addr, $this->cf["port"]);
        Log::info('HTTP SERVER LISTEN', ['addr' => $addr, 'port' => $this->cf["port"], 'pid' => $this->master_pid]);

        self::setMasterToRunning();

        Log::info("S YN", $this->childs);

        $loop_cnt = 1;
        while (TRUE) {
            if ($loop_cnt % 10 == 0) {
                $loop_cnt = 1;
                //$this->moniterYnPool();
            }
            $loop_cnt += 1;
            sleep(1);
        }
    }
}
