<?php

defined('YII_PATH') or define('YII_PATH', dirname(__FILE__));
defined('YII_BEGIN_TIME') or define('YII_BEGIN_TIME', microtime(true));

class YiiBase
{
  private static $_app;
  private static $_logger;

  public static function createWebApplication($config = null)
  {
    return self::createApplication('CWebApplication', $config);
  }


  public static function createApplication($class, $config = null)
  {
    return new $class($config);
  }

  //public static function app($config = null)
  public static function app($config = null): CApplication
  {
    if (is_null(self::$_app)) {
      self::$_app = new CApplication($config);
    }
    return self::$_app;
  }

  public static function setApplication($app)
  {
    if (self::$_app === null || $app === null) {
      self::$_app = $app;
    } else {
      throw new Exception('Yii application can only be created once.');
      //throw new CException(Yii::t('yii', 'Yii application can only be created once.'));
    }
  }

  public static function createComponent($config): object
  {
    if (is_string($config)) {
      $type = $config;
      $config = array();
    } elseif (isset($config['class'])) {
      $type = $config['class'];
      unset($config['class']);
    } else {
      print_r($config);
      //throw new CException(Yii::t('yii', 'Object configuration must be an array containing a "class" element.'));
      throw new Exception("0999-createComponent err");
    }

    //if (!class_exists($type, false))
    // $type = Yii::import($type, true);

    if (($n = func_num_args()) > 1) {
      //var_dump("func_num_args", $n);
      $args = func_get_args();
      //var_dump("args", $args);
      if ($n === 2)
        $object = new $type($args[1]);
      elseif ($n === 3)
        $object = new $type($args[1], $args[2]);
      elseif ($n === 4)
        $object = new $type($args[1], $args[2], $args[3]);
      else {
        unset($args[0]);
        $class = new ReflectionClass($type);
        // Note: ReflectionClass::newInstanceArgs() is available for PHP 5.1.3+
        // $object=$class->newInstanceArgs($args);
        $object = call_user_func_array(array($class, 'newInstance'), $args);
      }
    } else
      $object = new $type;

    foreach ($config as $key => $value)
      $object->$key = $value;

    //var_dump("object", $object);
    return $object;
  }


  public static function trace($msg, $category = 'core'): void
  {
    if (YII_DEBUG) {
      self::log($msg, CLogger::LEVEL_TRACE, $category);
    }
  }

  public static function logFlush(): void
  {
    if (self::$_logger === null) {
      self::$_logger = new CLogger;
    }
    self::$_logger->flush();
  }

  public static function log($msg, $level = CLogger::LEVEL_INFO, $category = 'core'): void
  {
    if (self::$_logger === null) {
      self::$_logger = new CLogger;
    }
    if (YII_DEBUG && YII_TRACE_LEVEL > 0 && $level !== CLogger::LEVEL_PROFILE) {
      $traces = debug_backtrace();
      $count = 0;
      foreach ($traces as $trace) {
        /*
        if (isset($trace['file'], $trace['line']) && strpos($trace['file'], YII_PATH) !== 0) {
          $msg .= "\nin " . $trace['file'] . ' (' . $trace['line'] . ')';
          if (++$count >= YII_TRACE_LEVEL)
            break;
        }
         */
      }
    }


    self::$_logger->log($msg, $level, $category);
    self::$_logger->flush();
  }
}
