<?php
class Cntl
{
    public $router = [];
    public $get = [];
    public $post = [];

    public function render($view, $args = [], $return = FALSE)
    {
        //Log::info("VIEW FILE {$view}");
        //Log::info("router", $this->router);
        $view_file = Router::absoluteViewFilePath($view, $this->router);
        Log::info("VIEW FILE {$view_file}");
        //$view_file = $view;
        if ($return) {
            return View::render($view_file, $args, $return);
        }
        View::render($view_file, $args, $return);
    }
}
