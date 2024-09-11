<?php

namespace Paranoid\Framework\Dbal;

use Doctrine\DBAL\Connection;
use Paranoid\Framework\Dbal\Event\PostPersist;
use Paranoid\Framework\EventDispatcher\EventDispatcher;

class DataMapper
{
    public function __construct(
        private Connection $connection,
        private EventDispatcher $eventDispatcher,
    )
    {
    }

    public function getConnection(): Connection
    {
        return $this->connection;
    }

    public function save(Entity $subject): int|string|null
    {
       $this->eventDispatcher->dispatch(new PostPersist($subject));

       return $this->connection->lastInsertId();
    }



}
