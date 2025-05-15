<?php
class CRedisCache extends CCache
{
  public $hostname = 'localhost';
  public $port = 6379;
  public $password;
  public $database = 0;
  public $timeout = null;

  private $_socket;

  protected function connect()
  {
    $this->_socket = @stream_socket_client(
      $this->hostname . ':' . $this->port,
      $errorNumber,
      $errorDescription,
      $this->timeout ? $this->timeout : ini_get("default_socket_timeout")
    );
    if ($this->_socket) {
      if ($this->password !== null)
        $this->executeCommand('AUTH', array($this->password));
      $this->executeCommand('SELECT', array($this->database));
    } else {
      throw new Exception('Failed to connect to redis: ' . $errorDescription, (int)$errorNumber);
      //throw new CException('Failed to connect to redis: ' . $errorDescription, (int)$errorNumber);
    }
  }

  public function executeCommand($name, $params = array())
  {
    if ($this->_socket === null)
      $this->connect();

    array_unshift($params, $name);
    $command = '*' . count($params) . "\r\n";
    foreach ($params as $arg)
      $command .= '$' . strlen($arg) . "\r\n" . $arg . "\r\n";

    fwrite($this->_socket, $command);

    return $this->parseResponse(implode(' ', $params));
  }


  private function parseResponse()
  {
    if (($line = fgets($this->_socket)) === false) {
      //throw new CException('Failed reading data from redis connection socket.');
      throw new Exception('Failed reading data from redis connection socket.');
    }
    $type = $line[0];
    $line = substr($line, 1, -2);
    switch ($type) {
      case '+': // Status reply
        return true;
      case '-': // Error reply
        //throw new CException('Redis error: ' . $line);
        throw new Exception('Redis error: ' . $line);
      case ':': // Integer reply
        // no cast to int as it is in the range of a signed 64 bit integer
        return $line;
      case '$': // Bulk replies
        if ($line == '-1')
          return null;
        $length = $line + 2;
        $data = '';
        while ($length > 0) {
          if (($block = fread($this->_socket, $length)) === false) {
            //throw new CException('Failed reading data from redis connection socket.');
            throw new Exception('Failed reading data from redis connection socket.');
          }
          $data .= $block;
          $length -= (function_exists('mb_strlen') ? mb_strlen($block, '8bit') : strlen($block));
        }
        return substr($data, 0, -2);
      case '*': // Multi-bulk replies
        $count = (int)$line;
        $data = array();
        for ($i = 0; $i < $count; $i++)
          $data[] = $this->parseResponse();
        return $data;
      default:
        throw new Exception('Unable to parse data received from redis.');
        //throw new CException('Unable to parse data received from redis.');
    }
  }

  protected function getValue($key)
  {
    return $this->executeCommand('GET', array($key));
  }

  protected function getValues($keys)
  {
    $response = $this->executeCommand('MGET', $keys);
    $result = array();
    $i = 0;
    foreach ($keys as $key)
      $result[$key] = $response[$i++];
    return $result;
  }

  protected function setValue($key, $value, $expire)
  {
    if ($expire == 0)
      return (bool)$this->executeCommand('SET', array($key, $value));
    return (bool)$this->executeCommand('SETEX', array($key, $expire, $value));
  }

  protected function addValue($key, $value, $expire)
  {
    if ($expire == 0)
      return (bool)$this->executeCommand('SETNX', array($key, $value));

    if ($this->executeCommand('SETNX', array($key, $value))) {
      $this->executeCommand('EXPIRE', array($key, $expire));
      return true;
    } else
      return false;
  }

  protected function deleteValue($key)
  {
    return (bool)$this->executeCommand('DEL', array($key));
  }

  protected function flushValues()
  {
    return $this->executeCommand('FLUSHDB');
  }
}
