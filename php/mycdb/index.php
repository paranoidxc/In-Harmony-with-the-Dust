<?php
defined('YII_DEBUG') or define('YII_DEBUG', true);
defined('YII_TRACE_LEVEL') or define('YII_TRACE_LEVEL', 1);


/*error_reporting(E_ALL);*/
/*ini_set('log_error', 1);*/
/*ini_set('display_errors', 'Off');*/
/*ini_set('log_level',  'notice');*/
/*ini_set('error_log',  './log');*/

//ini_set('error_log','D:\'.date('Y-m-d').'_weblog.txt');

$fw_classes = [
  "interfaces.php",
  "CList.php",
  "CFilter.php",
  "CFilterChain.php",
  "CCache.php",
  "CRedisCache.php",
  "CErrorHandler.php",
  "CAction.php",
  "CInlineAction.php",
  "CController.php",
  "CHttpRequest.php",
  "CUrlManager.php",
  "CLogger.php",
  "CDbTransaction.php",
  "CMap.php",
  "CModule.php",
  "YiiBase.php",
  "CModel.php",
  "CDbCriteria.php",
  "CActiveRecord.php",
  "CApplication.php",
  "CWebApplication.php",
  "CDbConnection.php",
  "CDbSchema.php",
  "CMysqlSchema.php",
  "CDbCommand.php",
  "CDbCommandBuilder.php",
  "Yii.php",
];


$classes = [];
foreach ($fw_classes as $cls) {
  $classes[] = "." . DIRECTORY_SEPARATOR . "fw" . DIRECTORY_SEPARATOR . "{$cls}";
}

$classes = array_merge(
  $classes,
  [
    "./TestFilter.php",
    "./GController.php",
    "./YzgUser.php",
    "./Student.php",
  ],
);

$config = [
  //'basePath' => dirname(__FILE__) . DIRECTORY_SEPARATOR . '..',
  'basePath' => dirname(__FILE__) . DIRECTORY_SEPARATOR,
  'errorHandler' => [
    'class' => 'CErrorHandler',
    'errorAction' => 'index/error',
  ],
  'redis' => [
    'class' => 'CRedisCache',
    'hostname' => 'host.docker.internal',
    'port' => 6383,
    'database' => 0,
    'options' => STREAM_CLIENT_CONNECT,
    'hashKey' => false,
    //'serializer' => false,
  ],
  'db' => [
    'class' => 'CDbConnection',
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

Yii::createWebApplication($config);
Yii::app()->run();
