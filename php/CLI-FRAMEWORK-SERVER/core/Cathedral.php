<?php
class Cathedral
{
    public static function socket(&$listen_socket)
    {
        $listen_socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
        socket_set_option($listen_socket, SOL_SOCKET, SO_REUSEADDR, 1);
        socket_set_option($listen_socket, SOL_SOCKET, SO_REUSEPORT, 1);

        socket_bind($listen_socket, TheOldHunter::$CUR_CF['host'], TheOldHunter::$CUR_CF['port']);
        socket_listen($listen_socket);
        socket_set_nonblock($listen_socket);
    }

    public static function read(&$fd, &$client_content)
    {
        while (($cnt_ret = socket_recv($fd, $tmp_recv_content, 1024, MSG_DONTWAIT))) {
            $client_content .= $tmp_recv_content;
            if ((int)$cnt_ret < 1024) {
                break;
            }
        }
    }
}
