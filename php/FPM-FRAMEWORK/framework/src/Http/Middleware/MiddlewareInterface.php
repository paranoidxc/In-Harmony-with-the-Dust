<?php
namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;

interface MiddlewareInterface
{
    public function process(Request $request, RequestHandlerInterface $requestHandler): Response;
}

