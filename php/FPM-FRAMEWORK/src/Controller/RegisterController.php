<?php

namespace App\Controller;

use App\Form\User\RegistrationForm;
use App\Repository\UserMapper;
use Paranoid\Framework\Authentication\SessionAuthentication;
use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\RedirectResponse;
use Paranoid\Framework\Http\Response;

class RegisterController extends AbstractController
{
    public function __construct(
        private UserMapper $userMapper,
        private SessionAuthentication $authComponent,
    )
    {
    }

    public function index(): Response
    {
        return $this->render('register.html.twig');
    }

    public function register(): Response
    {
        $form = new RegistrationForm($this->userMapper);
        $form ->setFields(
            $this->request->input('username'),
            $this->request->input('password'),
        );

        if ($form->hasValidateionErrors()) {
            foreach($form->getValidationErrors() as $error) {
                $this->request->getSession()->setFlash('error', $error);
            }

            return  new RedirectResponse('/register');
        }

        $user = $form->save();
        $this->request->getSession()->setFlash(
            'success',
            "user register ok id:".$user->getId(),
        );
        $this->authComponent->login($user);

        return new RedirectResponse('/dashboard');
    }
}