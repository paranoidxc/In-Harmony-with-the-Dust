<?php

namespace Paranoid\Framework\Console;

use Paranoid\Framework\Console\Command\CommandInterface;
use Psr\Container\ContainerInterface;

final class Kernel
{
    public function __construct(
        private ContainerInterface $container,
        private Application $application
    )
    {
    }

    public function handle(): int
    {
        $this->registerCommands();

        $status = $this->application->run();
        //dd($status);

        return $status;
    }

    private  function registerCommands(): void
    {
        $commandFiles = new \DirectoryIterator(__DIR__.'/Command');
        $namespace = $this->container->get('base-commands-namespace');
        foreach ($commandFiles as $commandFile) {
            if (!$commandFile->isFile()) {
                continue;
            }

            $command = $namespace.pathinfo($commandFile, PATHINFO_FILENAME);

            if (is_subclass_of($command, CommandInterface::class)) {
                $commandName = (new \ReflectionClass($command))->getProperty('name')->getDefaultValue();
                $this->container->add($commandName, $command);
            }
        }

        //dd($commandFiles);
    }
}