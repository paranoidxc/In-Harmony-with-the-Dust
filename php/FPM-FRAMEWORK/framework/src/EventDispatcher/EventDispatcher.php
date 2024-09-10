<?php

namespace Paranoid\Framework\EventDispatcher;

use Psr\EventDispatcher\EventDispatcherInterface;
use Psr\EventDispatcher\StoppableEventInterface;

class EventDispatcher implements EventDispatcherInterface
{
    private iterable $listeners = [];

    public function dispatch(object $event): object
    {
        foreach ($this->getListenersForEvent($event) as $listener) {

            if ($event instanceof StoppableEventInterface && $event->isPropagationStopped()) {
                return $event;
            }

            $listener($event);
        }

        return $event;
    }

    public function addListener(string $eventName, callable $listener)
    {
        $this->listeners[$eventName][] = $listener;

        return $this;
    }

    public function getListenersForEvent(object $event): iterable
    {
        $eventName = get_class($event);
        if (array_key_exists($eventName, $this->listeners)) {
            return $this->listeners[$eventName];
        }

        return [];
    }
}