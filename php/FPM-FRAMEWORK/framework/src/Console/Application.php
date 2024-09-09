<?php


namespace Paranoid\Framework\Console;

use Psr\Container\ContainerInterface;

class Application
{
    public function  __construct(private ContainerInterface $container)
    {

    }

    public function run(): int
    {
        $argv = $_SERVER['argv'];

        $commandName = $argv[1] ?? null;

        if (!$commandName) {
            throw new ConsoleException("A Command name must be provided");
        }
        //dd($argv);

        $command = $this->container->get($commandName);
        //dd($command);

        $args = array_slice($argv, 2);
        //dd($args);

        $options = $this->parseOptions($args);
        //dd($options);

        $status = $command->execute($options);
        //dd($status);

        return $status;
    }

    private function parseOptions(array $args): array
    {
        $options = [];
        foreach($args as $arg) {
            if (str_starts_with($arg, '--')) {
                $option = explode("=", substr($arg, '2'));
                $options[$option[0]] = $option[1] ?? true;
            }
        }

        return $options;
    }
}