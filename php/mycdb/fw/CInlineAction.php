<?php

class CInlineAction extends CAction
{
  /**
   * Runs the action.
   * The action method defined in the controller is invoked.
   * This method is required by {@link CAction}.
   */
  public function run()
  {
    $method = 'action' . $this->getId();
    $this->getController()->$method();
  }

  /**
   * Runs the action with the supplied request parameters.
   * This method is internally called by {@link CController::runAction()}.
   * @param array $params the request parameters (name=>value)
   * @return boolean whether the request parameters are valid
   * @since 1.1.7
   */
  public function runWithParams($params)
  {
    $methodName = 'action' . $this->getId();
    $controller = $this->getController();
    $method = new ReflectionMethod($controller, $methodName);
    if ($method->getNumberOfParameters() > 0)
      return $this->runWithParamsInternal($controller, $method, $params);
    else
      return $controller->$methodName();
  }
}
