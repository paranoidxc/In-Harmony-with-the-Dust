<?php
class Http
{
    // 定义下目前支持的http方法们，目前只支持get和post
    private static $a_method = array('get', 'post');
    public static function decode($s_raw_http_content)
    {
        $s_http_method       = '';
        $s_http_version      = '';
        $s_http_pathinfo     = '';
        $s_http_querystring  = '';
        $s_http_body_boundry = '';  // 当post方法且为form-data的时候.
        $a_http_post         = [];
        $a_http_get          = [];
        $a_http_header       = [];
        $a_http_file         = [];
        $json                = [];

        // 先通过两个 \r\n\r\n 把 请求行+请求头 与 请求体 分割开来.
        list($s_http_line_and_header, $s_http_body) = explode("\r\n\r\n", $s_raw_http_content, 2);
        // 再分解$s_http_line_and_header数组
        // 数组的第一个元素一定是 请求行
        // 数组剩余所有元素就是 请求头
        $a_http_line_header = explode("\r\n", $s_http_line_and_header);
        $s_http_line = $a_http_line_header[0];
        unset($a_http_line_header[0]);
        $a_http_raw_header = $a_http_line_header;
        // 好了，请求行 + 请求头数组 + 请求体 都有了
        // 先从请求行分解 method + pathinfo + querystring + http版本
        list($s_http_method, $s_http_pathinfo_querystring, $s_http_version) = explode(' ', $s_http_line);
        if (false === strpos($s_http_pathinfo_querystring, "?")) {
            $s_http_pathinfo = $s_http_pathinfo_querystring;
        } else {
            list($s_http_pathinfo, $s_http_querystring) = explode('?', $s_http_pathinfo_querystring);
        }
        // 处理querystring为数组
        if ('' != $s_http_querystring) {
            $a_raw_http_get = explode('&', $s_http_querystring);
            foreach ($a_raw_http_get as $s_http_get_item) {
                if ('' != trim($s_http_get_item)) {
                    list($s_get_key, $s_get_value) = explode('=', $s_http_get_item);
                    $a_http_get[$s_get_key] = $s_get_value;
                }
            }
        }
        // 处理$s_http_header
        foreach ($a_http_raw_header as $trash => $a_raw_http_header_item) {
            if ('' != trim($a_raw_http_header_item)) {
                list($s_http_header_key, $s_http_header_value) = explode(":", $a_raw_http_header_item);
                $a_http_header[strtoupper($s_http_header_key)] = $s_http_header_value;
            }
        }

        ////mylog("HEADER");
        foreach ($a_http_header as $k => $v) {
            //mylog("{$k}={$v}");
        }
        ////mylog("METHOD");
        ////mylog($s_http_method);
        ////mylog($a_http_header['CONTENT-TYPE']);
        //application/x-www-form-urlencoded
        // 如果是post方法，处理post body
        if ('post' === strtolower($s_http_method)) {
            ////mylog("POST");
            ////mylog($a_http_header['CONTENT-TYPE']);
            // post 方法里要关注几种不同的content-type
            if ('application/x-www-form-urlencoded' == trim($a_http_header['CONTENT-TYPE'])) {
                ////mylog($s_http_body);
                $a_http_raw_post = explode("&", $s_http_body);
                // 解析http body
                foreach ($a_http_raw_post as $s_http_raw_body_item) {
                    if ('' != $s_http_raw_body_item) {
                        ////mylog($s_http_raw_body_item);
                        list($s_http_raw_body_key, $s_http_raw_body_value) = explode("=", $s_http_raw_body_item);
                        $a_http_post[$s_http_raw_body_key] = $s_http_raw_body_value;
                    }
                }
            }

            // json
            if ('application/json' == trim($a_http_header['CONTENT-TYPE'])) {
                if (strlen($s_http_body)) {
                    $json = json_decode($s_http_body, TRUE);
                    if ($json == null) {
                        $json = array();
                    }
                }
            }

            // form-data
            if (false !== strpos($a_http_header['CONTENT-TYPE'], 'multipart/form-data')) {
                ////mylog("FORM-DATA");
                list($s_http_header_content_type, $s_http_body_raw_boundry) = explode(';', $a_http_header['CONTENT-TYPE']);
                ////mylog("BOUNDARY");
                ////mylog($s_http_body_raw_boundry);
                $a_http_header['CONTENT-TYPE'] = trim($s_http_header_content_type);
                list($nil, $s_http_body_boundry) = explode('=', $s_http_body_raw_boundry);

                ////mylog("RREAL BOUNDARY");
                ////mylog($s_http_body_boundry);

                $s_http_body_boundry_end = '--' . $s_http_body_boundry . '--';
                ////mylog("BBBBODY");
                ////mylog($s_http_body);
                $a_http_raw_post     = explode('--' . $s_http_body_boundry . "\r\n", $s_http_body);
                foreach ($a_http_raw_post as $s_http_raw_body_item) {
                    ////mylog("====START");
                    ////mylog($s_http_raw_body_item);
                    //mylog("====END");
                    //$s_http_raw_body_item = trim($s_http_raw_body_item);
                    if ('' != $s_http_raw_body_item) {
                        $pos = strpos($s_http_raw_body_item, "\r\n");
                        $first_line = substr($s_http_raw_body_item, 0, $pos);
                        //mylog("FIRST_LINE" . $first_line);
                        $first_line_info = explode(";", $first_line);
                        if (count($first_line_info) == 2) {
                            $other_line = substr($s_http_raw_body_item, $pos + 4);
                            $other_line = substr($other_line, 0, strlen($other_line) - 2);
                            //mylog("OTHER_LINE" . $other_line);
                            $field_items = explode('"', $first_line);
                            $field_name = $field_items[1];
                            $field_val = str_replace("\r\n" . $s_http_body_boundry_end, '', $other_line);
                            $a_http_post[$field_name] = $field_val;
                        } else if (count($first_line_info) == 3) {
                            //mylog("UPLOAD=========");
                            $field_items = explode('"', $first_line);
                            $field_name = $field_items[1];
                            $field_fname = $field_items[3];
                            $other_line = substr($s_http_raw_body_item, $pos + 2);

                            //mylog($other_line);
                            $pos = strpos($other_line, "\r\n");
                            $first_line = substr($other_line, 0, $pos);

                            //mylog("FIRSTLINE");
                            //mylog($first_line);
                            $other_line = substr($other_line, $pos + 4);
                            $other_line = substr($other_line, 0, strlen($other_line) - 2);

                            //mylog("FFFFFFFFFFFFFFFF");
                            //mylog($other_line);

                            file_put_contents("/tmp/{$field_fname}", $other_line);

                            $a_http_post[$field_name] = $field_fname;
                        }
                    }
                    //mylog("");
                    //mylog("");
                }
            }
        }

        $a_ret = array(
            'method'   => $s_http_method,
            'version'  => $s_http_version,
            'pathinfo' => $s_http_pathinfo,
            'post'     => $a_http_post,
            'get'      => $a_http_get,
            'header'   => $a_http_header,
            'file'     => [],
            'json'     => $json,
        );

        return $a_ret;
    }

    public static function encode($a_data)
    {
        //echo "RESPONSE " . PHP_EOL;
        //$s_data        = json_encode( $a_data );
        $s_data = "<html><body>" . $a_data . "</body></html>";
        $s_http_line   = "HTTP/1.1 200 OK";
        $a_http_header = array(
            "Content-Type"   => "text/html;charset=UTF-8",
            "Connection" => "keep-alive",
            //"Date"           => gmdate("M d Y H:i:s", time()),
            //"Content-Type"   => "application/html",
            "Server"   => "ftw",
            "Content-Length" => strlen($s_data),
        );
        $s_http_header = '';
        foreach ($a_http_header as $s_http_header_key => $s_http_header_item) {
            $_s_header_line = $s_http_header_key . ': ' . $s_http_header_item;
            $s_http_header  = $s_http_header . $_s_header_line . "\r\n";
        }
        $s_ret = $s_http_line . "\r\n" . $s_http_header . "\r\n" . $s_data;
        //echo $s_ret;
        return $s_ret;
    }
}
