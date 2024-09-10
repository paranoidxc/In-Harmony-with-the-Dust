<?php
namespace Paranoid\Framework\Routing;

use Paranoid\Framework\Http\Request;
use Psr\Container\ContainerInterface;

interface RouterInterface
{
    public function dispatcher(Request $request, ContainerInterface $container);

    //public function setRoutes(array $routes): void;
}
