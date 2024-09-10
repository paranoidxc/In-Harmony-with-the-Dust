<?php

namespace Paranoid\Framework\Authentication;

use Paranoid\Framework\Session\Session;
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

        if (!password_verify($password, $user->getPassword())) {
            return false;
        }

        // login
        $this->login($user);
        return true;
    }

    public function login(AuthUserInterface $user)
    {
        //dd($user);
        $this->session->start();
        $this->session->set(Session::AUTH_KEY, $user->getAuthId());

        $this->user = $user;
    }

    public function logout()
    {
        $this->session->remove(Session::AUTH_KEY);
    }

    public function getUser(): AuthUserInterface
    {
        return $this->user;
    }
}
