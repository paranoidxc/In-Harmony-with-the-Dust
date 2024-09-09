<?php

namespace Paranoid\Framework\Http;

class RedirectResponse extends Response
{
    public function __construct(private string $url = '')
    {
        parent::__construct('', 302, ['location' => $url]);
    }

    public function send(): void
    {
        header("Location: " . $this->getHeaders('location'), true, $this->getStatus());
        exit;
    }
}
