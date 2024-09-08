<?php

declare(strict_types=1);

use \Paranoid\Framework\Http\Request;
use \Paranoid\Framework\Http\Kernel;
use \Paranoid\Framework\Routing\Router;

define('BASE_PATH', dirname(__DIR__));

require_once dirname(__DIR__) . '/vendor/autoload.php';

// request received
$request = Request::createFromGlobals();

//dd($request);

$router = new Router();
// perform some logic
$kernel = new Kernel($router);

// send response (string of content)
$response = $kernel->handle($request);

$response->send();
