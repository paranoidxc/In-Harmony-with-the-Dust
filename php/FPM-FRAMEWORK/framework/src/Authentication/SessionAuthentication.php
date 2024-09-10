<?php

namespace Paranoid\Framework\Authentication;

class SessionAuthentication implements SessionAuthInterface
{
    public function __construct(private AuthRepositoryInterface $authUserRepository)
    {
    }

    public function authenticate(string $username, string $password): bool
    {
        $user = $this->authUserRepository->findByUserName($username);
        if (!$user) {
            return false;
        }

        if (password_verify($password, $user->getPassword())) {
            // login
            $this->login($user);
            return true;
        }

        return false;
    }

    public function login(AuthUserInterface $user)
    {
    }

    public function logout()
    {
    }

    public function getUser(): AuthUserInterface
    {
    }
}
