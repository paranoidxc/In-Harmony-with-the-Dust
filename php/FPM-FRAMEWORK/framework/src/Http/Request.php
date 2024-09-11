<?php

namespace Paranoid\Framework\Http;

use Paranoid\Framework\Session\SessionInterface;

class Request
{
    private SessionInterface $session;
    private mixed $routeHandler;
    private array $routeHandlerArgs;
    private array $jsonParams = [];

    public const METHOD_GET = "GET";
    public const METHOD_POST = "POST";
    public const METHOD_PATCH = "PATCH";
    public const METHOD_DELETE = "DELETE";
    public const METHOD_UPDATE = "UPDATE";

    public function __construct(
        public readonly array $getParams,
        public readonly array $postParams,
        public readonly array $cookie,
        public readonly array $files,
        public readonly array $server
    ) {
    }

    public static function createFromGlobals(): static
    {
        return new static($_GET, $_POST, $_COOKIE, $_FILES, $_SERVER);
    }

    public function getPathInfo(): string
    {
        return strtok($this->server['REQUEST_URI'], '?');
    }

    public function getMethod(): string
    {
        return $this->server['REQUEST_METHOD'];
    }

    public function getSession(): ?SessionInterface
    {
        return $this->session ?? null;
    }

    public function setSession(SessionInterface $session): void
    {
        $this->session = $session;
    }

    public function json($key): mixed
    {
        return $this->jsonParams[$key] ?? "";
    }

    public function setJson(array $data=[])
    {
        return $this->jsonParams = $data;
    }

    public function get($key): mixed
    {
        return $this->getParams[$key];
    }

    public function input($key): mixed
    {
        return $this->postParams[$key];
    }

    public function getRouteHandler(): mixed
    {
        return $this->routeHandler;
    }

    public function setRouteHandler(mixed $routeHandler): void
    {
        $this->routeHandler = $routeHandler;
    }

    public function getRouteHandlerArgs(): array
    {
        return $this->routeHandlerArgs;
    }

    public function setRouteHandlerArgs(array $routeHandlerArgs): void
    {
        $this->routeHandlerArgs = $routeHandlerArgs;
    }

}
