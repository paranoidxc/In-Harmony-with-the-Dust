<?php
class YaNanSelect
{
    public $status = 1;
    public $pid = 0;
    public $key;
    public $share_fd;
    public $clients = [];
    public $listen_socket = NULL;
    public $cf = [];
    public $sig_queue = [];
    public $that;

    public function __construct($share_fd, $listen_socket, $cf, $master)
    {
        $this->pid           = posix_getpid();

        $this->that          = $master;
        $this->share_fd      = $share_fd;
        $this->listen_socket = $listen_socket;
        $this->cf            = $cf;
        $this->clients       = array($this->listen_socket);
        //Log::info('CF', ['data'=>$master->cf]);
    }

    public function run()
    {
        set_error_handler(function ($errno, $errstr, $errfile, $errline) {
            $d = ['e' => $errno, 'f' => $errfile, 'l' => $errline, 's' => $errstr];
            Log::info("ERR", ['data' => $d]);
            return false;
        });

        $yn_hunt_cnt = 0;
        while (TRUE) {
            $share_data = shmop_read($this->share_fd, 0, 10);
            $s_status = $share_data[0];
            $is_reload = $share_data[$this->that->reload_idx];

            //Log::info("YN RUN POS = {$this->that->reload_idx} status = $share_data", ["pid" => $this->pid]);
            $read      = $this->clients;
            $write     = [];
            $exception = [];
            $ret       = socket_select($read, $write, $exception, 0);
            if ($ret <= 0) {
                // exit current worker to start new
                if ($yn_hunt_cnt > $this->cf['yanan_hunt_cnt_max']) {
                    Log::info("YN HUNT MAX EXIT", ['pid' => posix_getpid()]);
                    exit(0);
                }
                if (TheOldHunter::$SERVER_RELOAD == $is_reload) {
                    Log::info("YN WILL EXIT");
                    exit(0);
                    break;
                }
                usleep(1000);
                //sleep(3);
                continue;
            }
            // new accept
            if (in_array($this->listen_socket, $read)) {
                $connection_socket = socket_accept($this->listen_socket);
                if (!$connection_socket) {
                    continue;
                }
                $yn_hunt_cnt++;
                socket_getpeername($connection_socket, $client_ip, $client_port);
                Log::info("YN ACC", [
                    "ip"      => $client_ip,
                    "port"    => $client_port,
                    "pid"     => posix_getpid(),
                    "hunt_cnt" => $yn_hunt_cnt,
                ]);
                $this->clients[] = $connection_socket;
                $key      = array_search($this->listen_socket, $read);
                unset($read[$key]);
            }
            // other socket
            foreach ($read as $read_key => $read_fd) {
                //require ROOT_PATH . "Http.php";
                $client_content = '';
                $out = '';
                if (FALSE) {
                    $out = "response " . date("Y-m-d H:i:s");
                } else {
                    //$cnt_ret = socket_recv($read_fd, $client_content, 65535, 0);
                    //$s_content = socket_read($r_connection_socket, 1024);

                    /*
                    Log::info("MSG_OOB ".MSG_OOB);
                    Log::info("MSG_PEEK ".MSG_PEEK);
                    Log::info("MSG_WAIT ".MSG_WAITALL);
                    Log::info("MSG_NOWAIT ".MSG_DONTWAIT);
                    */

                    Cathedral::read($read_fd, $client_content);
                    /*
                    $client_content = Cathedral::read($client_content);
                    while (1) {
                        $cnt_ret = socket_recv($read_fd, $tmp_recv_content, 1024, MSG_DONTWAIT);
                        if ($cnt_ret > 0 && $cnt_ret <= 1024) {
                            $client_content .= $tmp_recv_content;
                        } else {
                            break;
                        }
                    }
                    */

                    $decode_ret = Http::decode($client_content);
                    //$_GET = $_POST = $_REQUEST = $_FILES = [];
                    //$_GET = $decode_ret['get'];
                    //$path = str_replace("/", "", $decode_ret['pathinfo']);
                    $out = Router::run($decode_ret);

                    //TODO:ROUTE
                    /*
                    $out = '';
                    $php_file = 'index.php';
                    if (in_array($path, ["a.php", "b.php", "c.php"])) {
                        $php_file = $path;
                    }
                    if ($php_file) {
                        ob_start();
                        require ROOT_PATH . "htdocs/{$php_file}";
                        $out = ob_get_contents();
                        ob_end_clean();
                    }
                    */
                }

                $encode_ret = Http::encode($out);
                socket_write($read_fd, $encode_ret, strlen($encode_ret));

                socket_close($read_fd);
                unset($read[$read_key]);
                $key = array_search($read_fd, $this->clients);
                unset($this->clients[$read_key]);
            }
        }
    }
}
