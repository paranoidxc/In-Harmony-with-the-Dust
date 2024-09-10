<?php

namespace Paranoid\Framework\Authentication;

interface AuthUserInterface
{
    public function getAuthId(): int|string;

    public function getUsername(): string;

    public function getPassword(): string;
}