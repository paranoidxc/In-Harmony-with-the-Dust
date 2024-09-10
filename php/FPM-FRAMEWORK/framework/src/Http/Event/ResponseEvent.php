<?php
namespace Paranoid\Framework\Http\Event;

use Paranoid\Framework\EventDispatcher\Event;
use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;

class ResponseEvent extends Event
{
    public function __construct(
        private Request $request,
        private Response $response,
    )
    {
    }

    public function getRequest(): Request
    {
        return $this->request;
    }

    public function getResponse(): Response
    {
        return $this->response;
    }
}