<?php
class CModel
{

  private $_errors = array();  // attribute name => array of errors
  private $_validators;      // validators
  private $_scenario = '';    // scenario

  public function getScenario()
  {
    return $this->_scenario;
  }

  public function setScenario($value)
  {
    $this->_scenario = $value;
  }


  public function setAttributes($values, $safeOnly = true)
  {
    if (!is_array($values)) {
      return;
    }

    //$attributes = array_flip($safeOnly ? $this->getSafeAttributeNames() : $this->attributeNames());
    $attributes = array_flip($this->getColumnsName());
    //print_r("-----------------xxxxxxxxxxxxxxxxxxxxxx");
    //print_r($attributes);
    //print_r($values);

    foreach ($values as $name => $value) {
      if (isset($attributes[$name])) {
        $this->$name = $value;
      } elseif ($safeOnly) {
        //$this->onUnsafeAttribute($name, $value);
      }
    }
  }
}
