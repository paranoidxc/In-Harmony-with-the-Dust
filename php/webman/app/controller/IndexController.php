<?php

namespace app\controller;

use support\Request;

class IndexController
{
  public function index(Request $request)
  {
    static $readme;
    if (!$readme) {
      $readme = file_get_contents(base_path('README.md'));
    }
    put_log("index call");
    return $readme;
  }

  public function view(Request $request)
  {
    return view('index/view', ['name' => 'webman']);
  }

  public function json(Request $request)
  {
    put_log("json call");
    return json(['code' => 0, 'msg' => 'ok ' . date("Y-m-d H:i:s")]);
  }
}
