<?php
class CErrorHandler
{
  private $_error;
  public function init() {}
  public function handle($event)
  {
    //echo '004-CErrorHandler';
    //if ($event instanceof CExceptionEvent) {
    if ($event instanceof Exception) {
      //$this->handleException($event->exception);
      $this->handleException($event);
    } else // CErrorEvent
      $this->handleError($event);
  }


  public function getError()
  {
    return $this->_error;
  }

  protected function handleException($exception)
  {
    //echo "999-handleException";
    $app = Yii::app();
    if ($app instanceof CWebApplication) {
      //TODO 
      $this->_error = $data = $exception;
      $this->render('error', $data);
      /*if ($exception instanceof CHttpException || !YII_DEBUG) {*/
      /*} else {*/
      /*  if ($this->isAjaxRequest())*/
      /*    $app->displayException($exception);*/
      /*  else*/
      /*    $this->render('exception', $data);*/
      /*}*/
    } else {
      $app->displayException($exception);
    }
  }

  protected function handleError($event)
  {
    echo "999-handleError";
    exit;
  }


  protected function render($view, $data)
  {
    if ($view === 'error' && $this->errorAction !== null)
      Yii::app()->runController($this->errorAction);
    else {
      // additional information to be passed to view
      $data['version'] = "ver.0"; //$this->getVersionInfo();
      $data['time'] = time();
      $data['admin'] = "admin.info"; //$this->adminInfo;
      var_dump("render", $data);
      //include($this->getViewFile($view, $data['code']));
    }
  }
}
