<?php
class YzgUser extends CActiveRecord
{
  public $primaryKey = 'id';
  public $rawName = 'yzg_user';

  public function primaryKey()
  {
    return 'id';
  }

  public function tableName()
  {
    return 'yzg_user';
  }

  public function getColumnsName()
  {
    return ['customer_id'];
  }
}
