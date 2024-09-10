<?php

namespace App\Controller;

use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\Response;

class DashboardController extends AbstractController
{
    public function index(): Response
    {
        return $this->render('dashboard.html.twig');
    }

}
