<?php
namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;
use Paranoid\Framework\Routing\RouterInterface;
use Psr\Container\ContainerInterface;

class RouterDispatch implements MiddlewareInterface
{
    public function __construct(
        private RouterInterface $router,
        private ContainerInterface $container,
    )
    {
    }

    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
        [$routeHandler, $vars] = $this->router->dispatcher($request, $this->container);

        $response = call_user_func_array($routeHandler, $vars);

        return $response;
    }
}