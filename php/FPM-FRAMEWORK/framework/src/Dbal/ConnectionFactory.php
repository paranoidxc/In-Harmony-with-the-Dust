<?php

namespace Paranoid\Framework\Dbal;

use Doctrine\DBAL\Connection;
use Doctrine\DBAL\DriverManager;
use SebastianBergmann\CodeCoverage\Driver\Driver;

class ConnectionFactory
{
    public function __construct(private string $databaseUrl)
    {

    }

    public function create():Connection
    {
        return DriverManager::getConnection(['url'=>$this->databaseUrl]);
    }
}