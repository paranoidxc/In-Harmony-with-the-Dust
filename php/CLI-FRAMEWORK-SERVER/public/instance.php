<?php

define('ROOT_PATH',  dirname(__FILE__) . '/'); // 站点应用目录

require_once ROOT_PATH . "../core/TheOldHunter.php";

$cf                         = [
    'port'                 => 9000,
    'master_process_title' => 'HTTP SERVER',
    'worker_process_title' => "WORKER",
];
$cf2 = ['port' => 9000, 'TheOldHunter', "YaNan"];
$hunter = new TheOldHunter();
$hunter->hunter();
