<?php

namespace App\Entity;

class Post
{
    public function __construct(
        private ?int               $id,
        private string             $title,
        private string             $body,
        private \DateTimeImmutable $createdAt,
    )
    {

    }

    public static function create(
        string              $title,
        string              $body,
        ?int                $id = null,
        ?\DateTimeImmutable $createdAt = null,
    ): Post
    {
        return new self($id, $title, $body, $createdAt?? new \DateTimeImmutable());
    }

    /**
     * @return int|null
     */
    public function getId(): ?int
    {
        return $this->id;
    }

    /**
     * @param int|null $id
     */
    public function setId(?int $id): void
    {
        $this->id = $id;
    }

    /**
     * @return string
     */
    public function getTitle(): string
    {
        return $this->title;
    }

    /**
     * @param string $title
     */
    public function setTitle(string $title): void
    {
        $this->title = $title;
    }

    /**
     * @return string
     */
    public function getBody(): string
    {
        return $this->body;
    }

    /**
     * @param string $body
     */
    public function setBody(string $body): void
    {
        $this->body = $body;
    }

    /**
     * @return \DateTimeImmutable
     */
    public function getCreatedAt(): \DateTimeImmutable
    {
        return $this->createdAt;
    }

    /**
     * @param \DateTimeImmutable $createdAt
     */
    public function setCreatedAt(\DateTimeImmutable $createdAt): void
    {
        $this->createdAt = $createdAt;
    }



}