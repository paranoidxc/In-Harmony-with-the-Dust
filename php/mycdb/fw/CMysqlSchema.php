<?php
class CMysqlSchema extends CDbSchema
{
  public function quoteSimpleTableName($name)
  {
    return '`' . $name . '`';
  }
}
