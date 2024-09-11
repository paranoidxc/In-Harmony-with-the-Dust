<?php

namespace Paranoid\Framework\Utils;

class TreeChildrenIds
{
    private $idName = 'id';
    private $parentIdName = 'parent_id';
    private $findId;
    private $rs = [];
    private $mapDtChildrenIds = [];

    function __construct($rs, $findId, $idName = '', $parentIdName = '')
    {
        if (is_array($rs)) {
            $this->rs = $rs;
        }

        if (strlen(trim($idName))) {
            $this->idName = trim($idName);
        }

        if (strlen(trim($parentIdName))) {
            $this->parentIdName = trim($parentIdName);
        }

        $this->findId = $findId;
    }

    public static function run($rs, $findId, $idName = '', $parentIdName = '')
    {
        if (!strlen($findId)) {
            return [];
        }

        $o = new self($rs, $findId, $idName, $parentIdName);
        return $o->process();
    }

    private function process()
    {
        $mapDtChildrenIds = [];

        foreach ($this->rs as $row) {
            $id = $row[$this->idName];
            $parentId = $row[$this->parentIdName];
            if (!isset($mapDtChildrenIds[$parentId])) {
                $mapDtChildrenIds[$parentId] = [];
            }
            $mapDtChildrenIds[$parentId][] = $id;
        }

        $this->mapDtChildrenIds = $mapDtChildrenIds;

        $allIds = $this->getAllChildrenIds($this->findId);

        return $allIds;
    }

    private function getAllChildrenIds($id)
    {
        $r = [];

        if (isset($this->mapDtChildrenIds[$id])) {
            $dt_subs = $this->mapDtChildrenIds[$id];
            foreach ($dt_subs as $_id) {
                $r = array_merge($this->getAllChildrenIds($_id), $r);
            }
        }
        $r[] = $id;

        return $r;
    }
}