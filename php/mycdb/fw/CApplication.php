<?php
class CApplication extends CModule
{
  //private $_components = array();
  //private $_componentConfig = array();


  public $name = 'FuckingTheWorld';
  private $_basePath;
  private $_id;

  public function __construct($config = null)
  {
    Yii::setApplication($this);

    // set basePath at early as possible to avoid trouble
    if (is_string($config)) {
      $config = require($config);
    }
    if (isset($config['basePath'])) {
      $this->setBasePath($config['basePath']);
      unset($config['basePath']);
    } else {
      //$this->setBasePath('protected');
    }
    //Yii::setPathOfAlias('application', $this->getBasePath());
    //Yii::setPathOfAlias('webroot', dirname($_SERVER['SCRIPT_FILENAME']));
    if (isset($config['extensionPath'])) {
      //$this->setExtensionPath($config['extensionPath']);
      //unset($config['extensionPath']);
    } else {
      //  Yii::setPathOfAlias('ext', $this->getBasePath() . DIRECTORY_SEPARATOR . 'extensions');
    }
    if (isset($config['aliases'])) {
      //$this->setAliases($config['aliases']);
      //unset($config['aliases']);
    }
    $this->preinit();
    $this->initSystemHandlers();
    $this->registerCoreComponents();
    $this->configure($config);
    /*$this->attachBehaviors($this->behaviors);*/
    $this->preloadComponents();
    $this->init();
  }

  protected function initSystemHandlers()
  {
    //if (YII_ENABLE_EXCEPTION_HANDLER)
    //set_exception_handler(array($this, 'handleException'));
    //if (YII_ENABLE_ERROR_HANDLER)
    //set_error_handler(array($this, 'handleError'), error_reporting());
  }

  public function run()
  {
    /*if ($this->hasEventHandler('onBeginRequest')) {*/
    /*  $this->onBeginRequest(new CEvent($this));*/
    /*}*/
    //register_shutdown_function(array($this, 'end'), 0, false);
    $this->processRequest(); //call subclass method
    /*if ($this->hasEventHandler('onEndRequest')) {*/
    /*  $this->onEndRequest(new CEvent($this));*/
    /*}*/
  }

  public function end($status = 0, $exit = true)
  {
    /*if ($this->hasEventHandler('onEndRequest'))*/
    /*  $this->onEndRequest(new CEvent($this));*/
    if ($exit)
      exit($status);
  }

  //public function getDb()
  public function getDb(): CDbConnection
  {
    //echo "getdb\n";
    return $this->getComponent('db');
  }

  public function getRds()
  {
    //echo "getdb\n";
    return $this->getComponent('redis');
  }


  public function getId()
  {
    if ($this->_id !== null)
      return $this->_id;
    else
      return $this->_id = sprintf('%x', crc32($this->getBasePath() . $this->name));
  }

  public function setId($id)
  {
    $this->_id = $id;
  }

  public function getErrorHandler()
  {
    return $this->getComponent('errorHandler');
  }

  public function getBasePath()
  {
    return $this->_basePath;
  }

  public function getRequest()
  {
    return $this->getComponent('request');
  }

  public function getUrlManager()
  {
    return $this->getComponent('urlManager');
  }

  protected function registerCoreComponents()
  {
    $components = array(
      /*'coreMessages' => array(*/
      /*  'class' => 'CPhpMessageSource',*/
      /*  'language' => 'en_us',*/
      /*  'basePath' => YII_PATH . DIRECTORY_SEPARATOR . 'messages',*/
      /*),*/
      'db' => array(
        'class' => 'CDbConnection',
      ),
      /*'messages' => array(*/
      /*  'class' => 'CPhpMessageSource',*/
      /*),*/
      'errorHandler' => array(
        'class' => 'CErrorHandler',
      ),
      /*'securityManager' => array(*/
      /*  'class' => 'CSecurityManager',*/
      /*),*/
      /*'statePersister' => array(*/
      /*  'class' => 'CStatePersister',*/
      /*),*/
      'urlManager' => array(
        'class' => 'CUrlManager',
      ),
      'request' => array(
        'class' => 'CHttpRequest',
      ),
      /*'format' => array(*/
      /*  'class' => 'CFormatter',*/
      /*),*/
    );

    $this->setComponents($components);
  }

  public function setBasePath($path)
  {
    if (($this->_basePath = realpath($path)) === false || !is_dir($this->_basePath)) {
      throw new Exception((
        'yii' .
        'Base path "{path}" is not a valid directory.'
        //        array('{path}' => $path)
      ));
      /*throw new CException(Yii::t(*/
      /*  'yii',*/
      /*  'Base path "{path}" is not a valid directory.',*/
      /*  array('{path}' => $path)*/
      /*));*/
    }
  }


  public function handleException($exception)
  {
    // disable error capturing to avoid recursive errors
    restore_error_handler();
    restore_exception_handler();

    $category = 'exception.' . get_class($exception);
    /*if ($exception instanceof CHttpException) {*/
    /*  $category .= '.' . $exception->statusCode;*/
    /*}*/
    // php <5.2 doesn't support string conversion auto-magically
    $message = $exception->__toString();
    if (isset($_SERVER['REQUEST_URI']))
      $message .= "\nREQUEST_URI=" . $_SERVER['REQUEST_URI'];
    if (isset($_SERVER['HTTP_REFERER']))
      $message .= "\nHTTP_REFERER=" . $_SERVER['HTTP_REFERER'];
    $message .= "\n---";
    //Yii::log($message, CLogger::LEVEL_ERROR, $category);


    /*echo "001---------------handleException<br>\n";*/
    /*var_dump("exception", $exception, "<br>");*/
    /*var_dump("message", $message, "<br>");*/


    if (($handler = $this->getErrorHandler()) !== null) {
      //$handler->handle($event);
      $handler->handle($exception);
    } else {
      $this->displayException($exception);
    }

    /*try {*/
    /*  $event = new CExceptionEvent($this, $exception);*/
    /*  $this->onException($event);*/
    /*  if (!$event->handled) {*/
    /*    // try an error handler*/
    /*    if (($handler = $this->getErrorHandler()) !== null)*/
    /*      $handler->handle($event);*/
    /*    else*/
    /*      $this->displayException($exception);*/
    /*  }*/
    /*} catch (Exception $e) {*/
    /*  $this->displayException($e);*/
    /*}*/

    /*
    try {
      $this->end(1);
    } catch (Exception $e) {
      // use the most primitive way to log error
      $msg = get_class($e) . ': ' . $e->getMessage() . ' (' . $e->getFile() . ':' . $e->getLine() . ")\n";
      $msg .= $e->getTraceAsString() . "\n";
      $msg .= "Previous exception:\n";
      $msg .= get_class($exception) . ': ' . $exception->getMessage() . ' (' . $exception->getFile() . ':' . $exception->getLine() . ")\n";
      $msg .= $exception->getTraceAsString() . "\n";
      $msg .= '$_SERVER=' . var_export($_SERVER, true);
      error_log($msg);
      exit(1);
    }
    */
  }


  public function handleError($code, $message, $file, $line)
  {
    if ($code & error_reporting()) {
      // disable error capturing to avoid recursive errors
      restore_error_handler();
      restore_exception_handler();

      $log = "$message ($file:$line)\nStack trace:\n";
      $trace = debug_backtrace();
      // skip the first 3 stacks as they do not tell the error position
      if (count($trace) > 3)
        $trace = array_slice($trace, 3);
      foreach ($trace as $i => $t) {
        if (!isset($t['file']))
          $t['file'] = 'unknown';
        if (!isset($t['line']))
          $t['line'] = 0;
        if (!isset($t['function']))
          $t['function'] = 'unknown';
        $log .= "#$i {$t['file']}({$t['line']}): ";
        if (isset($t['object']) && is_object($t['object']))
          $log .= get_class($t['object']) . '->';
        $log .= "{$t['function']}()\n";
      }
      if (isset($_SERVER['REQUEST_URI']))
        $log .= 'REQUEST_URI=' . $_SERVER['REQUEST_URI'];
      //Yii::log($log, CLogger::LEVEL_ERROR, 'php');


      echo "002---------------handleError<br>\n";
      var_dump($message);

      /*try {*/
      /*  Yii::import('CErrorEvent', true);*/
      /*  $event = new CErrorEvent($this, $code, $message, $file, $line);*/
      /*  $this->onError($event);*/
      /*  if (!$event->handled) {*/
      /*    // try an error handler*/
      /*    if (($handler = $this->getErrorHandler()) !== null)*/
      /*      $handler->handle($event);*/
      /*    else*/
      /*      $this->displayError($code, $message, $file, $line);*/
      /*  }*/
      /*} catch (Exception $e) {*/
      /*  $this->displayException($e);*/
      /*}*/

      /*try {*/
      /*  $this->end(1);*/
      /*} catch (Exception $e) {*/
      /*  // use the most primitive way to log error*/
      /*  $msg = get_class($e) . ': ' . $e->getMessage() . ' (' . $e->getFile() . ':' . $e->getLine() . ")\n";*/
      /*  $msg .= $e->getTraceAsString() . "\n";*/
      /*  $msg .= "Previous error:\n";*/
      /*  $msg .= $log . "\n";*/
      /*  $msg .= '$_SERVER=' . var_export($_SERVER, true);*/
      /*  error_log($msg);*/
      /*  exit(1);*/
      /*}*/
    }
  }


  public function displayException($exception)
  {
    if (YII_DEBUG) {
      echo '<h1>' . get_class($exception) . "</h1>\n";
      echo '<p>' . $exception->getMessage() . ' (' . $exception->getFile() . ':' . $exception->getLine() . ')</p>';
      echo '<pre>' . $exception->getTraceAsString() . '</pre>';
    } else {
      echo '<h1>' . get_class($exception) . "</h1>\n";
      echo '<p>' . $exception->getMessage() . '</p>';
    }
  }

  public function createUrl($route, $params = array(), $ampersand = '&')
  {
    return $this->getUrlManager()->createUrl($route, $params, $ampersand);
  }

  public function createAbsoluteUrl($route, $params = array(), $schema = '', $ampersand = '&')
  {
    $url = $this->createUrl($route, $params, $ampersand);
    if (strpos($url, 'http') === 0)
      return $url;
    else
      return $this->getRequest()->getHostInfo($schema) . $url;
  }
}
