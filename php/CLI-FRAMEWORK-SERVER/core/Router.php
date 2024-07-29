<?php
class Router
{
    public static function run($content)
    {
        $pathinfo = $content['pathinfo'] ?? '';
        $out = '';
        try {
            ob_start();
            self::callCntl(self::parse($pathinfo), $content);
            $out = ob_get_contents();
            ob_end_clean();
        } catch (Exception $e) {
            $out = $e->getMessage();
        }

        return $out;
    }

    // $this->add('', ['controller'=>'Home', 'action'=> 'index'])
    // $this->add('posts', ['controller'=>'Posts', 'action'=> 'index'])
    // $this->add('{controller}/{action}', [])
    // $this->add('{controller}/{id:\d+}/{action}', [])
    private static function add($rou, $params = [])
    {
        //chapter 20,21
        // convert the route to a regular expression: escape forward slashes
        $rou = preg_replace('/\//', '\\/', $rou);

        // convert variables
        $rou = preg_replace('/\{([a-z]+)\}/', '(?P<\1>[a-z-]+)', $rou);

        // convert variables with custom regular expressions
        $rou = preg_replace('/\{([a-z]+):([^\]+)\}/', '(?P<\1>\2)', $rou);

        // add start and end delimiters and case insensitive flag
        $rou = '/^' . $rou . '$/i';

        //match url format /controller/action
        //$reg_exp = '/^(?P<controller>[a-z-]+)\/(?P<action>[a-z-]+)$/';
    }

    private static function parse($pathinfo)
    {

        Log::info("info", $pathinfo);
        $info = explode("/", str_replace('.php', '', $pathinfo));

        $module = '';
        $controller = NULL;
        $action = NULL;
        $is_controller_done = FALSE;

        $map_kv = [];
        $modules = [];

        Log::info("info", $info);
        foreach ($info as $comp) {
            if (strlen($comp)) {
                if (in_array($comp, $modules)) {
                    $module = $comp;
                    continue;
                }
                if ($is_controller_done) {
                    if (!is_null($action)) {
                        $map_kv[] = $comp;
                    } else {
                        $action = $comp;
                    }
                } else {
                    $controller = $comp;
                }

                if (!in_array($comp, $modules)) {
                    $is_controller_done = TRUE;
                }
            }
        }

        Log::info("Controller", $controller);

        $controller = ucfirst(strtolower($controller ?? 'Default'));
        $action = $action ?? 'index';

        $cntl_class_name  = $controller . "Cntl";
        $cntl_method_name = "act" . ucwords(strtolower($action));

        $params = [];
        for ($i = 0; $i < count($map_kv); $i += 2) {
            $params[$map_kv[$i]] = ($map_kv[$i + 1]) ?? "";
        }
        //Log::info('params', $params);

        $router                  = [
            'module'           => $module,
            'cntl'             => $controller,
            'act'              => $action,
            'cntl_class_name'  => $cntl_class_name,
            'cntl_method_name' => $cntl_method_name,
            'params'           => $params,
        ];

        Log::info("router info", $router);
        return $router;
    }

    private static function callCntl($router, $content)
    {
        $cntl_class_name = $router['cntl_class_name'];
        $cntl_method_name = $router['cntl_method_name'];

        if (!class_exists($cntl_class_name)) {
            self::loadFile(
                self::absoluteControllFilePath($cntl_class_name, $router),
                $cntl_class_name
            );
        }

        $c = new $cntl_class_name($router);
        $c->get = $content['get'];
        $c->post = $content['post'];
        $c->router = $router;
        Log::info("RRRRRRRRr", $router);

        if (method_exists($c, $cntl_method_name)) {
            if (method_exists($c, "befAct")) {
                call_user_func([$c, "befAct"]);
            }
            call_user_func([$c, $cntl_method_name]);
            if (method_exists($c, "afAct")) {
                call_user_func([$c, "afAct"]);
            }
        } else {
            Log::err(
                'Cntl Method Not Exist',
                [
                    'class'  => $cntl_class_name,
                    'method' => $cntl_method_name,
                ]
            );
            throw new Exception("Controller:{$cntl_class_name} Method:$cntl_method_name Not Exist");
        }
    }

    protected static function loadFile($file, $class_name = '')
    {
        if (strlen($class_name)) {
            if (class_exists($class_name)) {
                return;
            }
        }
        if (file_exists($file)) {
            require_once($file);
        } else {
            Log::err("File Not Exist", ['file' => $file, 'class_name' => $class_name]);
            throw new Exception("File Not Exist {$file}");
        }
    }

    public static function absoluteViewFilePath($file_name, $router = [])
    {
        $dir = [
            "view",
            lcfirst($router['cntl']),
            $file_name . ".php"
        ];
        return  self::absoluteFilePath($router, $dir, FALSE);
    }

    private static function absoluteControllFilePath($file_name, $router = [])
    {
        $dir = [
            "cntl",
            $file_name . ".php"
        ];
        return self::absoluteFilePath($router, $dir, FALSE);
    }

    private static function absoluteFilePath($router = [], $merge_dir = [])
    {
        $module = $router['module'];
        $dir = [APP_PATH];
        if (isset($module[0])) {
            $dir = array_merge(
                $dir,
                ["module", $module]
            );
        }

        $file = join(
            DIRECTORY_SEPARATOR,
            array_merge(
                $dir,
                $merge_dir
            )
        );

        return $file;
    }
}
