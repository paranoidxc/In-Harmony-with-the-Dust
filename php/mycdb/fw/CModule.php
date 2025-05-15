<?php
class CModule
{
  private $_basePath;
  public $preload = array();
  private $_components = array();
  private $_componentConfig = array();

  public function configure($config)
  {
    if (is_array($config)) {
      foreach ($config as $key => $value) {
        $this->$key = $value;
        $this->_componentConfig[$key] = $value;
      }
    }
  }

  public function getComponent($id, $createIfNull = true)
  {
    /*$config = [*/
    /*  'class' => 'CDbConnection',*/
    /*  //'connectionString' => 'mysql:host=192.168.1.232;port=3307;dbname=store_memberdb',*/
    /*  //'connectionString' => 'mysql:host=192.168.1.232;port=3307;dbname=clouddb',*/
    /*  'connectionString' => 'mysql:host=192.168.1.245;port=3306;dbname=test',*/
    /*  'schemaCachingDuration' => 86400,*/
    /*  'emulatePrepare' => 1,*/
    /*  'autoConnect' => '',*/
    /*  'username' => 'root',*/
    /*  'password' => '123456',*/
    /*  'charset' => 'utf8',*/
    /*];*/

    //echo "getComponent {$id}";
    if (isset($this->_components[$id])) {
      //echo "isset({$id})";
      return $this->_components[$id];
    } elseif (isset($this->_componentConfig[$id])) {
      //echo "isset(_componentConfig[{$id}])";
      //echo "003-{$id}";
      $config = $this->_componentConfig[$id];
      //var_dump("config", $config, "<BR>");

      $component = Yii::createComponent($config);
      $component->init();
      return $this->_components[$id] = $component;
    }

    /*if (isset($this->_components[$id]))*/
    /*  return $this->_components[$id];*/
    /*elseif (isset($this->_componentConfig[$id]) && $createIfNull) {*/
    /*  $config = $this->_componentConfig[$id];*/
    /*  GLMonolog::debug($config, 'config');*/
    /*  if (!isset($config['enabled']) || $config['enabled']) {*/
    /*    Yii::trace("Loading \"$id\" application component", 'system.CModule');*/
    /*    unset($config['enabled']);*/
    /*    $component = Yii::createComponent($config);*/
    /*    $component->init();*/
    /*    return $this->_components[$id] = $component;*/
    /*  }*/
    /*}*/
  }

  public function setComponents($components, $merge = true)
  {
    foreach ($components as $id => $component)
      $this->setComponent($id, $component, $merge);
  }

  public function setComponent($id, $component, $merge = true)
  {
    //var_dump($id, $component, $merge);
    if ($component === null) {
      unset($this->_components[$id]);
      return;
      /*} elseif ($component instanceof IApplicationComponent) {*/
      /*  $this->_components[$id] = $component;*/
      /**/
      /*  if (!$component->getIsInitialized()) {*/
      /*    $component->init();*/
      /*  }*/
      /**/
      /*  return;*/
    } elseif (isset($this->_components[$id])) {
      if (isset($component['class']) && get_class($this->_components[$id]) !== $component['class']) {
        unset($this->_components[$id]);
        $this->_componentConfig[$id] = $component; //we should ignore merge here
        return;
      }

      foreach ($component as $key => $value) {
        if ($key !== 'class')
          $this->_components[$id]->$key = $value;
      }
    } elseif (
      isset($this->_componentConfig[$id]['class'], $component['class'])
      && $this->_componentConfig[$id]['class'] !== $component['class']
    ) {
      $this->_componentConfig[$id] = $component; //we should ignore merge here
      return;
    }

    if (isset($this->_componentConfig[$id]) && $merge) {
      $this->_componentConfig[$id] = CMap::mergeArray($this->_componentConfig[$id], $component);
      //$this->_componentConfig[$id] = $component;
    } else {
      $this->_componentConfig[$id] = $component;
      //print_r($this->_componentConfig);
    }
  }

  protected function preloadComponents()
  {
    foreach ($this->preload as $id)
      $this->getComponent($id);
  }

  protected function preinit() {}

  protected function init() {}
}
