<?php

namespace App\Controller;

use Paranoid\Framework\Http\Response;

class PostsController
{
    public function show(int $id): Response
    {
        $content = "<h1>This is post {$id} From PostsController</h1>";

        return new Response($content);
    }
}
