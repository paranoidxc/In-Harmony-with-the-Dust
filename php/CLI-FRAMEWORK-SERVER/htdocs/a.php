<?php

echo "FROM " . __FILE__ . date("Y-m-d H:i:s") . "<br/>";
echo "GET PARAMS <BR>";
print_r($_GET);
echo " <BR>";
print_r ($decode_ret['get']);
echo " <BR>";
echo " <BR>";
echo " <BR>";
echo " <BR>";
/*
foreach (decode_ret['xget'] as $k => $v) {
    echo " {$k} = {$v}<BR>";
}
*/

echo "POST PARAMS <BR>";
foreach ($decode_ret['post'] as $k => $v) {
    echo " {$k} = {$v}<BR>";
}

echo "JSON PARAMS <BR>";
