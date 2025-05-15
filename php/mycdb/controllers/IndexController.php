<?php
class IndexController extends GController
{

  public function filters()
  {
    return [
      ['TestFilter - test']
    ];
  }

  public function actionTest()
  {
    //var_dump($this->con
    var_dump($this->getId());
    var_dump($this->getAction()->getId());
    //exit;
    $string = file_get_contents('php://input');
    echo json_encode([
      'method' => 'test',
      'input' => $string,
      'get' => $_GET,
      'post' => $_POST,

    ]);
    exit;
  }

  public function actionIndex()
  {
    $old_model = Student::model()->findByPk(12);
    $db = Yii::app()->getDb();
    $list = $db->createCommand("select * from student where sname = :sname")->bindValue(':sname', 'abc')->queryAll();
    $url = Yii::app()->createUrl('a/b', ['foo' => 'bar']);
    $abs_url = Yii::app()->createAbsoluteUrl('a/b', ['foo' => 'bar']);
    //print_r("<pre>");
    echo json_encode([
      'method' => 'index',
      'student' => $old_model->getAttributes(),
      'list_student' => $list,
      'tip' => 'okay',
      'url' => $url,
      'abs_url' => $abs_url
    ]);
    exit;
  }

  public function actionError()
  {

    $error = Yii::app()->getErrorHandler()->getError();

    $err_dump = is_object($error) ? serialize($error) : "undefine";
    $err_dump = is_array($error) ? $error : $err_dump;

    echo json_encode(['method' => 'error', 'tip' => 'handle error', 'err' => $err_dump]);
    exit;
  }
}
