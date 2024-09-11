<?php

$providers = [
    \App\Provider\EventServiceProvider::class
];

foreach ($providers as $providerClass) {
    $providers = $container->get($providerClass);
    $providers->register();
}