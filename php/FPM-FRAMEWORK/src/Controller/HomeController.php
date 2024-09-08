<?php

namespace App\Controller;

use Paranoid\Framework\Http\Response;

class HomeController
{
    public function index(): Response
    {
        $content = '<h1>Hello World From HomeController</h1>';

        return new Response($content);
    }
}
