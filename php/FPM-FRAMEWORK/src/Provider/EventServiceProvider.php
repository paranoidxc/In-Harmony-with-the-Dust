<?php
namespace App\Provider;

use App\EventListener\ContentLengthListener;
use App\EventListener\InternalErrorListener;
use Paranoid\Framework\Dbal\Event\PostPersist;
use Paranoid\Framework\EventDispatcher\EventDispatcher;
use Paranoid\Framework\Http\Event\ResponseEvent;
use Paranoid\Framework\ServiceProvider\ServiceProviderInterface;

class EventServiceProvider implements ServiceProviderInterface
{
    private array $listen = [
        ResponseEvent::class => [
            InternalErrorListener::class,
            ContentLengthListener::class,
        ],
        PostPersist::class => [],
    ];

    public function __construct(private EventDispatcher $eventDispatcher)
    {
    }

    public function register(): void
    {
        foreach ($this->listen as $eventName => $listeners) {
            foreach (array_unique($listeners) as $listener) {
                $this->eventDispatcher->addListener($eventName, new $listener);
            }
        }
    }
}
