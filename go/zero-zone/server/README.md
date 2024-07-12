# Zeor-Zone-Admin

# 项目启动

## 添加 v10验证 

```
$ go get github.com/go-playground/validator/v10
$ go get github.com/go-playground/universal-translator
```

## 定时库

```
$ go get github.com/robfig/cron/v3@v3.0.0
```

doc: https://pkg.go.dev/github.com/robfig/cron#section-readme

## 文档生成 Swagger API

 [GitHub - sliveryou/goctl-swagger: 自定义 goctl-swagger](https://github.com/sliveryou/goctl-swagger)

```
$ go install github.com/sliveryou/goctl-swagger@latest
```

## 后台启动 8008 端口

```
cd ~/applet/api
go mod tidy
go run core.go

测试后台是否正常
http://ip:port/admin/user/login/captcha
```

## 前台启动 5100 端口 

```
cd ~/web
npm install
npm run dev //会自动打开浏览器显示登录页面
```

### 引入md5

```
npm install --save js-md5
import md5 from 'js-md5'
md5('xxxxxxxx')
```



# 注意点

## 数据库表统一个格式

数据库表 尽量统一维持这4字段

id 自增主键，创建时间 更新时间 删除时间

```
ALTER TABLE tb_name
ADD  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';

ALTER TABLE tb_name
ADD  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';

ALTER TABLE tb_name
ADD  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间';
```

## CRUD 自动化

前端自动生成文件

    - web/src/api/feat/xxx.js //ajax 直接请求后台接口的文件
    - web/src/view/feat/xxx.vue // vue文件 回调用 xxx.js

后端代码自动生成

    - 生成 goctl 需要的api文件
    - 生成 handler/feat 目录下的控制器文件
    - 生成 logic/feat目录下的逻辑文件
    - 生成 model 目录下的模型文件

# 如何初始化一个功能的CURD

1. 在 model/tmpAutoFeat/ 建立一个文件 

   1. 里面的结构体以 Tmp 打头 
   2. 模仿 TmpDemoCurd 即可 里面的 gorm tag 只是用来分析使用，暂时没有用到 gorm 来操作数据库
   3. 第一个字段需要是主键
2. 在 model/autoCurd.go 下建立map关系
3. 启动项目
4. 登录后台使用页面操作
5. 在 svc/servicecontext.go 下配置对应的模型和变量
6. 重启系统
7. 之后如果需要重新初始化该功能的话 需要手动删除生成的文件 
8. 如果后期是功能的添加 那么只需要 修改 desc/xxx.api 然后使用 goctl 手动生成 （这样只会添加新的handle和logic 不会覆盖旧的）
9. 生成的 vue 文件基本是 text 类型的 需要自己按需调整
10. 生成的 vue 文件搜索页面 对应请求到的 logic 需要到 logic 中 把注释取消 按需调整

## goctl 手动生成路由/控制器/逻辑文件    控制器/逻辑文件 可能需要删除额外信息多导入的包

```
cd api
goctl api go -api core.api --style goZero -dir . --home=../tpl/
```

## goctl 手动生成 模型文件

```
cd applet
goctl model mysql datasource -url="root:fucking@tcp(127.0.0.1:13307)/zero_zone" -table="demo_curd"  -dir=./model -cache true --style=goZero --home ./tpl -i "created_at,updated_at,deleted_at" 
```

# 启动普罗米修斯 查看指标

普罗米修斯配置文件 里面的配置

```
 ~/etc/prometheus.yml  
 
global:
  scrape_interval: 15s # 默认情况下每15秒抓取一次数据
  evaluation_interval: 15s # 默认规则每15秒评估一次

scrape_configs:
  - job_name: 'verification_system'
    static_configs:
      - targets: ['host.docker.internal:6060'] # 从 docker 中请求到本机host
```

docker 方式安装 prometheus 修改 -v 的本地配置文件的地址

```
docker run -d \
-p 9090:9090 \
-v xxxxx/etc/prometheus.yml:/etc/prometheus/prometheus.yml \
prom/prometheus
```

访问 http://localhost:9090/

输入框输入 

```
{path="/admin/feat/demoCurdPage/page"}
{method="GET"}
{code="200"}
```



# 启动 链路追踪

```
etc/core-api.yaml

# 链路追踪配置
Telemetry:
  Name: verification.system.admin.api
  Endpoint: http://127.0.0.1:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger
```



```text
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 16686:16686 -p 14268:14268  -p 14269:14269   -p 9411:9411 -p 6831:6831/udp jaegertracing/all-in-one:latest
```

访问 http://127.0.0.1:16686/search

**先启动 jaeger， 然后再启动项目**

访问项目 查看日志 在日志的输入中会带有  trace=ecad00f54c41433e432bd7c9931b77cd  span=10efcec8ff2a2814

这个 trace 串就是用来查询的