<?php

namespace App\Controller;

use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\Response;

class LoginController extends AbstractController
{
    public function index(): Response
    {
        return $this->render('login.html.twig');
    }


    public function login(): Response
    {

    }
}
