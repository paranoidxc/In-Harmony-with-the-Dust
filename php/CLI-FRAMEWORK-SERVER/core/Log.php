<?php
class Log
{
    public static $LEVEL_DEF = 0;
    public static $LEVEL_INFO = 0;
    public static $LEVEL_DEBUG = 1;
    public static $LEVEL_WARN = 2;
    public static $LEVEL_ERR = 3;

    public static $LEVEL_OPTS = [
        0 => 'INFO',
        1 => 'DEBUG',
        2 => 'WARN',
        3 => 'ERR',
    ];

    private static function _log($level = 0, $msg, $data)
    {
        $level_str = self::$LEVEL_OPTS[$level] ?? 'UNT';
        $tmp = [
            date("Y-m-d H:i:s"),
            $level_str,
        ];

        if (is_array($msg) && count($msg)) {
            $tmp[] =  json_encode($msg, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
        } else if (is_string($msg) && strlen($msg)) {
            $tmp[] = $msg;
        }

        if (is_array($data) && count($data)) {
            $tmp[] =  json_encode($data, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
        } else if (is_string($data) && strlen($data)) {
            $tmp[] = $data;
        }
        file_put_contents(ROOT_PATH . '../logs/log', join(" | ", $tmp) . PHP_EOL, FILE_APPEND);
        //file_put_contents("/Users/xc/log", join(" | ", $tmp) . PHP_EOL, FILE_APPEND);
    }

    public static function info($msg, $data = [])
    {
        self::_log(self::$LEVEL_INFO, $msg, $data);
    }

    public static function debug($msg, $data = [])
    {
        self::_log(self::$LEVEL_DEBUG, $msg, $data);
    }

    public static function warn($msg, $data = [])
    {
        self::_log(self::$LEVEL_WARN, $msg, $data);
    }

    public static function err($msg, $data = [])
    {
        self::_log(self::$LEVEL_ERR, $msg, $data);
    }
}
