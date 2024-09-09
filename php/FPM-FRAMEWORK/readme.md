docker compose up -d

ab -n 1000 -c 100 http://localhost:8001/

docker compose exec app composer dump-autoload
docker compose exec app composer require symfony/var-dumper

    "require-dev": {
        "symfony/var-dumper": "6.3.x-dev"
    }


docker compose exec app composer require nikic/fast-route



docker compose exec app composer update


cd framework
docker compose exec app composer require --dev phpunit/phpunit:10.1.x-dev
10.1.x-dev

docker compose exec app composer require psr/container
2.0.x-dev


docker compose exec app composer dump-autoload

docker compose exec app php bin/console
docker compose exec app php bin/console database:mirgrations:mirgrate
docker compose exec app php bin/console database:mirgrations:mirgrate --up=1 --foo
