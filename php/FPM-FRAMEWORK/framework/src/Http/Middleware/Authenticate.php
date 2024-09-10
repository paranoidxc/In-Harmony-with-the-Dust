<?php
namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;

class Authenticate implements MiddlewareInterface
{
    private bool $authenticated = true;

    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
        if (!$this->authenticated) {
            return new Response("Auth failed", 401);
        }

        return $requestHandler->handle($request);
    }
}
