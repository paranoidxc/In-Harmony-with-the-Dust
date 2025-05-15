<?php
class CWebApplication extends CApplication
{
  public $controllerNamespace;
  public $defaultController = 'site';
  private $_controllerPath;
  private $_controller;

  public function processRequest()
  {
    //Yii::trace(['msg' => "okay"]);
    /*
    if (is_array($this->catchAllRequest) && isset($this->catchAllRequest[0])) {
      $route = $this->catchAllRequest[0];
      foreach (array_splice($this->catchAllRequest, 1) as $name => $value)
        $_GET[$name] = $value;
    } else
      $route = $this->getUrlManager()->parseUrl($this->getRequest());
    $this->runController($route);

     */

    //$this->getRequest();
    $route = $this->getUrlManager()->parseUrl($this->getRequest());
    //echo "############";
    //print_r($route);
    //$route = 
    //$route = $this->getUrlManager()->parseUrl($this->getRequest());
    $this->runController($route);
  }

  public function beforeControllerAction($controller, $action)
  {
    return true;
  }

  public function afterControllerAction($controller, $action) {}

  public function runController($route)
  {
    if (($ca = $this->createController($route)) !== null) {
      list($controller, $actionID) = $ca;
      Yii::trace(['md' => 'runController', 'ca' => $ca]);
      $oldController = $this->_controller;

      $this->_controller = $controller;
      $controller->init();
      $controller->run($actionID);
      $this->_controller = $oldController;
    } else {
      $ttt =  $route === '' ? $this->defaultController : $route;
      throw new Exception((
        '404-yii' .
        'Unable to resolve the request "{route} "' . $ttt
        //array('{route}' => $route === '' ? $this->defaultController : $route)
      ));

      /*throw new CHttpException(404, Yii::t(*/
      /*  'yii',*/
      /*  'Unable to resolve the request "{route}".',*/
      /*  array('{route}' => $route === '' ? $this->defaultController : $route)*/
      /*));*/
    }
  }

  public function getController()
  {
    return $this->_controller;
  }

  public function getControllerPath()
  {
    if ($this->_controllerPath !== null) {
      return $this->_controllerPath;
    } else {
      //echo  $this->getBasePath() . DIRECTORY_SEPARATOR . 'controllers';
      return $this->_controllerPath = $this->getBasePath() . DIRECTORY_SEPARATOR . 'controllers';
    }
  }

  public function createController($route, $owner = null)
  {
    if ($owner === null)
      $owner = $this;
    if (($route = trim($route, '/')) === '')
      $route = $owner->defaultController;
    $caseSensitive = $this->getUrlManager()->caseSensitive;

    $route .= '/';
    $controllerID = ''; // own add
    while (($pos = strpos($route, '/')) !== false) {
      $id = substr($route, 0, $pos);
      if (!preg_match('/^\w+$/', $id))
        return null;
      if (!$caseSensitive)
        $id = strtolower($id);
      $route = (string)substr($route, $pos + 1);
      if (!isset($basePath))  // first segment
      {
        //echo "first seg";
        if (isset($owner->controllerMap[$id])) {
          return array(
            Yii::createComponent($owner->controllerMap[$id], $id, $owner === $this ? null : $owner),
            $this->parseActionParams($route),
          );
        }

        $module = null;
        /*
        if (($module = $owner->getModule($id)) !== null) {
          return $this->createController($route, $module);
        }
        */

        $basePath = $owner->getControllerPath();
        Yii::trace(['basePath' => $basePath]);
        $controllerID = '';
      } else {
        $controllerID .= '/';
      }
      $className = ucfirst($id) . 'Controller';
      Yii::trace(['className' => $className]);

      $classFile = $basePath . DIRECTORY_SEPARATOR . $className . '.php';
      Yii::trace(['classFile' => $classFile]);

      if ($owner->controllerNamespace !== null)
        $className = $owner->controllerNamespace . '\\' . $className;

      if (is_file($classFile)) {
        if (!class_exists($className, false)) {
          Yii::trace(['classFile' => $classFile, 'is_file' => 1, 'class_exist' => 1]);
          require($classFile);
        }
        if (class_exists($className, false) && is_subclass_of($className, 'CController')) {
          $id[0] = strtolower($id[0]);
          return array(
            new $className($controllerID . $id, $owner === $this ? null : $owner),
            $this->parseActionParams($route),
          );
        }
        return null;
      }
      $controllerID .= $id;
      $basePath .= DIRECTORY_SEPARATOR . $id;
    }
  }


  protected function parseActionParams($pathInfo)
  {
    if (($pos = strpos($pathInfo, '/')) !== false) {
      $manager = $this->getUrlManager();
      $manager->parsePathInfo((string)substr($pathInfo, $pos + 1));
      $actionID = substr($pathInfo, 0, $pos);
      return $manager->caseSensitive ? $actionID : strtolower($actionID);
    } else
      return $pathInfo;
  }
}
