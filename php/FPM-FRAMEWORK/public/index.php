<?php

declare(strict_types=1);

use \Paranoid\Framework\Http\Request;
use \Paranoid\Framework\Http\Kernel;
use \Paranoid\Framework\Routing\Router;

define('BASE_PATH', dirname(__DIR__));

require_once BASE_PATH . '/vendor/autoload.php';

$container = require BASE_PATH . '/config/services.php';

//dd($container);


// request received
$request = Request::createFromGlobals();

//dd($request);

//$router = new Router();
// perform some logic
//$kernel = new Kernel($router);

$kernel = $container->get(Kernel::class);

// send response (string of content)
$response = $kernel->handle($request);

$response->send();

$kernel->terminate($request, $response);