<?php

class GController extends CController
{
  public function init()
  {
    parent::init();
    //echo 'init';
  }

  protected function beforeAction($action)
  {
    //echo 'beforeAction';
    Yii::trace(['action' => $action]);
    return true;
  }
}
