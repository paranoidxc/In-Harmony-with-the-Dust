<?php

declare(strict_types=1);

use \Paranoid\Framework\Http\Request;
use \Paranoid\Framework\Http\Kernel;

define('BASE_PATH', dirname(__DIR__));

require_once BASE_PATH . '/vendor/autoload.php';

$container = require BASE_PATH . '/config/services.php';
//dd($container);

// bootstrapping
require  BASE_PATH.'/bootstrap/bootstrap.php';

/*
$eventDispatcher = $container->get(\Paranoid\Framework\EventDispatcher\EventDispatcher::class);
$eventDispatcher
    ->addListener(
        \Paranoid\Framework\Http\Event\ResponseEvent::class,
        new \App\EventListener\InternalErrorListener(),
    )
    ->addListener(
    \Paranoid\Framework\Http\Event\ResponseEvent::class,
    new \App\EventListener\ContentLengthListener(),
);
*/

$request = Request::createFromGlobals();
//dd($request);

$kernel = $container->get(Kernel::class);

// send response (string of content)
$response = $kernel->handle($request);

$response->send();

$kernel->terminate($request, $response);