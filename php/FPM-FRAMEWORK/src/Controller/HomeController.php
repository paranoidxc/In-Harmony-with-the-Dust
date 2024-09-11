<?php

namespace App\Controller;

use App\Widget;
use Paranoid\Framework\Controller\AbstractController;
use Paranoid\Framework\Http\JsonResponse;
use Paranoid\Framework\Http\Response;
use Paranoid\Framework\Utils\TreeChildrenIds;

class HomeController extends AbstractController
{
    public function __construct(private Widget $widget)
    {

    }

    public function index(): Response
    {
        /*
        $rs = [
            ["id"=>1, "parent_id"=> 0],
            ["id"=>2, "parent_id"=> 0],
            ["id"=>3, "parent_id"=> 1],
            ["id"=>4, "parent_id"=> 2],
            ["id"=>5, "parent_id"=> 0],
        ];

        $fids = TreeChildrenIds::run($rs,0);
        dd($fids);
        */

        return $this->render('home.html.twig');
    }

    public function json(): Response
    {
        $data = [
            "key"=>"测试 The Fucking World",
        ];
        return new JsonResponse($data);
    }

    /*
    public function index(): Response
    {
        //dd($this->container->get('twig'));
        //
        //$this->container->get('twig')->render();

        $content = "<h1>Hello {$this->widget->name} From HomeController</h1>";

        return new Response($content);
    }
    */
}
