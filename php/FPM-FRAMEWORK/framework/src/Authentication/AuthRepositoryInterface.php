<?php

namespace Paranoid\Framework\Authentication;

interface AuthRepositoryInterface
{
    public function findByUsername(string $user): AuthUserInterface;
}