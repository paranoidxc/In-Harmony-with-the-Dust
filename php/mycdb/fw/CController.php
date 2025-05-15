<?php

class CController
{
  public $defaultAction = 'index';

  private $_id;
  private $_action;
  private $_pageTitle;
  private $_cachingStack;
  private $_clips;
  private $_dynamicOutput;
  private $_pageStates;
  private $_module;

  public function __construct($id, $module = null)
  {
    $this->_id = $id;
    $this->_module = $module;
    //$this->attachBehaviors($this->behaviors());
  }

  public function init() {}

  public function getId()
  {
    return $this->_id;
  }

  public function getAction()
  {
    return $this->_action;
  }

  public function setAction($value)
  {
    $this->_action = $value;
  }


  public function getModule()
  {
    return $this->_module;
  }

  public function filters()
  {
    // https://blog.csdn.net/hudeyong926/article/details/99540471
    return array();
  }


  protected function beforeAction($action)
  {
    return true;
  }

  protected function afterAction($action) {}

  public function run($actionID)
  {
    Yii::trace(['run actionId' => $actionID]);
    if (($action = $this->createAction($actionID)) !== null) {
      /*if (($parent = $this->getModule()) === null) {*/
      /*  $parent = Yii::app();*/
      /*}*/
      //print_r($action);
      Yii::trace(['createAction' => $action]);
      $this->setAction($actionID);
      $parent = Yii::app();
      if ($parent->beforeControllerAction($this, $action)) {
        Yii::trace(['tip' => "beforeControllerAction"]);
        $this->runActionWithFilters($action, $this->filters());
        $parent->afterControllerAction($this, $action);
      }
    } else {
      Yii::trace(['tip' => "missingAction"]);
      $this->missingAction($actionID);
    }
  }


  public function runActionWithFilters($action, $filters)
  {
    if (empty($filters)) {
      $this->runAction($action);
    } else {
      $priorAction = $this->_action;
      $this->_action = $action;
      CFilterChain::create($this, $action, $filters)->run();
      $this->_action = $priorAction;
    }
  }


  public function runAction($action)
  {
    Yii::trace(['md' => "runAction", 'action' => $action]);
    $priorAction = $this->_action;
    $this->_action = $action;
    if ($this->beforeAction($action)) {
      if ($action->runWithParams($this->getActionParams()) === false)
        $this->invalidActionParams($action);
      else
        $this->afterAction($action);
    }
    $this->_action = $priorAction;
  }


  public function createAction($actionID)
  {
    if ($actionID === '')
      $actionID = $this->defaultAction;
    if (method_exists($this, 'action' . $actionID) && strcasecmp($actionID, 's')) // we have actions method 
    {
      //echo 1111;
      return new CInlineAction($this, $actionID);
    } else {
      //echo 222;
      //echo "10001-Exit;";
      //exit;
      /*$action = $this->createActionFromMa($this->actions(), $actionID, $actionID);*/
      /*if ($action !== null && !method_exists($action, 'run'))*/
      /*  throw new CException(Yii::t('yii', 'Action class {class} must implement the "run" method.', array('{class}' => get_class($action))));*/
      /*return $action;*/
    }
  }

  public function getActionParams()
  {
    return $_GET;
  }


  public function missingAction($actionID)
  {
    throw new Exception((
      'yii-' .
      'The system is unable to find the requested action "{action}".' .
      ($actionID == '' ? $this->defaultAction : $actionID)
    ));
    /*
    throw new CHttpException(404, Yii::t(
      'yii',
      'The system is unable to find the requested action "{action}".',
      array('{action}' => $actionID == '' ? $this->defaultAction : $actionID)
    ));
     */
  }

  public function createUrl($route, $params = array(), $ampersand = '&')
  {
    if ($route === '')
      $route = $this->getId() . '/' . $this->getAction()->getId();
    elseif (strpos($route, '/') === false)
      $route = $this->getId() . '/' . $route;
    if ($route[0] !== '/' && ($module = $this->getModule()) !== null)
      $route = $module->getId() . '/' . $route;
    return Yii::app()->createUrl(trim($route, '/'), $params, $ampersand);
  }

  public function redirect($url, $terminate = true, $statusCode = 302)
  {
    if (is_array($url)) {
      $route = isset($url[0]) ? $url[0] : '';
      $url = $this->createUrl($route, array_splice($url, 1));
    }
    Yii::app()->getRequest()->redirect($url, $terminate, $statusCode);
  }
}
