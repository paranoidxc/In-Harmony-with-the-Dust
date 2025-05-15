<?php
abstract class CActiveRecord extends CModel
{
  public static $db;

  private static $_models = array();      // class name => model
  private $_new = false;            // whether this instance is new or not
  private $_attributes = array();        // attribute name => attribute value
  private $_pk;                // old primary key value
  private $_alias = 't';            // the table alias being used for query

  abstract public function tableName();

  public static function model($className = __CLASS__)
  {
    if (isset(self::$_models[$className]))
      return self::$_models[$className];
    else {
      $model = self::$_models[$className] = new $className(null);
      //$model->attachBehaviors($model->behaviors());
      return $model;
    }
  }

  public function __construct($scenario = 'insert')
  {
    if ($scenario === null) // internally used by populateRecord() and model()
      return;

    $this->setScenario($scenario);
    $this->setIsNewRecord(true);
    /*
    $this->_attributes = $this->getMetaData()->attributeDefaults;

    $this->init();

    $this->attachBehaviors($this->behaviors());
    $this->afterConstruct();
    */
  }

  public function primaryKey()
  {
    return 'id';
  }

  public function getIsNewRecord()
  {
    return $this->_new;
  }

  public function setIsNewRecord($value)
  {
    $this->_new = $value;
  }

  public function getDbConnection()
  {
    if (self::$db !== null)
      return self::$db;
    else {
      self::$db = Yii::app()->getDb();
      return self::$db;
      /*if (self::$db instanceof CDbConnection)*/
      /*  return self::$db;*/
      /*else*/
      /*  throw new CDbException(Yii::t('yii', 'Active Record requires a "db" CDbConnection application component.'));*/
    }
  }

  public function getCommandBuilder()
  {
    return $this->getDbConnection()->getSchema()->getCommandBuilder();
  }

  public function count($condition = '', $params = array())
  {
    //Yii::trace(get_class($this) . '.count()', 'system.db.ar.CActiveRecord');
    $builder = $this->getCommandBuilder();
    //$this->beforeCount();
    $criteria = $builder->createCriteria($condition, $params);
    //$this->applyScopes($criteria);

    if (empty($criteria->with))
      return $builder->createCountCommand($this->getTableSchema(), $criteria)->queryScalar();
    else {
      //$finder = $this->getActiveFinder($criteria->with);
      //return $finder->count($criteria);
    }
  }


  public function find($condition = '', $params = array())
  {
    //Yii::trace(get_class($this) . '.find()', 'system.db.ar.CActiveRecord');
    $criteria = $this->getCommandBuilder()->createCriteria($condition, $params);
    return $this->query($criteria);
  }

  public function findAll($condition = '', $params = array())
  {
    $criteria = $this->getCommandBuilder()->createCriteria($condition, $params);
    return $this->query($criteria, true);
  }


  public function findByPk($pk, $condition = '', $params = array())
  {
    //Yii::trace(get_class($this) . '.findByPk()', 'system.db.ar.CActiveRecord');
    $prefix = $this->getTableAlias(true) . '.';
    $criteria = $this->getCommandBuilder()->createPkCriteria($this->getTableSchema(), $pk, $condition, $params, $prefix);
    return $this->query($criteria);
  }

  public function findAllByPk($pk, $condition = '', $params = array())
  {
    //Yii::trace(get_class($this) . '.findAllByPk()', 'system.db.ar.CActiveRecord');
    $prefix = $this->getTableAlias(true) . '.';
    $criteria = $this->getCommandBuilder()->createPkCriteria($this->getTableSchema(), $pk, $condition, $params, $prefix);
    return $this->query($criteria, true);
  }

  protected function query($criteria, $all = false)
  {
    if (!$all) {
      $criteria->limit = 1;
    }
    //$command = $this->getCommandBuilder()->createFindCommand($this->getTableSchema(), $criteria, $this->getTableAlias());
    //print_r($command->queryAll());
    //return $this->populateRecords($command->queryAll(), true, $criteria->index);
    //return $all ? $this->populateRecords($command->queryAll(), true, $criteria->index) : $this->populateRecord($command->queryRow());

    //$this->beforeFind();
    //$this->applyScopes($criteria);
    if (empty($criteria->with)) {
      if (!$all) {
        $criteria->limit = 1;
      }
      $command = $this->getCommandBuilder()->createFindCommand($this->getTableSchema(), $criteria, $this->getTableAlias());
      //$tmp = $all ? $command->queryAll() : $command->queryRow();
      //return $tmp;
      return $all ? $this->populateRecords($command->queryAll(), true, $criteria->index) : $this->populateRecord($command->queryRow());
    } else {
      //$finder = $this->getActiveFinder($criteria->with);
      //return $finder->query($criteria, $all);
    }
  }

  protected function instantiate($attributes)
  {
    $class = get_class($this);
    $model = new $class(null);
    return $model;
  }

  public function populateRecord($attributes, $callAfterFind = true)
  {
    //print_r("___________________________populateRecord:");
    //print_r($attributes);
    if ($attributes !== false) {
      $record = $this->instantiate($attributes);
      $record->setScenario('update');
      //var_dump("new", $record->_new);
      //$record->init();
      //$md = $record->getMetaData();
      foreach ($attributes as $name => $value) {
        /*  if (property_exists($record, $name))*/
        $record->$name = $value;
        /*  elseif (isset($md->columns[$name]))*/
        $record->_attributes[$name] = $value;
      }
      /*$record->_pk = $record->getPrimaryKey();*/
      /*$record->attachBehaviors($record->behaviors());*/
      /*if ($callAfterFind)*/
      /*  $record->afterFind();*/
      return $record;
    } else
      return null;
  }


  public function populateRecords($data, $callAfterFind = true, $index = null)
  {
    $records = array();
    foreach ($data as $attributes) {
      if (($record = $this->populateRecord($attributes, $callAfterFind)) !== null) {
        if ($index === null)
          $records[] = $record;
        else
          $records[$record->$index] = $record;
      }
    }
    return $records;
  }

  public function save($runValidation = true, $attributes = null)
  {
    //echo "model-SAVE";
    if ($this->getIsNewRecord()) {
      //echo "-- insert ";
      return $this->insert($attributes);
    } else {
      //echo "-- update ";
      return $this->update($attributes);
    }
    /*if (!$runValidation || $this->validate($attributes))*/
    /*  return $this->getIsNewRecord() ? $this->insert($attributes) : $this->update($attributes);*/
    /*else*/
    /*  return false;*/
  }

  public function delete()
  {
    if (!$this->getIsNewRecord()) {
      //Yii::trace(get_class($this) . '.delete()', 'system.db.ar.CActiveRecord');
      //if ($this->beforeDelete()) {
      if (true) {
        $result = $this->deleteByPk($this->getPrimaryKey()) > 0;
        //$this->afterDelete();
        return $result;
      } else {
        return false;
      }
    } else {
      //throw new CDbException(Yii::t('yii', 'The active record cannot be deleted because it is new.'));
    }
  }


  public function deleteByPk($pk, $condition = '', $params = array())
  {
    //Yii::trace(get_class($this) . '.deleteByPk()', 'system.db.ar.CActiveRecord');
    $builder = $this->getCommandBuilder();
    $criteria = $builder->createPkCriteria($this->getTableSchema(), $pk, $condition, $params);
    $command = $builder->createDeleteCommand($this->getTableSchema(), $criteria);
    return $command->execute();
  }

  public function insert($attributes = null)
  {
    if (!$this->getIsNewRecord()) {
      echo "getIsNewRecord";
      exit;
      //throw new CDbException(Yii::t('yii', 'The active record cannot be inserted to database because it is not new.'));
    }
    //if ($this->beforeSave()) {
    if (TRUE) {
      //Yii::trace(get_class($this) . '.insert()', 'system.db.ar.CActiveRecord');
      $builder = $this->getCommandBuilder();
      //$table = $this->getMetaData()->tableSchema;
      $table = $this->getTableSchema();
      $command = $builder->createInsertCommand($table, $this->getAttributes($attributes));
      if ($command->execute()) {
        exit;
        $primaryKey = $table->primaryKey;
        if ($table->sequenceName !== null) {
          if (is_string($primaryKey) && $this->$primaryKey === null)
            $this->$primaryKey = $builder->getLastInsertID($table);
          elseif (is_array($primaryKey)) {
            foreach ($primaryKey as $pk) {
              if ($this->$pk === null) {
                $this->$pk = $builder->getLastInsertID($table);
                break;
              }
            }
          }
        }
        $this->_pk = $this->getPrimaryKey();
        //$this->afterSave();
        $this->setIsNewRecord(false);
        $this->setScenario('update');
        return true;
      }
    }
    return false;
  }

  public function update($attributes = null)
  {
    if ($this->getIsNewRecord()) {
      echo  "IsNewRecord";
      exit;
      //throw new CDbException(Yii::t('yii', 'The active record cannot be updated because it is new.'));
    }
    //    if ($this->beforeSave()) {
    if (true) {
      //Yii::trace(get_class($this) . '.update()', 'system.db.ar.CActiveRecord');
      if ($this->_pk === null) {
        $this->_pk = $this->getPrimaryKey();
      }
      //var_dump("_pk:, ", $this->_pk);
      $this->updateByPk($this->getOldPrimaryKey(), $this->getAttributes($attributes));
      $this->_pk = $this->getPrimaryKey();
      //$this->afterSave();
      return true;
    } else
      return false;
  }

  public function updateByPk($pk, $attributes, $condition = '', $params = array())
  {
    //Yii::trace(get_class($this) . '.updateByPk()', 'system.db.ar.CActiveRecord');
    $builder = $this->getCommandBuilder();
    $table = $this->getTableSchema();
    $criteria = $builder->createPkCriteria($table, $pk, $condition, $params);
    $command = $builder->createUpdateCommand($table, $attributes, $criteria);
    return $command->execute();
  }

  public function getAttributes($names = true)
  {
    $attributes = $this->_attributes;
    $columns = $this->getColumnsName();
    //var_dump("----------columns:--------------", $columns);
    foreach ($columns as $name) {
      if (property_exists($this, $name)) {
        $attributes[$name] = $this->$name;
      } elseif ($names === true && !isset($attributes[$name])) {
        $attributes[$name] = null;
      }
    }
    /*
    foreach ($this->getMetaData()->columns as $name => $column) {
      if (property_exists($this, $name))
        $attributes[$name] = $this->$name;
      elseif ($names === true && !isset($attributes[$name]))
        $attributes[$name] = null;
    }
    */
    if (is_array($names)) {
      $attrs = array();
      foreach ($names as $name) {
        if (property_exists($this, $name))
          $attrs[$name] = $this->$name;
        else
          $attrs[$name] = isset($attributes[$name]) ? $attributes[$name] : null;
      }
      return $attrs;
    } else
      return $attributes;
  }

  public function setAttribute($name, $value)
  {
    if (property_exists($this, $name)) {
      $this->$name = $value;
      //} elseif (isset($this->getMetaData()->columns[$name])) {
      //$this->_attributes[$name] = $value;
    } else {
      $this->_attributes[$name] = $value;
      //return false;
    }
    return true;
  }

  public function getTableAlias($quote = false, $checkScopes = true)
  {
    $alias = $this->_alias;
    return $alias;

    /*if ($checkScopes && ($criteria = $this->getDbCriteria(false)) !== null && $criteria->alias != '')*/
    /*  $alias = $criteria->alias;*/
    /*else*/
    /*  $alias = $this->_alias;*/
    /*return $quote ? $this->getDbConnection()->getSchema()->quoteTableName($alias) : $alias;*/
  }

  public function getTableSchema()
  {
    return $this;
    //return $this->tableName();
    //return "yzg_user";
    //return $this->getMetaData()->tableSchema;
  }

  public function getPrimaryKey()
  {
    $table = $this;
    if (is_string($table->primaryKey))
      return $this->{$table->primaryKey};
    elseif (is_array($table->primaryKey)) {
      $values = array();
      foreach ($table->primaryKey as $name)
        $values[$name] = $this->$name;
      return $values;
    } else
      return null;
  }

  public function getOldPrimaryKey()
  {
    return $this->_pk;
  }
}
