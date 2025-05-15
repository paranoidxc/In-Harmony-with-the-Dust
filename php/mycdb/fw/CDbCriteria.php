<?php
class CDbCriteria
{
  const PARAM_PREFIX = ':ycp';
  public static $paramCount = 0;

  public $select = '*';
  public $distinct = false;
  public $condition = '';

  public $params = array();

  public $limit = -1;
  public $offset = -1;
  public $order = '';
  public $group = '';
  public $join = '';
  public $having = '';

  public $with;
  public $alias;
  public $together;
  public $index;

  public function addCondition($condition, $operator = 'AND')
  {
    if (is_array($condition)) {
      if ($condition === array())
        return $this;
      $condition = '(' . implode(') ' . $operator . ' (', $condition) . ')';
    }
    if ($this->condition === '')
      $this->condition = $condition;
    else
      $this->condition = '(' . $this->condition . ') ' . $operator . ' (' . $condition . ')';

    return $this;
  }

  public function addSearchCondition($column, $keyword, $escape = true, $operator = 'AND', $like = 'LIKE')
  {
    if ($keyword == '')
      return $this;
    if ($escape)
      $keyword = '%' . strtr($keyword, array('%' => '\%', '_' => '\_', '\\' => '\\\\')) . '%';
    $condition = $column . " $like " . self::PARAM_PREFIX . self::$paramCount;
    $this->params[self::PARAM_PREFIX . self::$paramCount++] = $keyword;
    return $this->addCondition($condition, $operator);
  }

  public function addInCondition($column, $values, $operator = 'AND')
  {
    if (($n = count($values)) < 1)
      $condition = '0=1'; // 0=1 is used because in MSSQL value alone can't be used in WHERE
    elseif ($n === 1) {
      $value = reset($values);
      if ($value === null)
        $condition = $column . ' IS NULL';
      else {
        $condition = $column . '=' . self::PARAM_PREFIX . self::$paramCount;
        $this->params[self::PARAM_PREFIX . self::$paramCount++] = $value;
      }
    } else {
      $params = array();
      foreach ($values as $value) {
        $params[] = self::PARAM_PREFIX . self::$paramCount;
        $this->params[self::PARAM_PREFIX . self::$paramCount++] = $value;
      }
      $condition = $column . ' IN (' . implode(', ', $params) . ')';
    }
    return $this->addCondition($condition, $operator);
  }

  public function addNotInCondition($column, $values, $operator = 'AND')
  {
    if (($n = count($values)) < 1)
      return $this;
    if ($n === 1) {
      $value = reset($values);
      if ($value === null)
        $condition = $column . ' IS NOT NULL';
      else {
        $condition = $column . '!=' . self::PARAM_PREFIX . self::$paramCount;
        $this->params[self::PARAM_PREFIX . self::$paramCount++] = $value;
      }
    } else {
      $params = array();
      foreach ($values as $value) {
        $params[] = self::PARAM_PREFIX . self::$paramCount;
        $this->params[self::PARAM_PREFIX . self::$paramCount++] = $value;
      }
      $condition = $column . ' NOT IN (' . implode(', ', $params) . ')';
    }
    return $this->addCondition($condition, $operator);
  }

  public function addBetweenCondition($column, $valueStart, $valueEnd, $operator = 'AND')
  {
    if ($valueStart === '' || $valueEnd === '')
      return $this;

    $paramStart = self::PARAM_PREFIX . self::$paramCount++;
    $paramEnd = self::PARAM_PREFIX . self::$paramCount++;
    $this->params[$paramStart] = $valueStart;
    $this->params[$paramEnd] = $valueEnd;
    $condition = "$column BETWEEN $paramStart AND $paramEnd";

    return $this->addCondition($condition, $operator);
  }

  public function compare($column, $value, $partialMatch = false, $operator = 'AND', $escape = true)
  {
    if (is_array($value)) {
      if ($value === array())
        return $this;
      return $this->addInCondition($column, $value, $operator);
    } else
      $value = "$value";


    if (preg_match('/^(?:\s*(<>|<=|>=|<|>|=))?(.*)$/', $value, $matches)) {
      $value = $matches[2];
      $op = $matches[1];
    } else
      $op = '';

    if ($value === '')
      return $this;

    if ($partialMatch) {
      if ($op === '')
        return $this->addSearchCondition($column, $value, $escape, $operator);
      if ($op === '<>')
        return $this->addSearchCondition($column, $value, $escape, $operator, 'NOT LIKE');
    } elseif ($op === '')
      $op = '=';

    $this->addCondition($column . $op . self::PARAM_PREFIX . self::$paramCount, $operator);

    $this->params[self::PARAM_PREFIX . self::$paramCount++] = $value;

    return $this;
  }

  public function toArray()
  {
    $result = array();
    foreach (array('select', 'condition', 'params', 'limit', 'offset', 'order', 'group', 'join', 'having', 'distinct', 'scopes', 'with', 'alias', 'index', 'together') as $name)
      $result[$name] = $this->$name;
    return $result;
  }
}
