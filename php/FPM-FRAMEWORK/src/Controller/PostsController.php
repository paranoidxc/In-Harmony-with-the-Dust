<?php

namespace App\Controller;

use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\Response;

class PostsController extends AbstractController
{
    public function show(int $id): Response
    {
        return $this->render('posts.html.twig',[
            //'postId' => $id
            'postId' => "<Script>alert('xx');</Script>"
        ]);
    }

    public function create(): Response
    {
        return $this->render('create-post.html.twig');
    }
}
