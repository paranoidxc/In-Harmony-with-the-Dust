<?php

namespace Paranoid\Framework\Routing;

use FastRoute\Dispatcher;
use Paranoid\Framework\Http\HttpException;
use Paranoid\Framework\Http\HttpRequestMethodException;
use Paranoid\Framework\Http\Request;
use Psr\Container\ContainerInterface;

class Router implements RouterInterface
{
    private array $routes = [];

    public function dispatcher(Request $request, ContainerInterface $container): array
    {
        $routeInfo = $this->extractRouteInfo($request);
        [$handler, $vars] = $routeInfo;

        if (is_array($handler)) {
            [$controllerId, $method] = $handler;
            $controller = $container->get($controllerId);
            $handler = [$controller, $method];
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
            // The /{title} suffix is optional
            // $r->addRoute('GET', '/articles/{id:\d+}[/{title}]', 'get_article_handler');
            foreach ($this->routes as $route) {
                $r->addRoute(...$route);
            }
        });

        //dd($request);
        //dd($dispatcher);
        // dispatcher a url, to obtain the route info
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
}
