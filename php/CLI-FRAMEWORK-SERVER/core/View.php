<?php
class View
{
    public static function render($view_file, $args = [], $return = FALSE)
    {
        //Log::info("PARAMS", $args);
        extract($args);
        if ($return) {
            ob_start();
        }
        include $view_file;
        if ($return) {
            $out = ob_get_contents();
            ob_end_clean();

            return $out;
        }
    }
}
