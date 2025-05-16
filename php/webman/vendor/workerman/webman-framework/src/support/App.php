<?php

namespace support;

use Dotenv\Dotenv;
use RuntimeException;
use Throwable;
use Webman\Config;
use Webman\Util;
use Workerman\Connection\TcpConnection;
use Workerman\Worker;
use function base_path;
use function call_user_func;
use function is_dir;
use function opcache_get_status;
use function opcache_invalidate;
use const DIRECTORY_SEPARATOR;

class App
{
  /**
   * Run.
   * @return void
   * @throws Throwable
   */
  public static function run()
  {
    ini_set('display_errors', 'on');
    error_reporting(E_ALL);

    if (class_exists(Dotenv::class) && file_exists(run_path('.env'))) {
      if (method_exists(Dotenv::class, 'createUnsafeImmutable')) {
        Dotenv::createUnsafeImmutable(run_path())->load();
      } else {
        Dotenv::createMutable(run_path())->load();
      }
    }

    if (!$appConfigFile = config_path('app.php')) {
      throw new RuntimeException('Config file not found: app.php');
    }
    $appConfig = require $appConfigFile;
    //var_dump("appConfig", $appConfig);
    if ($timezone = $appConfig['default_timezone'] ?? '') {
      date_default_timezone_set($timezone);
    }

    static::loadAllConfig(['route', 'container']);

    // 是 PHP 提供的一个内置类，用于创建和操作 PHAR 文件 （PHP Archive）。
    // PHAR 文件是一种将多个 PHP 文件、资源文件（如图片、CSS、JS 等）以及元数据打包成一个单独文件的格式。
    // 它类似于 Java 的 JAR 文件或 ZIP 压缩文件，但专门为 PHP 设计
    if (!is_phar() && DIRECTORY_SEPARATOR === '\\' && empty(config('server.listen'))) {
      echo "Please run 'php windows.php' on windows system." . PHP_EOL;
      exit;
    }

    $errorReporting = config('app.error_reporting');
    if (isset($errorReporting)) {
      error_reporting($errorReporting);
    }

    $runtimeLogsPath = runtime_path() . DIRECTORY_SEPARATOR . 'logs';
    if (!file_exists($runtimeLogsPath) || !is_dir($runtimeLogsPath)) {
      if (!mkdir($runtimeLogsPath, 0777, true)) {
        throw new RuntimeException("Failed to create runtime logs directory. Please check the permission.");
      }
    }

    $runtimeViewsPath = runtime_path() . DIRECTORY_SEPARATOR . 'views';
    if (!file_exists($runtimeViewsPath) || !is_dir($runtimeViewsPath)) {
      if (!mkdir($runtimeViewsPath, 0777, true)) {
        throw new RuntimeException("Failed to create runtime views directory. Please check the permission.");
      }
    }

    // 这段代码是在 Webman 框架中定义的一个主重载（master reload）回调函数。
    // 它的主要作用是在 Workerman 的主进程重新加载时清除 PHP 的 OPcache 缓存
    // 在 Worker.php 的 reload 方法中发现了对 onMasterReload 回调的调用： protected static function reload(): void
    Worker::$onMasterReload = function () {
      if (function_exists('opcache_get_status')) {
        if ($status = opcache_get_status()) {
          if (isset($status['scripts']) && $scripts = $status['scripts']) {
            foreach (array_keys($scripts) as $file) {
              opcache_invalidate($file, true);
            }
          }
        }
      }
    };

    $config = config('server');
    Worker::$pidFile = $config['pid_file'];
    Worker::$stdoutFile = $config['stdout_file'];
    Worker::$logFile = $config['log_file'];
    Worker::$eventLoopClass = $config['event_loop'] ?? '';
    TcpConnection::$defaultMaxPackageSize = $config['max_package_size'] ?? 10 * 1024 * 1024;
    if (property_exists(Worker::class, 'statusFile')) {
      Worker::$statusFile = $config['status_file'] ?? '';
    }
    if (property_exists(Worker::class, 'stopTimeout')) {
      Worker::$stopTimeout = $config['stop_timeout'] ?? 2;
    }

    put_log([
      'tip' => 'config info',
      'listen' => $config['listen'] ?? false,
      'event_loop' => $config['event_loop'],
    ]);
    if ($config['listen'] ?? false) {
      put_log([
        'tip' => 'config set listen work',
      ]);
      $worker = new Worker($config['listen'], $config['context']);
      $propertyMap = [
        'name',
        'count',
        'user',
        'group',
        'reusePort',
        'transport',
        'protocol'
      ];
      foreach ($propertyMap as $property) {
        if (isset($config[$property])) {
          $worker->$property = $config[$property];
        }
      }

      // 这段代码的主要作用是将 Webman 应用与 Workerman 的事件驱动模型连接起来，
      // 使 HTTP 请求能够被 Webman 框架正确处理。
      // 当 HTTP 请求到达时，会触发 onMessage 回调，进而执行 Webman 应用的路由分发和请求处理流程。
      // 实际使用时，这段代码会在应用启动时执行，为每个工作进程设置好处理 HTTP 请求的环境和回调函数，
      // 确保应用能够正确响应用户请求。
      $worker->onWorkerStart = function ($worker) {
        require_once base_path() . '/support/bootstrap.php';
        $app = new \Webman\App(config('app.request_class', Request::class), Log::channel('default'), app_path(), public_path());
        //var_dump(" ---- worker->onWorkerStart ---- ;");
        //var_dump(" ---- app ---- ;", $app, "\n");
        $worker->onMessage = [$app, 'onMessage'];
        call_user_func([$app, 'onWorkerStart'], $worker);
      };
    }

    $windowsWithoutServerListen = is_phar() && DIRECTORY_SEPARATOR === '\\' && empty($config['listen']);
    $process = config('process', []);

    put_log([
      'tip' => "[main] ----- info",
      "windowsWithoutServerListen" => $windowsWithoutServerListen,
      "count process" => count($process),
    ]);
    //var_dump("process", $process);
    if ($windowsWithoutServerListen && $process) {
      $processName = isset($process['webman']) ? 'webman' : key($process);
      worker_start($processName, $process[$processName]);
    } else if (DIRECTORY_SEPARATOR === '/') {
      put_log([
        'tip' => "[main] ----- is linux",
      ]);
      foreach (config('process', []) as $processName => $config) {
        put_log([
          'tip' => "[main] ----- start process",
          'processName' => $processName,
        ]);
        worker_start($processName, $config);
      }

      foreach (config('plugin', []) as $firm => $projects) {
        foreach ($projects as $name => $project) {
          if (!is_array($project)) {
            continue;
          }
          foreach ($project['process'] ?? [] as $processName => $config) {
            worker_start("plugin.$firm.$name.$processName", $config);
          }
        }
        foreach ($projects['process'] ?? [] as $processName => $config) {
          worker_start("plugin.$firm.$processName", $config);
        }
      }
    }

    Worker::runAll();
  }

  /**
   * LoadAllConfig.
   * @param array $excludes
   * @return void
   */
  public static function loadAllConfig(array $excludes = [])
  {
    Config::load(config_path(), $excludes);
    $directory = base_path() . '/plugin';
    foreach (Util::scanDir($directory, false) as $name) {
      $dir = "$directory/$name/config";
      if (is_dir($dir)) {
        Config::load($dir, $excludes, "plugin.$name");
      }
    }
  }
}
