<?php

namespace Paranoid\Framework\Console\Command;

class MigrateDatabase implements CommandInterface
{
    public string $name = 'database:mirgrations:mirgrate';

    public function execute(array $params = []): int
    {
        echo "Exec".PHP_EOL;
        return 0;
    }
}
