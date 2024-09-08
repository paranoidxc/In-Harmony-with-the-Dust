docker compose up -d

docker compose exec app composer dump-autoload
docker compose exec app composer require symfony/var-dumper

    "require-dev": {
        "symfony/var-dumper": "6.3.x-dev"
    }


docker compose exec app composer require nikic/fast-route

