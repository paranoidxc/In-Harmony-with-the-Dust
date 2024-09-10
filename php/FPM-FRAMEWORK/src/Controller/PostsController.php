<?php

namespace App\Controller;

use App\Entity\Post;
use App\Repository\PostMapper;
use App\Repository\PostRepository;
use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\RedirectResponse;
use Paranoid\Framework\Http\Response;
use Paranoid\Framework\Session\SessionInterface;

class PostsController extends AbstractController
{
    public function __construct(
        private PostMapper $postMapper,
        private PostRepository $postRepository,
        private SessionInterface $session,
    )
    {
    }

    public function show(int $id): Response
    {
        $post = $this->postRepository->findOrFail($id);
        //dd($post);

        return $this->render('posts.html.twig',[
            'post' => $post
        ]);
    }

    public function create(): Response
    {
        return $this->render('create-post.html.twig');
    }

    public function store(): RedirectResponse
    {
        $title = $this->request->postParams['title'];
        $body = $this->request->postParams['body'];

        $post= Post::create($title, $body);

        $this->postMapper->save($post);
        //dd($post);

        $this->session->setFlash("success", sprintf("Post %s successfully created", $title));
        return new RedirectResponse('/posts');
    }
}
