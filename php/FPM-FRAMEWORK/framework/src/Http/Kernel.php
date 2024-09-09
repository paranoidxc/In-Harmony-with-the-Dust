<?php

namespace Paranoid\Framework\Http;

use Exception;
use Paranoid\Framework\Routing\Router;
use Paranoid\Framework\Routing\RouterInterface;
use Psr\Container\ContainerInterface;

class Kernel
{

    private string $appEnv;

    //public function __construct(private Router $router)
    public function __construct(
        private RouterInterface $router,
        private ContainerInterface $container
    )
    {
        $this->appEnv = $container->get("APP_ENV");
    }

    public function handle(Request $request): Response
    {
        try {
            //throw new \Exception("Kernel ERR");

            [$routeHandler, $vars] = $this->router->dispatcher($request, $this->container);

            $response = call_user_func_array($routeHandler, $vars);
        } catch (\Exception $exception) {
            $response = $this->createExceptionResponse($exception);
        }
        return $response;
    }

    /**
     * @throws \Exception $exception
     */
    private function createExceptionResponse(\Exception $exception): Response
    {
        if (in_array($this->appEnv, ['dev', 'test'])) {
            throw $exception;
        }

        if ($exception instanceof HttpException) {
            return $response = new Response($exception->getMessage(), $exception->getStatusCode());
        }

        $response = new Response("Server error", Response::HTTP_INTERNAL_SERVER_ERROR);

        return $response;
    }
}
