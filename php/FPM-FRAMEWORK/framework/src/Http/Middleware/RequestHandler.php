<?php
namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;
use Psr\Container\ContainerInterface;

class RequestHandler implements RequestHandlerInterface
{
    private array $middleware = [
        ExtractRouteInfo::class,
        StartSession::class,
        VerifyCsrfToken::class,
        RouterDispatch::class,
    ];

    public function __construct(
        private ContainerInterface $container
    )
    {
    }

    public function handle(Request $request): Response
    {
        if (empty($this->middleware)) {
            return new Response("empty middleware", 500);
        }

        $middlewareClass = array_shift($this->middleware);

        $middleware = $this->container->get($middlewareClass);

        $response = $middleware->process($request, $this);

        return $response;
    }

    public function injectMiddleware(array $middlerware): void
    {
        array_splice($this->middleware, 0, 0, $middlerware);
        //dd($this->middleware);
    }

}

