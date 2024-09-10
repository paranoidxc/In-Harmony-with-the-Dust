<?php

use Paranoid\Framework\Http\Response;

// $r->addRoute('GET', '/articles/{id:\d+}[/{title}]', 'get_article_handler');
return [
    ['GET', '/', [\App\Controller\HomeController::class, 'index']],
    ['GET', '/posts/{id:\d+}', [\App\Controller\PostsController::class, 'show']],
    ['GET', '/posts', [\App\Controller\PostsController::class, 'create']],
    ['POST', '/posts', [\App\Controller\PostsController::class, 'store']],
    ['GET', '/register', [\App\Controller\RegisterController::class, 'index',[
        \Paranoid\Framework\Http\Middleware\Guest::class,
    ]]],
    ['POST', '/register', [\App\Controller\RegisterController::class, 'register']],
    ['GET', '/login', [\App\Controller\LoginController::class, 'index',[
        \Paranoid\Framework\Http\Middleware\Guest::class,
    ]]],
    ['POST', '/login', [\App\Controller\LoginController::class, 'login']],
    ['GET', '/logout', [\App\Controller\LoginController::class, 'logout',[
        \Paranoid\Framework\Http\Middleware\Authenticate::class,
    ]]],
    ['GET', '/dashboard', [\App\Controller\DashboardController::class, 'index', [
            \Paranoid\Framework\Http\Middleware\Authenticate::class,
            \Paranoid\Framework\Http\Middleware\Dummy::class,
    ]]],
    ['GET', '/hello/{name:.+}', function(string $name) {
        return new Response("Hello {$name}");
    }],
    ['GET', '/ping', function() {
        return new Response("pong ".date("Y-m-d H:i:s"));
    }],

];
