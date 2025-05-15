<?php

class Student extends CActiveRecord
{
  public $primaryKey = 'id';
  public $rawName = 'student';

  public static function model($className = __CLASS__): Student
  {
    return parent::model($className);
  }

  public function primaryKey()
  {
    return $this->primaryKey;
  }

  /*  public function getPrimaryKey()
  {
    return $this->primaryKey;
  }
  */

  public function tableName(): string
  {
    return $this->rawName;
  }

  public function getColumnsName()
  {
    return ['sno', 'sname', 'age'];
  }

  public function getColumn($name): mixed
  {
    $r = array_flip($this->getColumnsName());
    return isset($r[$name]) ? $name : NULL;
  }
}
