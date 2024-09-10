<?php

use Paranoid\Framework\Http\Response;

return [
    ['GET', '/', [\App\Controller\HomeController::class, 'index']],
    ['GET', '/posts/{id:\d+}', [\App\Controller\PostsController::class, 'show']],
    ['GET', '/posts', [\App\Controller\PostsController::class, 'create']],
    ['POST', '/posts', [\App\Controller\PostsController::class, 'store']],
    ['GET', '/register', [\App\Controller\RegisterController::class, 'index']],
    ['POST', '/register', [\App\Controller\RegisterController::class, 'register']],
    ['GET', '/hello/{name:.+}', function(string $name) {
        return new Response("Hello {$name}");
    }],
    ['GET', '/ping', function() {
        return new Response("pong ".date("Y-m-d H:i:s"));
    }],

];
