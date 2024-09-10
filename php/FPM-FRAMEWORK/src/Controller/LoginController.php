<?php

namespace App\Controller;

use Paranoid\Framework\Authentication\SessionAuthentication;
use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\RedirectResponse;
use Paranoid\Framework\Http\Response;

class LoginController extends AbstractController
{
    public function __construct(private SessionAuthentication $authComponent)
    {
    }

    public function index(): Response
    {
        return $this->render('login.html.twig');
    }


    public function login(): Response
    {
        $userIsAuthenticated =  $this->authComponent->authenticate(
            $this->request->input("username"),
            $this->request->input("password"),
        );

        if (!$userIsAuthenticated) {
            $this->request->getSession()->setFlash('error', 'bad creds');
            return new RedirectResponse("/login");
        }

        $user = $this->authComponent->getUser();
        $this->request->getSession()->setFlash('success', 'you are now login');
        return new RedirectResponse("/dashboard");
    }

    public function logout(): Response
    {
        $this->authComponent->logout();

        $this->request->getSession()->setFlash('success', 'bye.. see you soon!');

        return new RedirectResponse("/login");
    }
}
