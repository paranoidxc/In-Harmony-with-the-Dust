<?php
namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\RedirectResponse;
use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;
use Paranoid\Framework\Session\Session;
use Paranoid\Framework\Session\SessionInterface;

class Guest implements MiddlewareInterface
{

    public function __construct(private SessionInterface $session)
    {
    }

    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
        $this->session->start();
        if ($this->session->has(Session::AUTH_KEY)) {
            return new RedirectResponse("/dashboard");
        }

        return $requestHandler->handle($request);
    }
}
