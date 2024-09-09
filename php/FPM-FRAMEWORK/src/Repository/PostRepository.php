<?php

namespace App\Repository;

use App\Entity\Post;
use Doctrine\DBAL\Connection;
use Paranoid\Framework\Http\NotFouondException;

class PostRepository
{
    public function __construct(private Connection $connection)
    {

    }

    public function findById(int $id): ?Post
    {
        $queryBuilder = $this->connection->createQueryBuilder();
        $queryBuilder
            ->select('id', 'title', 'body', 'created_at')
            ->from('posts')
            ->where("id = :id")
            ->setParameter('id', $id);
        $result = $queryBuilder->executeQuery();

        $row = $result->fetchAssociative();

        if (!$row) {
            return null;
        }

        //dd($row);

        $post = Post::create(
            id: $row['id'],
            title: $row['title'],
            body: $row['body'],
            createdAt: new \DateTimeImmutable($row['created_at']),
        );

        return $post;
    }

    public function findOrFail(int $id): ?Post
    {
        $post = $this->findById($id);
        if (!$post) {
           throw new NotFouondException(sprintf("Post ID %d not found", $id));
        }

        return $post;
    }
}

