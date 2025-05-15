<?php
class CDbSchema
{

  public $columnTypes = array();

  private $_tableNames = array();
  private $_tables = array();
  private $_connection;
  private $_builder;
  private $_cacheExclude = array();

  public function __construct($conn)
  {
    $this->_connection = $conn;
    foreach ($conn->schemaCachingExclude as $name)
      $this->_cacheExclude[$name] = true;
  }

  public function getDbConnection()
  {
    return $this->_connection;
  }

  public function getCommandBuilder()
  {
    if ($this->_builder !== null)
      return $this->_builder;
    else
      return $this->_builder = $this->createCommandBuilder();
  }

  protected function createCommandBuilder()
  {
    return new CDbCommandBuilder($this);
  }


  public function quoteTableName($name)
  {
    if (strpos($name, '.') === false)
      return $this->quoteSimpleTableName($name);
    $parts = explode('.', $name);
    foreach ($parts as $i => $part)
      $parts[$i] = $this->quoteSimpleTableName($part);
    return implode('.', $parts);
  }

  public function quoteSimpleTableName($name)
  {
    return "'" . $name . "'";
  }
}
