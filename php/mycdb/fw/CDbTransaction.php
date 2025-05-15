<?php

class CDbTransaction
{
  private $_connection = null;
  private $_active;

  public function __construct(CDbConnection $connection)
  {
    $this->_connection = $connection;
    $this->_active = true;
  }

  public function commit()
  {
    if ($this->_active && $this->_connection->getActive()) {
      //Yii::trace('Committing transaction', 'system.db.CDbTransaction');
      $this->_connection->getPdoInstance()->commit();
      $this->_active = false;
    } else {
      //throw new CDbException(Yii::t('yii', 'CDbTransaction is inactive and cannot perform commit or roll back operations.'));
      throw new Exception("commit");
    }
  }

  public function rollback()
  {
    if ($this->_active && $this->_connection->getActive()) {
      //Yii::trace('Rolling back transaction', 'system.db.CDbTransaction');
      $this->_connection->getPdoInstance()->rollBack();
      $this->_active = false;
    } else {
      //throw new CDbException(Yii::t('yii', 'CDbTransaction is inactive and cannot perform commit or roll back operations.'));
      throw new Exception("rollback");
    }
  }

  public function getConnection()
  {
    return $this->_connection;
  }

  public function getActive()
  {
    return $this->_active;
  }


  protected function setActive($value)
  {
    $this->_active = $value;
  }
}
