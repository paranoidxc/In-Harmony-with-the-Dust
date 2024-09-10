<?php

namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;
use Paranoid\Framework\Session\SessionInterface;

class StartSession implements MiddlewareInterface
{

    public function __construct(private SessionInterface $session)
    {
    }

    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
        $this->session->start();

        $request->setSession($this->session);

        return $requestHandler->handle($request);
    }
}
