<?php

class CDbCommandBuilder
{
  const PARAM_PREFIX = ':yp';

  private $_schema;
  private $_connection;


  public function __construct($schema)
  {
    $this->_schema = $schema;
    $this->_connection = $schema->getDbConnection();
  }

  public function createCriteria($condition = '', $params = array())
  {
    if (is_array($condition)) {
      $criteria = new CDbCriteria($condition);
    } elseif ($condition instanceof CDbCriteria) {
      $criteria = clone $condition;
    } else {
      $criteria = new CDbCriteria;
      $criteria->condition = $condition;
      $criteria->params = $params;
    }
    return $criteria;
  }

  public function createPkCriteria($table, $pk, $condition = '', $params = array(), $prefix = null)
  {
    //$this->ensureTable($table);
    $criteria = $this->createCriteria($condition, $params);
    //print_r($criteria->toArray());
    if ($criteria->alias != '') {
      $prefix = $this->_schema->quoteTableName($criteria->alias) . '.';
    }
    if (!is_array($pk)) {
      // single key
      $pk = array($pk);
    }
    //if (is_array($table->primaryKey) && !isset($pk[0]) && $pk !== array()) {
    // single composite key 
    // $pk = array($pk);
    //}
    //$condition = $this->createInCondition($table, $table->primaryKey, $pk, $prefix);
    $condition = $this->createInCondition($table, $table->primaryKey(), $pk, $prefix);
    if ($criteria->condition != '') {
      $criteria->condition = $condition . ' AND (' . $criteria->condition . ')';
    } else {
      $criteria->condition = $condition;
    }

    //var_dump("------------ condition ------------------------", $criteria->condition);
    return $criteria;
  }

  public function createInCondition($table, $columnName, $values, $prefix = null)
  {
    if (($n = count($values)) < 1)
      return '0=1';

    //$this->ensureTable($table);

    if ($prefix === null)
      $prefix = $table->rawName . '.';

    $db = $this->_connection;

    if (is_array($columnName) && count($columnName) === 1) {
      $columnName = reset($columnName);
    }

    //print_r("============== columnName :");
    //print_r($columnName);

    if (is_string($columnName)) // simple key
    {
      //if (!isset($table->columns[$columnName])) {
      /*throw new Exception(*/
      /*  'Table "{table}" does not have a column named "{column}"'*/
      /*);*/
      /*throw new CDbException(Yii::t(*/
      /*  'yii',*/
      /*  'Table "{table}" does not have a column named "{column}".',*/
      /*  array('{table}' => $table->name, '{column}' => $columnName)*/
      /*));*/
      //}
      //$column = $table->columns[$columnName];
      $column = $columnName;

      $values = array_values($values);
      /*
      foreach ($values as &$value) {
        $value = $column->typecast($value);
        if (is_string($value))
          $value = $db->quoteValue($value);
      }
      */
      if ($n === 1) {
        return $prefix . $column . ($values[0] === null ? ' IS NULL' : '=' . $values[0]);
        //return $prefix . $column->rawName . ($values[0] === null ? ' IS NULL' : '=' . $values[0]);
      } else {
        return $prefix . $column . ' IN (' . implode(', ', $values) . ')';
        //return $prefix . $column->rawName . ' IN (' . implode(', ', $values) . ')';
      }
    } elseif (is_array($columnName)) // composite key: $values=array(array('pk1'=>'v1','pk2'=>'v2'),array(...))
    {
      foreach ($columnName as $name) {
        if (!isset($table->columns[$name])) {
          throw new Exception(
            'Table "{table}" does not have a column named "{column}"'
          );
        }
        /*throw new CDbException(Yii::t(*/
        /*  'yii',*/
        /*  'Table "{table}" does not have a column named "{column}".',*/
        /*  array('{table}' => $table->name, '{column}' => $name)*/
        /*));*/

        for ($i = 0; $i < $n; ++$i) {
          if (isset($values[$i][$name])) {
            $value = $table->columns[$name]->typecast($values[$i][$name]);
            if (is_string($value))
              $values[$i][$name] = $db->quoteValue($value);
            else
              $values[$i][$name] = $value;
          } else {
            throw new Exception(
              'Table "{table}" does not have a column named "{column}"'
            );
            /*throw new CDbException(Yii::t(*/
            /*  'yii',*/
            /*  'The value for the column "{column}" is not supplied when querying the table "{table}".',*/
            /*  array('{table}' => $table->name, '{column}' => $name)*/
            /*));*/
          }
        }
      }
      if (count($values) === 1) {
        $entries = array();
        foreach ($values[0] as $name => $value)
          $entries[] = $prefix . $table->columns[$name]->rawName . ($value === null ? ' IS NULL' : '=' . $value);
        return implode(' AND ', $entries);
      }

      return $this->createCompositeInCondition($table, $values, $prefix);
    } else {
      throw new Exception('Column name must be either a string or an array.');
      //throw new CDbException(Yii::t('yii', 'Column name must be either a string or an array.'));
    }
  }

  protected function createCompositeInCondition($table, $values, $prefix)
  {
    $keyNames = array();
    foreach (array_keys($values[0]) as $name)
      $keyNames[] = $prefix . $table->columns[$name]->rawName;
    $vs = array();
    foreach ($values as $value)
      $vs[] = '(' . implode(', ', $value) . ')';
    return '(' . implode(', ', $keyNames) . ') IN (' . implode(', ', $vs) . ')';
  }

  public function createFindCommand($table, $criteria, $alias = 't')
  {
    //$this->ensureTable($table);
    $select = is_array($criteria->select) ? implode(', ', $criteria->select) : $criteria->select;
    if ($criteria->alias != '')
      $alias = $criteria->alias;
    $alias = $this->_schema->quoteTableName($alias);

    // issue 1432: need to expand * when SQL has JOIN
    if ($select === '*' && !empty($criteria->join)) {
      /*$prefix = $alias . '.';*/
      /*$select = array();*/
      /*foreach ($table->getColumnNames() as $name)*/
      /*  $select[] = $prefix . $this->_schema->quoteColumnName($name);*/
      /*$select = implode(', ', $select);*/
    }

    $table_name = $table->tableName();
    $sql = ($criteria->distinct ? 'SELECT DISTINCT' : 'SELECT') . " {$select} FROM {$table_name} $alias";
    //$sql = ($criteria->distinct ? 'SELECT DISTINCT' : 'SELECT') . " {$select} FROM {$table->rawName} $alias";
    $sql = $this->applyJoin($sql, $criteria->join);
    $sql = $this->applyCondition($sql, $criteria->condition);
    $sql = $this->applyGroup($sql, $criteria->group);
    $sql = $this->applyHaving($sql, $criteria->having);
    $sql = $this->applyOrder($sql, $criteria->order);
    $sql = $this->applyLimit($sql, $criteria->limit, $criteria->offset);
    Yii::trace(['createFindCommand.sql' => $sql]);
    $command = $this->_connection->createCommand($sql);
    $this->bindValues($command, $criteria->params);
    return $command;
  }

  public function createCountCommand($table, $criteria, $alias = 't')
  {
    //$this->ensureTable($table);
    if ($criteria->alias != '') {
      $alias = $criteria->alias;
    }
    $alias = $this->_schema->quoteTableName($alias);

    $table_name = $table->tableName();

    if (!empty($criteria->group) || !empty($criteria->having)) {
      $select = is_array($criteria->select) ? implode(', ', $criteria->select) : $criteria->select;
      if ($criteria->alias != '') {
        $alias = $criteria->alias;
      }
      $sql = ($criteria->distinct ? 'SELECT DISTINCT' : 'SELECT') . " {$select} FROM {$table_name} $alias";
      //$sql = ($criteria->distinct ? 'SELECT DISTINCT' : 'SELECT') . " {$select} FROM {$table->rawName} $alias";
      $sql = $this->applyJoin($sql, $criteria->join);
      $sql = $this->applyCondition($sql, $criteria->condition);
      $sql = $this->applyGroup($sql, $criteria->group);
      $sql = $this->applyHaving($sql, $criteria->having);
      $sql = "SELECT COUNT(*) FROM ($sql) sq";
    } else {
      if (is_string($criteria->select) && stripos($criteria->select, 'count') === 0) {
        $sql = "SELECT " . $criteria->select;
      } elseif ($criteria->distinct) {
        if (is_array($table->primaryKey())) {
          $pk = array();
          foreach ($table->primaryKey() as $key)
            $pk[] = $alias . '.' . $key;
          $pk = implode(', ', $pk);
        } else
          $pk = $alias . '.' . $table->primaryKey();
        $sql = "SELECT COUNT(DISTINCT $pk)";
      } else
        $sql = "SELECT COUNT(*)";
      //$sql .= " FROM {$table->rawName} $alias";
      $sql .= " FROM {$table_name} $alias";
      $sql = $this->applyJoin($sql, $criteria->join);
      $sql = $this->applyCondition($sql, $criteria->condition);
    }

    // Suppress binding of parameters belonging to the ORDER clause. Issue #1407.
    if ($criteria->order && $criteria->params) {
      $params1 = array();
      preg_match_all('/(:\w+)/', $sql, $params1);
      $params2 = array();
      preg_match_all('/(:\w+)/', $this->applyOrder($sql, $criteria->order), $params2);
      foreach (array_diff($params2[0], $params1[0]) as $param) {
        unset($criteria->params[$param]);
      }
    }

    // Do the same for SELECT part.
    if ($criteria->select && $criteria->params) {
      $params1 = array();
      preg_match_all('/(:\w+)/', $sql, $params1);
      $params2 = array();
      preg_match_all('/(:\w+)/', $sql . ' ' . (is_array($criteria->select) ? implode(', ', $criteria->select) : $criteria->select), $params2);
      foreach (array_diff($params2[0], $params1[0]) as $param)
        unset($criteria->params[$param]);
    }

    Yii::trace(['createCountCommand.sql' => $sql]);
    $command = $this->_connection->createCommand($sql);
    $this->bindValues($command, $criteria->params);
    return $command;
  }

  public function createDeleteCommand($table, $criteria)
  {
    //$this->ensureTable($table);
    $sql = "DELETE FROM {$table->rawName}";
    $sql = $this->applyJoin($sql, $criteria->join);
    $sql = $this->applyCondition($sql, $criteria->condition);
    $sql = $this->applyGroup($sql, $criteria->group);
    $sql = $this->applyHaving($sql, $criteria->having);
    $sql = $this->applyOrder($sql, $criteria->order);
    $sql = $this->applyLimit($sql, $criteria->limit, $criteria->offset);
    Yii::trace(['createDeleteCommand.sql' => $sql]);
    $command = $this->_connection->createCommand($sql);
    $this->bindValues($command, $criteria->params);
    return $command;
  }

  public function createInsertCommand($table, $data)
  {
    //$this->ensureTable($table);
    $fields = array();
    $values = array();
    $placeholders = array();
    $i = 0;
    foreach ($data as $name => $value) {
      $fields[] = $name;
      $placeholders[] = self::PARAM_PREFIX . $i;
      $values[self::PARAM_PREFIX . $i] = $value; //$column->typecast($value);
      $i++;
      //列存在
      //列允许空 || 值不为空
      /*if (($column = $table->getColumn($name)) !== null && ($value !== null || $column->allowNull)) {*/
      /*  $fields[] = $column->rawName;*/
      /*  if ($value instanceof CDbExpression) {*/
      /*    $placeholders[] = $value->expression;*/
      /*    foreach ($value->params as $n => $v)*/
      /*      $values[$n] = $v;*/
      /*  } else {*/
      /*    $placeholders[] = self::PARAM_PREFIX . $i;*/
      /*    $values[self::PARAM_PREFIX . $i] = $column->typecast($value);*/
      /*    $i++;*/
      /*  }*/
      /*}*/
    }
    print_r($fields);
    if ($fields === array()) {
      $pks = is_array($table->primaryKey) ? $table->primaryKey : array($table->primaryKey);
      foreach ($pks as $pk) {
        //$fields[] = $table->getColumn($pk)->rawName;
        $fields[] = $pk;
        $placeholders[] = $this->getIntegerPrimaryKeyDefaultValue();
      }
    }
    $sql = "INSERT INTO {$table->rawName} (" . implode(', ', $fields) . ') VALUES (' . implode(', ', $placeholders) . ')';
    Yii::trace(['createInsertCommand.sql' => $sql]);
    $command = $this->_connection->createCommand($sql);
    Yii::trace(['bindValue' => $values]);
    foreach ($values as $name => $value) {
      $command->bindValue($name, $value);
    }

    return $command;
  }

  public function createUpdateCommand($table, $data, $criteria)
  {
    //$this->ensureTable($table);
    $fields = array();
    $values = array();
    $bindByPosition = isset($criteria->params[0]);
    $i = 0;
    foreach ($data as $name => $value) {
      if (($column = $table->getColumn($name)) !== null) {
        //if (false && $value instanceof CDbExpression) {
        if (FALSE) {
          $fields[] = $column->rawName . '=' . $value->expression;
          foreach ($value->params as $n => $v)
            $values[$n] = $v;
        } elseif ($bindByPosition) {
          //$fields[] = $column->rawName . '=?';
          //$values[] = $column->typecast($value);
          $fields[] = $column . '=?';
          $values[] = $value;
        } else {
          //$fields[] = $column->rawName . '=' . self::PARAM_PREFIX . $i;
          //$values[self::PARAM_PREFIX . $i] = $column->typecast($value);
          $fields[] = $column . '=' . self::PARAM_PREFIX . $i;
          $values[self::PARAM_PREFIX . $i] = $value;
          $i++;
        }
      }
    }
    if ($fields === array()) {
      echo "err-CDbException";
      exit;
      /*throw new CDbException(Yii::t(*/
      /*  'yii',*/
      /*  'No columns are being updated for table "{table}".',*/
      /*  array('{table}' => $table->name)*/
      /*));*/
    }
    $sql = "UPDATE {$table->rawName} SET " . implode(', ', $fields);
    $sql = $this->applyJoin($sql, $criteria->join);
    $sql = $this->applyCondition($sql, $criteria->condition);
    $sql = $this->applyOrder($sql, $criteria->order);
    $sql = $this->applyLimit($sql, $criteria->limit, $criteria->offset);

    Yii::trace(['createUpdateCommand.sql' => $sql]);
    $command = $this->_connection->createCommand($sql);
    $this->bindValues($command, array_merge($values, $criteria->params));

    return $command;
  }

  public function applyCondition($sql, $condition)
  {
    if ($condition != '')
      return $sql . ' WHERE ' . $condition;
    else
      return $sql;
  }

  public function applyJoin($sql, $join)
  {
    if ($join != '')
      return $sql . ' ' . $join;
    else
      return $sql;
  }

  public function applyOrder($sql, $orderBy)
  {
    if ($orderBy != '')
      return $sql . ' ORDER BY ' . $orderBy;
    else
      return $sql;
  }

  public function applyLimit($sql, $limit, $offset)
  {
    if ($limit >= 0)
      $sql .= ' LIMIT ' . (int)$limit;
    if ($offset > 0)
      $sql .= ' OFFSET ' . (int)$offset;
    return $sql;
  }

  public function applyGroup($sql, $group)
  {
    if ($group != '')
      return $sql . ' GROUP BY ' . $group;
    else
      return $sql;
  }

  public function applyHaving($sql, $having)
  {
    if ($having != '')
      return $sql . ' HAVING ' . $having;
    else
      return $sql;
  }

  public function bindValues($command, $values)
  {
    if (($n = count($values)) === 0)
      return;

    //var_dump("values", $values);

    if (isset($values[0])) // question mark placeholders
    {
      //var_dump("isset values[0]");
      for ($i = 0; $i < $n; ++$i) {
        $command->bindValue($i + 1, $values[$i]);
      }
    } else // named placeholders
    {
      //var_dump("named placeholders");
      foreach ($values as $name => $value) {
        if ($name[0] !== ':') {
          $name = ':' . $name;
        }
        $command->bindValue($name, $value);
      }
    }
  }

  protected function getIntegerPrimaryKeyDefaultValue()
  {
    return 'NULL';
  }
}
