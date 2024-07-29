<?php
ob_start();
echo "FROM BBBBB";
$a = ['aaaaaaaa'];
print_r($_GET);
var_dump($a);
$bd = ob_get_contents();
ob_flush();
echo $bd."<BR>";
//echo $_GET['ss'];
echo "FROM " . __FILE__ . date("Y-m-d H:i:s") . "<br/>";
