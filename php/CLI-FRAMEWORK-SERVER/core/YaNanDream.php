<?php
class YaNanDream
{
    public static function factory($share_fd, $listen_socket, $cf, $that)
    {
        //Log::info(TheOldHunter::$CUR_SERVER_TYPE);
        $obj = NULL;
        switch (TheOldHunter::$CUR_SERVER_TYPE) {
            case TheOldHunter::$TYPE_EPOLL:
                include CORE_PATH . "YaNanEpoll.php";
                $obj = new YaNanEpoll($share_fd, $listen_socket, $cf, $that);
                break;
            case TheOldHunter::$TYPE_SELECT:
                include CORE_PATH . "YaNanSelect.php";
                $obj = new YaNanSelect($share_fd, $listen_socket, $cf, $that);
                //$obj = new YaNanSelect($that);
                break;
            default:
                $obj = new YaNanDream();
                break;
        }
        return $obj;
    }

    public function run()
    {
        while (TRUE) {
            Log::err("Factory Faild");
            sleep(1);
        }
    }
}
