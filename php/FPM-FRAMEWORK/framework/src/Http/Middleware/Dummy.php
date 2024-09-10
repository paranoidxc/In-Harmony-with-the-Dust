<?php

namespace Paranoid\Framework\Http\Middleware;


use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;

class Dummy implements MiddlewareInterface
{
    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
        return $requestHandler->handle($request);
    }
}
