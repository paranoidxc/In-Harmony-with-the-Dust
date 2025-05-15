<?php
class RdsController extends GController
{
  public function actionGet()
  {
    echo "get" . date("Y-m-d H:i:s");
    $rds = Yii::app()->getRds();
    var_dump($rds);

    $key = 'abc';
    $val = $rds->get($key);
    var_dump("key:", $key, "      val:", $val);
    exit;
  }

  public function actionSet()
  {
    $key = 'abc';
    $rds = Yii::app()->getRds();
    $rds->set($key,  date("Y-m-d H:i:s"), 60 * 5);
    exit;
  }

  public function actionPush()
  {
    $LIST = "ABCDE";
    $val = "echo";
    $params = [$LIST, $val];
    $rds = Yii::app()->getRds();
    $r = $rds->executeCommand('RPUSH', $params);
    var_dump("push", $r);
  }

  public function actionPop()
  {
    $rds = Yii::app()->getRds();
    $LIST = "ABCDE";
    $params = [$LIST];

    $val = $rds->executeCommand("LPOP", $params);
    var_dump("pop", $val);
  }

  public function actionIndex()
  {
    $get = Yii::app()->createAbsoluteUrl('rds/getb', ['foo' => 'bar']);
    $set = Yii::app()->createAbsoluteUrl('rds/set', ['foo' => 'bar']);
    echo json_encode([
      'get' => $get,
      'set' => $set
    ]);
    exit;
  }
}
