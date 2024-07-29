<?php
class YaNanEpoll
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
    public $o_event_base;
    public $a_event_array = [];
    public $a_client_array = [];

    public function __construct($share_fd, $listen_socket, $cf, $that)
    {
        $this->pid           = posix_getpid();
        $this->that          = $that;
        $this->share_fd      = $share_fd;
        $this->listen_socket = $listen_socket;
        $this->cf            = $cf;
    }

    public function run()
    {
        $share_data = shmop_read($this->share_fd, 0, 10);
        Log::info("YN EPOLL RUN POS = {$this->that->reload_idx} status = $share_data", ["pid" => $this->pid]);

        $o_event_config = new EventConfig();
        // 通过requireFeatures方法来配置控制
        $o_event_config->requireFeatures(EventConfig::FEATURE_ET);
        $o_event_config->requireFeatures(EventConfig::FEATURE_O1);
        //$o_event_config->requireFeatures( EventConfig::FEATURE_FDS );
        $o_event_base = new EventBase($o_event_config);
        //$o_event_base  = new EventBase();
        $s_method_name = $o_event_base->getMethod();
        Log::info("METHOD", ['d' => $s_method_name]);
        /*
        if ('epoll' != $s_method_name) {
            exit("not epoll");
        }
        */
        Log::info("YN EPOLL 2");
        $o_event = new Event($o_event_base, $this->listen_socket, Event::READ | Event::PERSIST, function ($r_listen_socket, $i_event_flag, $o_event_base) {
            $r_connection_socket = socket_accept($r_listen_socket);
            if (false === $r_connection_socket) {
                Log::info("RE RUN");
                $this->run();
            }
            $this->a_client_array[] = $r_connection_socket;

            Log::info("Event::accept " . $this->pid);

            // 在这个客户端连接socket上添加 读事件
            $o_read_event = new Event(
                $o_event_base,
                $r_connection_socket,
                Event::READ | Event::PERSIST,
                function ($r_connection_socket, $i_event_flag, $o_event_base) {
                    Cathedral::read($r_connection_socket, $client_content);
                    //$s_content = socket_read($r_connection_socket, 1024);
                    Log::info("接受到：" . strlen($client_content) . " --- " . $client_content);

                    $decode_ret = Http::decode($client_content);
                    //$path = str_replace("/", "", $decode_ret['pathinfo']);
                    //$path = str_replace("/", "", $decode_ret['pathinfo']);
                    Log::info("content", $decode_ret);
                    $out = Router::run($decode_ret);

                    $o_event = $this->a_event_array[intval($r_connection_socket)]['read'];
                    $o_event->del();
                    unset($this->a_event_array[intval($r_connection_socket)]['read']);

                    //Log::info("out".$out);

                    // 当这个客户端连接socket一旦满足可写条件，我们就可以向socket中写数据了
                    $o_write_event = new Event(
                        $o_event_base,
                        $r_connection_socket,
                        Event::WRITE | Event::PERSIST,
                        function ($r_connection_socket, $i_event_flag, $out) {
                            Log::info("Event::write回调 " . $this->pid);
                            //Log::info("OUT =====".$out);
                            //$out = "EPOLL RES {$this->pid} " . date("Y-m-d H:i:s");
                            $s_content = Http::encode($out);

                            //Log::info($out);
                            //Log::info($s_content);
                            socket_write($r_connection_socket, $s_content, strlen($s_content));

                            // 在写回调中逻辑执行完毕后，将该写事件删除掉...
                            $o_event = $this->a_event_array[intval($r_connection_socket)]['write'];
                            $o_event->del();
                            unset($this->a_event_array[intval($r_connection_socket)]['write']);

                            socket_close($r_connection_socket);
                        }, $out);

                    $o_write_event->add();
                    $this->a_event_array[intval($r_connection_socket)]['write'] = $o_write_event;
                },
                $o_event_base
            );
            $o_read_event->add();
            $this->a_event_array[intval($r_connection_socket)]['read'] = $o_read_event;
        }, $o_event_base);

        /*
        $o_timer_event = new Event($o_event_base, -1, Event::TIMEOUT | Event::PERSIST, function () {
            Log::info("BINGO");
        });
        $o_timer_event->add(0.7);
        */

        $o_event->add();
        $o_event_base->loop();
    }

    public function testTimer()
    {
        $r_listen_socket = $this->listen_socket;
        $share_data = shmop_read($this->share_fd, 0, 10);
        Log::info("YN RUN POS = {$this->that->reload_idx} status = $share_data", ["pid" => $this->pid]);
        $o_event_base  = new EventBase();
        $s_method_name = $o_event_base->getMethod();
        if ('epoll' != $s_method_name) {
            exit("not epoll");
        }
        $o_timer_event = new Event($o_event_base, -1, Event::TIMEOUT | Event::PERSIST, function () {
            Log::info("BINGO");
        });
        $o_timer_event->add(0.7);
        $o_event_base->loop();
    }
}
