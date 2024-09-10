<?php


namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;

class Success implements MiddlewareInterface
{
    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
            return new Response("it's works", 200);
    }
}
