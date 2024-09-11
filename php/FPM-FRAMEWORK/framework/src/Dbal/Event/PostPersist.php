<?php

namespace Paranoid\Framework\Dbal\Event;

use Paranoid\Framework\Dbal\Entity;
use Paranoid\Framework\EventDispatcher\Event;

class PostPersist extends Event
{
    public function __construct(private Entity $subject)
    {
        //dd($subject);
    }
}
