<?php

namespace Paranoid\Framework\Controller;

use Paranoid\Framework\Http\Response;
use Psr\Container\ContainerInterface;

abstract class AbstractController
{
    protected ?ContainerInterface $container = null;

    public function setContainer(ContainerInterface $container): void
    {
        $this->container = $container;
    }

    public function render(string $template, array $parameters = [], Response $resp = null)
    {
        $content = $this->container->get('twig')->render($template, $parameters);

        $resp ??= new Response();

        $resp->setContent($content);

        return $resp;
    }
}
