<?php

namespace Paranoid\Framework\Routing;

use FastRoute\Dispatcher;
use Paranoid\Framework\Http\HttpException;
use Paranoid\Framework\Http\HttpRequestMethodException;
use Paranoid\Framework\Http\Request;

class Router implements RouterInterface
{
    public function dispatcher(Request $request): array
    {
        $routeInfo = $this->extractRouteInfo($request);
        [$handler, $vars] = $routeInfo;

        if (is_array($handler)) {
            [$controller, $method] = $handler;
            $handler = [new $controller, $method];
        }

        return [$handler, $vars];
    }

    private function extractRouteInfo(Request $request)
    {
        // create a dispatcher
        $dispatcher = \FastRoute\simpleDispatcher(function (\FastRoute\ConfigureRoutes $r) {
            // $r->addRoute('GET', '/users', 'get_all_users_handler');
            // {id} must be a number (\d+)
            // $r->addRoute('GET', '/user/{id:\d+}', 'get_user_handler');
            // The /{title} suffix is optional
            // $r->addRoute('GET', '/articles/{id:\d+}[/{title}]', 'get_article_handler');
            /*
            $r->addRoute('GET', '/', function () {
                $content = '<h1>Hello World</h1>';
                return new Response($content);
            });
            */

            /*
            $r->addRoute('GET', '/posts/{id:\d+}', function ($routeParams) {
                $content = "<h1>this is post {$routeParams['id']}</h1>";
                return new Response($content);
            });
            */

            $routes = include BASE_PATH . "/routes/web.php";
            foreach ($routes as $route) {
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
