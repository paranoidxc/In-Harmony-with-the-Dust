<?php

use Paranoid\Framework\Routing\RouterInterface;
use Paranoid\Framework\Session\SessionInterface;

$dotenv = new \Symfony\Component\Dotenv\Dotenv();
$dotenv->load(BASE_PATH . '/.env');

$container = new \League\Container\Container();

$container->delegate(new \League\Container\ReflectionContainer(true));

$routes = include BASE_PATH . '/routes/web.php';

$appEnv = $_SERVER['APP_ENV'];

$templatesPath = BASE_PATH . '/templates';

$container->add('APP_ENV', new \League\Container\Argument\Literal\StringArgument($appEnv));

$databaseUrl = 'sqlite:///'.BASE_PATH.'/var/db.sqlite';

$container->add(
    'base-commands-namespace',
    new \League\Container\Argument\Literal\StringArgument('Paranoid\\Framework\\Console\\Command\\')
);

$container->add(
    \Paranoid\Framework\Routing\RouterInterface::class,
    \Paranoid\Framework\Routing\Router::class,
);

$container->extend(\Paranoid\Framework\Routing\RouterInterface::class);
/*
$container->extend(\Paranoid\Framework\Routing\RouterInterface::class)
    ->addMethodCall(
        'setRoutes',
        [new \League\Container\Argument\Literal\ArrayArgument($routes)]
    );
*/

$container->add(
    \Paranoid\Framework\Http\Middleware\RequestHandlerInterface::class,
    \Paranoid\Framework\Http\Middleware\RequestHandler::class,
)->addArgument($container);


$container->addShared(\Paranoid\Framework\EventDispatcher\EventDispatcher::class);

$container->add(\Paranoid\Framework\Http\Kernel::class)
    ->addArguments([
        $container,
        \Paranoid\Framework\Http\Middleware\RequestHandlerInterface::class,
        \Paranoid\Framework\EventDispatcher\EventDispatcher::class,
    ]);


$container->add(\Paranoid\Framework\Console\Application::class)
    ->addArgument($container);

$container->add(\Paranoid\Framework\Console\Kernel::class)
    ->addArgument($container)
    ->addArgument( \Paranoid\Framework\Console\Application::class);


/*
$container->addShared('filesystem-loader', \Twig\Loader\FilesystemLoader::class)
    ->addArgument(new \League\Container\Argument\Literal\StringArgument($templatesPath));

$container->addShared('twig', \Twig\Environment::class)
    ->addArgument('filesystem-loader');
*/
$container->addShared(
    \Paranoid\Framework\Session\SessionInterface::class,
    \Paranoid\Framework\Session\Session::class,
);
$container->add('template-renderer-factory', \Paranoid\Framework\Template\TwigFactory::class)
    ->addArgument(\Paranoid\Framework\Session\SessionInterface::class)
    ->addArgument(new \League\Container\Argument\Literal\StringArgument($templatesPath));
$container->addShared('twig', function() use ($container) {
    return $container->get('template-renderer-factory')->create();
});

$container->add(\Paranoid\Framework\Controller\AbstractController::class);

$container->inflector(\Paranoid\Framework\Controller\AbstractController::class)
   ->invokeMethod('setContainer', [$container]);

$container->add(\Paranoid\Framework\Dbal\ConnectionFactory::class)
    ->addArgument(
        new \League\Container\Argument\Literal\StringArgument($databaseUrl)
    );
$container->addShared(\Doctrine\DBAL\Connection::class, function() use($container): \Doctrine\DBAL\Connection{
    return $container->get(\Paranoid\Framework\Dbal\ConnectionFactory::class)->create();
});

$container->add(\Paranoid\Framework\Http\Middleware\RouterDispatch::class)
    ->addArgument(RouterInterface::class)
    ->addArgument($container);


$container->add(\Paranoid\Framework\Authentication\SessionAuthentication::class)
    ->addArguments([
        \App\Repository\UserRepository::class,
        SessionInterface::class,
    ]);

$container->add(\Paranoid\Framework\Http\Middleware\ExtractRouteInfo::class)
    ->addArgument(new \League\Container\Argument\Literal\ArrayArgument($routes));

return $container;
