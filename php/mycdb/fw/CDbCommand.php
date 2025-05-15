<?php
class CDbCommand
{
  public $params = array();

  private $_connection;
  private $_text;
  private $_statement;
  private $_paramLog = array();
  private $_query;
  private $_fetchMode = array(PDO::FETCH_ASSOC);

  public function __construct(CDbConnection $connection, $query = null)
  {
    $this->_connection = $connection;
    if (is_array($query)) {
      foreach ($query as $name => $value)
        $this->$name = $value;
    } else
      $this->setText($query);
  }


  public function getConnection()
  {
    return $this->_connection;
  }


  public function queryScalar($params = array())
  {
    $result = $this->queryInternal('fetchColumn', 0, $params);
    if (is_resource($result) && get_resource_type($result) === 'stream')
      return stream_get_contents($result);
    else
      return $result;
  }


  public function queryAll($fetchAssociative = true, $params = array())
  {
    return $this->queryInternal('fetchAll', $fetchAssociative ? $this->_fetchMode : PDO::FETCH_NUM, $params);
  }


  public function queryRow($fetchAssociative = true, $params = array())
  {
    return $this->queryInternal('fetch', $fetchAssociative ? $this->_fetchMode : PDO::FETCH_NUM, $params);
  }


  private function queryInternal($method, $mode, $params = array())
  {
    $params = array_merge($this->params, $params);

    if ($this->_connection->enableParamLogging && ($pars = array_merge($this->_paramLog, $params)) !== array()) {
      $p = array();
      foreach ($pars as $name => $value)
        $p[$name] = $name . '=' . var_export($value, true);
      $par = '. Bound with ' . implode(', ', $p);
    } else {
      $par = '';
    }
    //echo "<br><br>".$this->getText()."<br><br>";

    /*if (isset(self::$SLOG) && self::$SLOG) {*/
    /*  //GLMonolog::debug($this->getText().$par);*/
    /*  //GLMonolog::debug($this->_paramLog,"_paramLog");*/
    /*  $_sql = $this->getText() . $par;*/
    /*  foreach ($this->_paramLog as $k => $v) {*/
    /*    $_sql = str_replace($k, "'{$v}'", $_sql);*/
    /*  }*/
    /*  GLMonolog::debug('Querying SQL:');*/
    /*  GLMonolog::debug($_sql);*/
    /*}*/
    /**/
    /*Yii::trace('Querying SQL: ' . $this->getText() . $par, 'system.db.CDbCommand');*/

    /*if (*/
    /*  $this->_connection->queryCachingCount > 0 && $method !== ''*/
    /*  && $this->_connection->queryCachingDuration > 0*/
    /*  && $this->_connection->queryCacheID !== false*/
    /*  && ($cache = Yii::app()->getComponent($this->_connection->queryCacheID)) !== null*/
    /*) {*/
    /*  $this->_connection->queryCachingCount--;*/
    /*  $cacheKey = 'yii:dbquery' . $this->_connection->connectionString . ':' . $this->_connection->username;*/
    /*  $cacheKey .= ':' . $this->getText() . ':' . serialize(array_merge($this->_paramLog, $params));*/
    /*  if (($result = $cache->get($cacheKey)) !== false) {*/
    /*    Yii::trace('Query result found in cache', 'system.db.CDbCommand');*/
    /*    return $result[0];*/
    /*  }*/
    /*}*/

    try {
      if ($this->_connection->enableProfiling) {
        //Yii::beginProfile('system.db.CDbCommand.query(' . $this->getText() . $par . ')', 'system.db.CDbCommand.query');
      }

      $this->prepare();
      if ($params === array())
        $this->_statement->execute();
      else
        $this->_statement->execute($params);

      if ($method === '') {
        echo "method empty";
        exit;
        //$result = new CDbDataReader($this);
      } else {
        $mode = (array)$mode;
        //var_dump("mode=============", $mode);
        //var_dump("statement===============", $this->_statement);
        call_user_func_array(array($this->_statement, 'setFetchMode'), $mode);
        $result = $this->_statement->$method();
        $this->_statement->closeCursor();
      }

      if ($this->_connection->enableProfiling) {
        //Yii::endProfile('system.db.CDbCommand.query(' . $this->getText() . $par . ')', 'system.db.CDbCommand.query');
      }

      if (isset($cache, $cacheKey))
        $cache->set($cacheKey, array($result), $this->_connection->queryCachingDuration, $this->_connection->queryCachingDependency);

      return $result;
    } catch (Exception $e) {
      if ($this->_connection->enableProfiling) {
        //Yii::endProfile('system.db.CDbCommand.query(' . $this->getText() . $par . ')', 'system.db.CDbCommand.query');
      }

      $errorInfo = $e instanceof PDOException ? $e->errorInfo : null;
      $message = $e->getMessage();
      /*Yii::log(Yii::t(*/
      /*  'yii',*/
      /*  'CDbCommand::{method}() failed: {error}. The SQL statement executed was: {sql}.',*/
      /*  array('{method}' => $method, '{error}' => $message, '{sql}' => $this->getText() . $par)*/
      /*), CLogger::LEVEL_ERROR, 'system.db.CDbCommand');*/

      /*if (YII_DEBUG)*/
      /*  $message .= '. The SQL statement executed was: ' . $this->getText() . $par;*/
      throw new Exception('yii' . 'CDbCommand failed to execute the SQL statement: {error}' . $message, (int)$e->getCode());
      /*throw new CDbException(Yii::t(*/
      /*  'yii',*/
      /*  'CDbCommand failed to execute the SQL statement: {error}',*/
      /*  array('{error}' => $message)*/
      /*), (int)$e->getCode(), $errorInfo);*/
    }
  }


  public function getText()
  {
    if ($this->_text == '' && !empty($this->_query)) {
      //$this->setText($this->buildQuery($this->_query));
      $this->setText($this->_query);
    }
    return $this->_text;
  }

  public function buildQuery($query)
  {
    echo "buildQuery not implement";
    exit;
  }

  public function prepare()
  {
    if ($this->_statement == null) {
      try {
        $this->_statement = $this->getConnection()->getPdoInstance()->prepare($this->getText());
        $this->_paramLog = array();
      } catch (Exception $e) {
        print_r("prepare");
        print_r($e);
        exit;
        /*Yii::log('Error in preparing SQL: ' . $this->getText(), CLogger::LEVEL_ERROR, 'system.db.CDbCommand');*/
        /*$errorInfo = $e instanceof PDOException ? $e->errorInfo : null;*/
        /*throw new CDbException(Yii::t(*/
        /*  'yii',*/
        /*  'CDbCommand failed to prepare the SQL statement: {error}',*/
        /*  array('{error}' => $e->getMessage())*/
        /*), (int)$e->getCode(), $errorInfo);*/
      }
    }
  }

  public function setText($value)
  {
    if ($this->_connection->tablePrefix !== null && $value != '')
      $this->_text = preg_replace('/{{(.*?)}}/', $this->_connection->tablePrefix . '\1', $value);
    else
      $this->_text = $value;
    $this->cancel();
    return $this;
  }

  public function cancel()
  {
    $this->_statement = null;
  }

  /* 
   * @see http://www.php.net/manual/en/function.PDOStatement-bindValue.php
   */
  public function bindValue($name, $value, $dataType = null)
  {
    $this->prepare();
    if ($dataType === null)
      $this->_statement->bindValue($name, $value, $this->_connection->getPdoType(gettype($value)));
    else
      $this->_statement->bindValue($name, $value, $dataType);
    $this->_paramLog[$name] = $value;
    return $this;
  }

  public function bindValues($values)
  {
    $this->prepare();
    foreach ($values as $name => $value) {
      $this->_statement->bindValue($name, $value, $this->_connection->getPdoType(gettype($value)));
      $this->_paramLog[$name] = $value;
    }
    return $this;
  }


  public function execute($params = array())
  {
    if ($this->_connection->enableParamLogging && ($pars = array_merge($this->_paramLog, $params)) !== array()) {
      $p = array();
      foreach ($pars as $name => $value)
        $p[$name] = $name . '=' . var_export($value, true);
      $par = '. Bound with ' . implode(', ', $p);
    } else
      $par = '';
    //Yii::trace('Executing SQL: ' . $this->getText() . $par, 'system.db.CDbCommand');

    /*if (isset(self::$SLOG) && self::$SLOG) {*/
    /*  $_sql = $this->getText() . $par;*/
    /*  foreach ($this->_paramLog as $k => $v) {*/
    /*    $_sql = str_replace($k, "'{$v}'", $_sql);*/
    /*  }*/
    /*  GLMonolog::debug('Executing SQL:');*/
    /*  GLMonolog::debug($_sql);*/
    /*}*/

    try {
      if ($this->_connection->enableProfiling) {
        //Yii::beginProfile('system.db.CDbCommand.execute(' . $this->getText() . $par . ')', 'system.db.CDbCommand.execute');
      }

      $this->prepare();
      if ($params === array())
        $this->_statement->execute();
      else
        $this->_statement->execute($params);
      $n = $this->_statement->rowCount();

      if ($this->_connection->enableProfiling) {
        //Yii::endProfile('system.db.CDbCommand.execute(' . $this->getText() . $par . ')', 'system.db.CDbCommand.execute');
      }

      return $n;
    } catch (Exception $e) {
      if ($this->_connection->enableProfiling) {
        //Yii::endProfile('system.db.CDbCommand.execute(' . $this->getText() . $par . ')', 'system.db.CDbCommand.execute');
      }

      $errorInfo = $e instanceof PDOException ? $e->errorInfo : null;
      $message = $e->getMessage();
      /*Yii::log(Yii::t(*/
      /*  'yii',*/
      /*  'CDbCommand::execute() failed: {error}. The SQL statement executed was: {sql}.',*/
      /*  array('{error}' => $message, '{sql}' => $this->getText() . $par)*/
      /*), CLogger::LEVEL_ERROR, 'system.db.CDbCommand');*/

      /*if (YII_DEBUG)*/
      /*  $message .= '. The SQL statement executed was: ' . $this->getText() . $par;*/

      /*throw new CDbException(Yii::t(*/
      /*  'yii',*/
      /*  'CDbCommand failed to execute the SQL statement: {error}',*/
      /*  array('{error}' => $message)*/
      /*), (int)$e->getCode(), $errorInfo);*/
    }
  }
}
