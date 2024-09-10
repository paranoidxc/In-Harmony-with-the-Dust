<?php

namespace Paranoid\Framework\Routing;

use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\Request;
use Psr\Container\ContainerInterface;

class Router implements RouterInterface
{
    private array $routes = [];

    public function dispatcher(Request $request, ContainerInterface $container): array
    {
        $routeHandler = $request->getRouteHandler();
        $routeHandlerArgs = $request->getRouteHandlerArgs();


        if (is_array($routeHandler)) {
            [$controllerId, $method] = $routeHandler;
            $controller = $container->get($controllerId);

            if (is_subclass_of($controller, AbstractController::class)) {
                $controller->setRequest($request);
            }
            $routeHandler = [$controller, $method];
        }

        return [$routeHandler, $routeHandlerArgs];
    }

    /*
    public function dispatcher(Request $request, ContainerInterface $container): array
    {
        $routeInfo = $this->extractRouteInfo($request);
        [$handler, $vars] = $routeInfo;

        if (is_array($handler)) {
            [$controllerId, $method] = $handler;
            $controller = $container->get($controllerId);
            $handler = [$controller, $method];

            if (is_subclass_of($controller, AbstractController::class)) {
                $controller->setRequest($request);
            }
        }


        return [$handler, $vars];
    }

    public function setRoutes(array $routes): void
    {
        $this->routes = $routes;
    }

    private function extractRouteInfo(Request $request)
    {
        // create a dispatcher
        $dispatcher = \FastRoute\simpleDispatcher(function (\FastRoute\ConfigureRoutes $r) {
            foreach ($this->routes as $route) {
                $r->addRoute(...$route);
            }
        });

        $routeInfo = $dispatcher->dispatch(
            $request->getMethod(),
            $request->getPathInfo()
        );

        switch ($routeInfo[0]) {
            case Dispatcher::FOUND:
                return [$routeInfo[1], $routeInfo[2]];
            case Dispatcher::METHOD_NOT_ALLOWED:
                $allowedMethods = implode(',', $routeInfo[1]);
                $e = new HttpRequestMethodException("The allow methods are $allowedMethods");
                $e->setStatusCode(405);
                throw $e;
            default:
                $e = new HttpException("Not Found");
                $e->setStatusCode(404);
                throw $e;
        }
        //print_r($routeInfo);
        //dd($routeInfo[1]);
        //dd($routeInfo[1]);
    }
    */
}
