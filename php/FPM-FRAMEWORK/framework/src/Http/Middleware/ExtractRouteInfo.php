<?php
namespace Paranoid\Framework\Http\Middleware;

use FastRoute\Dispatcher;
use FastRoute\RouteCollector;
use Paranoid\Framework\Http\HttpException;
use Paranoid\Framework\Http\HttpRequestMethodException;
use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;
use function FastRoute\simpleDispatcher;


class ExtractRouteInfo implements MiddlewareInterface
{
    public function __construct(private array $routes)
    {
    }

    public function process(Request $request, RequestHandlerInterface $requestHandler): Response
    {
        $dispatcher = simpleDispatcher(function (RouteCollector $r) {
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
                $request->setRouteHandler($routeInfo[1]);
                $request->setRouteHandlerArgs($routeInfo[2]);

                if (is_array($routeInfo[1]) && isset($routeInfo[1][2])) {
                    $requestHandler->injectMiddleware($routeInfo[1][2]);
                }
                break;
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
        return $requestHandler->handle($request);
    }
}
