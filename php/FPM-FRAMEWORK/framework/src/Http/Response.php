<?php
namespace Paranoid\Framework\Http;

class Response
{

    public const HTTP_INTERNAL_SERVER_ERROR = 500;

    public function __construct(
        private ?string $content = '',
        private int $status =200,
        private array $headers =[]
    ) {
        http_response_code($this->status);
    }

    public function send(): void {
        echo $this->content;
    }

    public function setContent(?string $content): void
    {
        $this->content = $content;
    }

    public function getStatus(): int
    {
        return $this->status;
    }

    public function setStatus(int $status): void
    {
        $this->status = $status;
    }

    public function getHeader(string $header): mixed
    {
        return $this->headers[$header];
    }

    public function getHeaders(): array
    {
        return $this->headers;
    }

    public function setHeader(string $key, string $val)
    {
        $this->headers[$key] = $val;
    }

    public function getContent(): ?string
    {
        return $this->content;
    }

}
