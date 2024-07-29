# TheOldHunter 

##居在 福州54北 找新工作 后端开发(PHP 或 GO)

## CLI-HTTP-SERVER-WITH-FRAMEWORK

** 实现一个基于 cli 常驻内存模式的 最小化 http 带个人 orm 和 lib 的服务框架 **

** to be continue... **

Usage:

- php >= 7.1
- 编译php configure --enable-pcntl --enable-sockets --enable-shmop
- php -m | grep event
  1. pecl install event
  2. configure event from source code
     1. install libevent on OS
     2. configure install event and add event.so to php.ini file
- 启动服务 php public/instance.php start|stop|reload
- 浏览器打开
  - 127.0.0.1:9000/
  - 127.0.0.1:9000/default/index?a=b&c=d

Todo List:

- [X] 守护进程模式运行
- [X] 监控子进程
- [X] 多进程模式
- [X] IO复用
  - [X] SELECT 模式
  - [X] EPOLL 模式
- [X] 支持 HTTP 协议
  - [X] POST
  - [X] GET
- [ ] 支持 HTTPS 协议
- [X] 代码重构成面向对象模式
- [ ] WEB 框架
  - [X] ROUTE
  - [ ] MVC
    - [x] Controller
    - [x] View
    - [ ] Model 
  - [ ] LIB
- [X] 自动热更新
  - [X] fswatch
  - [X] inotify
- [X] TINY LOG
