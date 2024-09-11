<?php

namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;
use Paranoid\Framework\Http\TokenMismatchException;
use Paranoid\Framework\Session\Session;

class VerifyCsrfToken implements MiddlewareInterface
{

    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
        if (!in_array($request->getMethod(), [Request::METHOD_POST])) {
            return $requestHandler->handle($request);
        }

        $tokenFromSession = $request->getSession()->get(Session::CSRF_KEY);
        $tokenFromRequest = $request->input('_token');

        if (!hash_equals($tokenFromRequest, $tokenFromSession)) {
            $exception = new TokenMismatchException("You request could not be validated. please try again");
            $exception->setStatusCode(Response::HTTP_FORBIDDEN);
            throw $exception;
        }

        //$request->getSession()->remove(Session::CSRF_KEY);

        return $requestHandler->handle($request);
    }
}

