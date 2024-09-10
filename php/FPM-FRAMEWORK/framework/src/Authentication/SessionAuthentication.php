<?php

namespace Paranoid\Framework\Authentication;

use Paranoid\Framework\Session\SessionInterface;

class SessionAuthentication implements SessionAuthInterface
{
    private AuthUserInterface $user;

    public function __construct(
        private AuthRepositoryInterface $authUserRepository,
        private SessionInterface $session,
    )
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
        //dd($user);
        $this->session->start();
        $this->session->set('auth_id', $user->getAuthId());

        $this->user = $user;
    }

    public function logout()
    {
    }

    public function getUser(): AuthUserInterface
    {
        return $this->user;
    }
}
