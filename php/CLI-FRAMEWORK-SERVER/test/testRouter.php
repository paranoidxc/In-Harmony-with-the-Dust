<?php

define('ROOT_PATH',  dirname(__FILE__) . '/../public/'); // 站点应用目录
define('CORE_PATH', ROOT_PATH . "../core/");
define('APP_PATH', ROOT_PATH . "../app/");


require_once ROOT_PATH."../core/Log.php";
require_once ROOT_PATH."../core/Router.php";

$pathinfo = "/u/login/name/hxc/";
$pathinfo = "/products";
$pathinfo = "/";

//Router::run($pathinfo);
