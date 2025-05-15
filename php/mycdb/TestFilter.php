<?php
class TestFilter extends CFilter
{
  public function init() {}

  protected function preFilter($filter_chain)
  {
    Yii::trace(['msg' => "TestFilter preFilter call"]);
    //var_dump($filter_chain);
    return true;
  }
}
