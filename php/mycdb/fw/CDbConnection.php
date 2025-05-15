<?php

class CDbConnection
{
  public $connectionString;
  public $username = '';
  public $password = '';
  public $schemaCachingDuration = 0;
  public $schemaCachingExclude = array();
  public $schemaCacheID = 'cache';
  public $queryCachingDuration = 0;
  public $queryCachingDependency;

  public $queryCachingCount = 0;

  public $queryCacheID = 'cache';
  public $autoConnect = true;

  public $charset;
  public $emulatePrepare;
  public $enableParamLogging = false;

  public $enableProfiling = false;

  public $tablePrefix;
  public $initSQLs;
  public $driverMap = array(
    'pgsql' => 'CPgsqlSchema',    // PostgreSQL
    'mysqli' => 'CMysqlSchema',   // MySQL
    'mysql' => 'CMysqlSchema',    // MySQL
    'sqlite' => 'CSqliteSchema',  // sqlite 3
    'sqlite2' => 'CSqliteSchema', // sqlite 2
    'mssql' => 'CMssqlSchema',    // Mssql driver on windows hosts
    'dblib' => 'CMssqlSchema',    // dblib drivers on linux (and maybe others os) hosts
    'sqlsrv' => 'CMssqlSchema',   // Mssql
    'oci' => 'COciSchema',        // Oracle driver
  );

  public $pdoClass = 'PDO';

  private $_attributes = array();
  private $_active = false;
  private $_pdo;
  private $_transaction;
  private $_schema;

  public function init()
  {
    //parent::init();
    //if ($this->autoConnect)
    $this->setActive(true);
  }

  public function getActive()
  {
    return $this->_active;
  }

  public function setActive($value)
  {
    if ($value != $this->_active) {
      if ($value) {
        $this->open();
      } else {
      }
      //$this->close();
    }
  }

  protected function open()
  {
    if ($this->_pdo === null) {
      if (empty($this->connectionString)) {
        //throw new CDbException('CDbConnection.connectionString cannot be empty.');
        throw new Exception('CDbConnection.connectionString cannot be empty.');
      }
      try {
        //Yii::trace('Opening DB connection', 'system.db.CDbConnection');
        $this->_pdo = $this->createPdoInstance();
        $this->initConnection($this->_pdo);
        $this->_active = true;
      } catch (PDOException $e) {
        //GLMonolog::debug($this->connectionString, "error cdbconnection");
        /*if (YII_DEBUG) {*/
        /*  throw new CDbException('CDbConnection failed to open the DB connection: ' .*/
        /*    $e->getMessage(), (int)$e->getCode(), $e->errorInfo);*/
        /*} else {*/
        /*  Yii::log($e->getMessage(), CLogger::LEVEL_ERROR, 'exception.CDbException');*/
        /*  throw new CDbException('CDbConnection failed to open the DB connection.', (int)$e->getCode(), $e->errorInfo);*/
        /*}*/
      }
    }
  }

  public function getPdoInstance()
  {
    return $this->_pdo;
  }

  protected function createPdoInstance()
  {
    $pdoClass = $this->pdoClass;
    if (($pos = strpos($this->connectionString, ':')) !== false) {
      $driver = strtolower(substr($this->connectionString, 0, $pos));
      if ($driver === 'mssql' || $driver === 'dblib')
        $pdoClass = 'CMssqlPdoAdapter';
      elseif ($driver === 'sqlsrv')
        $pdoClass = 'CMssqlSqlsrvPdoAdapter';
    }

    //var_dump("pdoClass", $pdoClass);
    if (!class_exists($pdoClass)) {
      var_dump("class_exists not", $pdoClass);
      exit;
      /*throw new CDbException(Yii::t(*/
      /*  'yii',*/
      /*  'CDbConnection is unable to find PDO class "{className}". Make sure PDO is installed correctly.',*/
      /*  array('{className}' => $pdoClass)*/
      /*));*/
    }

    @$instance = new $pdoClass($this->connectionString, $this->username, $this->password, $this->_attributes);

    if (!$instance) {
      //throw new CDbException(Yii::t('yii', 'CDbConnection failed to open the DB connection.'));
      throw new Exception('CDbConnection failed to open the DB connection.');
    }

    return $instance;
  }

  protected function initConnection($pdo)
  {
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    if ($this->emulatePrepare !== null && constant('PDO::ATTR_EMULATE_PREPARES'))
      $pdo->setAttribute(PDO::ATTR_EMULATE_PREPARES, $this->emulatePrepare);
    if ($this->charset !== null) {
      $driver = strtolower($pdo->getAttribute(PDO::ATTR_DRIVER_NAME));
      if (in_array($driver, array('pgsql', 'mysql', 'mysqli')))
        $pdo->exec('SET NAMES ' . $pdo->quote($this->charset));
    }
    if ($this->initSQLs !== null) {
      foreach ($this->initSQLs as $sql)
        $pdo->exec($sql);
    }
  }

  public function getCurrentTransaction()
  {
    if ($this->_transaction !== null) {
      if ($this->_transaction->getActive())
        return $this->_transaction;
    }
    return null;
  }

  public function beginTransaction()
  {
    //Yii::trace('Starting transaction', 'system.db.CDbConnection');
    $this->setActive(true);
    $this->_pdo->beginTransaction();
    return $this->_transaction = new CDbTransaction($this);
  }

  public function getSchema()
  {
    if ($this->_schema !== null) {
      return $this->_schema;
    } else {
      $driver = 'mysql';
      return $this->_schema = Yii::createComponent($this->driverMap[$driver], $this);

      /*  $driver = $this->getDriverName();*/
      /*  GLMonolog::debug($driver, "driver");*/
      /*  if (isset($this->driverMap[$driver]))*/
      /*    return $this->_schema = Yii::createComponent($this->driverMap[$driver], $this);*/
      /*  else*/
      /*    throw new CDbException(Yii::t(*/
      /*      'yii',*/
      /*      'CDbConnection does not support reading schema for {driver} database.',*/
      /*      array('{driver}' => $driver)*/
      /*    ));*/
      /*}*/
    }
  }

  public function sql($query = null)
  {
    $this->setActive(true);
    return new CDbCommand($this, $query);
  }

  public function createCommand($query = null)
  {
    $this->setActive(true);
    return new CDbCommand($this, $query);
  }


  public function getPdoType($type)
  {
    static $map = array(
      'boolean' => PDO::PARAM_BOOL,
      'integer' => PDO::PARAM_INT,
      'string' => PDO::PARAM_STR,
      'resource' => PDO::PARAM_LOB,
      'NULL' => PDO::PARAM_NULL,
    );
    return isset($map[$type]) ? $map[$type] : PDO::PARAM_STR;
  }
}
