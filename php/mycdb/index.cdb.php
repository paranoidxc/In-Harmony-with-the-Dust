<?php
defined('YII_DEBUG') or define('YII_DEBUG', true);
defined('YII_TRACE_LEVEL') or define('YII_TRACE_LEVEL', 1);


error_reporting(E_ALL);
ini_set('log_error', 1);
ini_set('display_errors', 'Off');
ini_set('log_level',  'notice');
ini_set('error_log',  './log');

//ini_set('error_log','D:\'.date('Y-m-d').'_weblog.txt');

$fw_classes = [
  "./CLogger.php",
  "./CDbTransaction.php",
  "./CMap.php",
  "./CModule.php",
  "./YiiBase.php",
  "./CModel.php",
  "./CDbCriteria.php",
  "./CActiveRecord.php",
  "./CApplication.php",
  "./CDbConnection.php",
  "./CDbSchema.php",
  "./CMysqlSchema.php",
  "./CDbCommand.php",
  "./CDbCommandBuilder.php",
  "./Yii.php",
];


$classes = [];
foreach ($fw_classes as $cls) {
  $classes[] = "." . DIRECTORY_SEPARATOR . "fw" . DIRECTORY_SEPARATOR . "{$cls}";
}

$classes = array_merge(
  $classes,
  [
    "./YzgUser.php",
    "./Student.php",
  ],
);

$config = [
  'db' => [
    'class' => 'CDbConnection',
    //'connectionString' => 'mysql:host=192.168.1.232;port=3307;dbname=store_memberdb',
    //'connectionString' => 'mysql:host=192.168.1.232;port=3307;dbname=clouddb',
    'connectionString' => 'mysql:host=192.168.1.245;port=3306;dbname=test',
    'schemaCachingDuration' => 86400,
    'emulatePrepare' => 1,
    'autoConnect' => '',
    'username' => 'root',
    'password' => '123456',
    'charset' => 'utf8',
  ],
];

foreach ($classes as $cls) {
  require_once $cls;
}

function pk()
{
  $m = new YzgUser();
  $r = $m->findByPk(85, "status = :status", [":status" => 0]);
  print_r("======================");
  print_r($r);
}

function pkList()
{
  $m = new YzgUser();
  $r = $m->findAllByPk([1, 2, 3, 85, 87], "status = :status", [":status" => 0]);
  print_r("======================");
  print_r($r);
}

function one()
{
  $c = new CDbCriteria();
  $c->addCondition("t.openid = ''");
  $c->compare("t.status", 0);
  $c->order = 't.id desc ';
  $m = new YzgUser();
  $r = $m->find($c);
  return $r;
}

function listAllOld()
{
  $c = new CDbCriteria();
  $c->addCondition("t.openid = ''");
  $c->compare("t.status", 0);
  $c->join = join(" ", [
    "LEFT JOIN yzg_user_quanxian as xx ON xx.user_id = t.id"
  ]);
  $c->order = 't.id desc ';

  $m = new YzgUser();
  $r = $m->findAll($c);
  //$r = $post->count($c);
  print_r("======================");
  print_r($r);
}

function listAll()
{
  $c = new CDbCriteria();
  $c->addCondition("t.sno = 1");
  // $c->compare("t.sno", 1);
  $c->join = join(" ", [
    "LEFT JOIN take as tk ON tk.sno = t.sno"
  ]);
  $c->order = 't.id desc ';

  $m = new Student();
  $r = $m->findAll($c);
  return $r;
}

function listLimit()
{
  $c = new CDbCriteria();
  $c->addCondition("t.openid = ''");
  //$c->addInCondition()
  $c->compare("t.status", 0);
  $c->join = join(" ", [
    "LEFT JOIN yzg_user_quanxian as xx ON xx.user_id = t.id"
  ]);
  //$c->compare("t.id", 125);
  $c->order = 't.id desc ';
  $c->offset = 1;
  $c->limit = 10;

  $m = new YzgUser();
  $r = $m->findAll($c);
  //$r = $post->count($c);
  print_r("======================");
  print_r($r);
}

function icount()
{
  $c = new CDbCriteria();
  $c->addCondition("t.openid = ''");
  //$c->addInCondition()
  $c->compare("t.status", 0);
  $c->join = join(" ", [
    "LEFT JOIN yzg_user_quanxian as xx ON xx.user_id = t.id"
  ]);
  //$c->compae("t.id", 125);
  $c->order = 't.id desc ';
  $c->offset = 1;
  $c->limit = 10;

  $m = new YzgUser();
  $r = $m->count($c);
  print_r("======================");
  print_r($r);
}

function iprint($r)
{
  print_r("====================== iprint ======================\n");
  var_dump("object", $r);
  var_dump("atts", $r->getAttributes());
}

function printModelAtts($list)
{
  foreach ($list as $obj) {
    print_r($obj->getAttributes());
  }
}

print_r("<pre>");
//pkList();
//pk();
//icount();
//$m = one();
//iprint($m);
//$m->save();

// new record
/*$new = new Student();*/
/*$new->setAttributes([*/
/*  'sno' => 2,*/
/*]);*/
/*$new->save();*/

// update record
/*$model = new Student();*/
/*$model = $model->findByPk(11);*/
/*$model->sno = 22;*/
/*$model->sname = 'abc';*/
/*$model->save();*/


//list data
//$r = listAll();
//foreach ($r as $obj) {
//print_r($obj->getAttributes());
//}
//printModelAtts(Student::model()->findAll());



$app = Yii::app($config);
/*
$db = $app->getDb();
$rows = $db->createCommand("select * from student where sname = :sname")->bindValue(':sname', 'abc')->queryAll();
print_r($rows);
$rows = $db->createCommand("select * from student where sname = :sname AND id = :id")
  ->bindValues([
    ':sname' => 'abc',
    ':id' => 12,
  ])
  ->queryAll();
print_r($rows);


$cnt = $db->createCommand("UPDATE student SET age = :age WHERE id = :id")
  ->bindValues([
    ':age' => 110,
    ':id' => 12,
  ])
  ->execute();
var_dump("update affected rows: ", $cnt);


$cnt = $db->createCommand("DELETE FROM student WHERE id = :id")
  ->bindValues([
    ':id' => 1,
  ])
  ->execute();
var_dump("xxxxxxxxxxxxxxxxxxxx-j-------------------------delete affected rows: ", $cnt);
 */

$old_model = Student::model()->findByPk(12);
var_dump("exist model: ", $old_model);

$old_model = Student::model()->findByPk(2);
var_dump("old model: ", $old_model);
if (!is_null($old_model)) {
  //$cnt = $old_model->delete();
  //var_dump("delete by model affected rows: ", $cnt);
}
/*
$trans = $db->beginTransaction();
try {
  $cnt = $db->createCommand("UPDATE student SET age = :age WHERE id = :id")
    ->bindValues([
      ':age' => 199,
      ':id' => 12,
    ])
    ->execute();
  //$trans->commit();
  $trans->rollback();
} catch (Exception $e) {
  $trans->rollback();
}

$db->createCommand("UPDATE student SET age = :age WHERE id = :id")
  ->bindValues([
    ':age' => 29,
    ':id' => 12,
  ])
  ->execute();


echo 'xxxxxxxx';
Yii::trace("XXXXXXXXXXXXXXXXX");
Yii::trace(['msg' => "BBBBBBBBBBBB"]);
Yii::logFlush();

//$sql = "select * from student";
//getAttributes

//$sql = "update student set age = 4 where id = 12";
//$row = $db->createCommand($sql)->execute();
*/

Yii::trace(['msg' => "okay"]);
Yii::logFlush();

echo "\nokay";
