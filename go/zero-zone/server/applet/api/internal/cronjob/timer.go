package cronjob

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest"
	"zero-zone/applet/api/internal/svc"
)

func InitTimer(server *rest.Server, serverCtx *svc.ServiceContext) {
	// Start the scheduler in a separate goroutine
	fmt.Println("cronJob initTimer")

	/*
		go func(serverCtx *svc.ServiceContext) {
			cronName := "test"
			cacheKey := fmt.Sprintf("cache:verificationSystem:cronjob:task:%s", cronName)
			err := AddTaskByFuncWithCompensation(serverCtx,
				cacheKey,
				cronName,
				"测试",
				"@every 1m",
				60,
				func() {
					serverCtx.Redis.Set(cacheKey, time.Now().Format("2006-01-02 15:04:05"))
					fmt.Println("Hello, World!")
				},
			)
			if err != nil {
				fmt.Println("定时任务添加失败", err)
			}
		}(serverCtx)
	*/

	/*
		go func(serverCtx *svc.ServiceContext) {
			cronName := "thirdClientSync"
			cacheKey := fmt.Sprintf("cache:verificationSystem:cronjob:task:%s", cronName)
			err := AddTaskByFuncWithCompensation(serverCtx,
				cacheKey,
				cronName,
				"授权列表同步功能",
				"@hourly",
				3600,
				func() {
					serverCtx.Redis.Set(cacheKey, time.Now().Format("2006-01-02 15:04:05"))
					taskClientSync(serverCtx)
				},
			)
			if err != nil {
				fmt.Println("定时任务 授权列表同步功能开启失败", err)
			}
		}(serverCtx)
	*/
	/*
		go func(server *rest.Server, svcCtx *svc.ServiceContext) {
			cronScheduler := cron.New(cron.WithSeconds())
			// ============= add job end =============
			_, err := cronScheduler.AddFunc("@hourly", func() {
				taskClientSync(svcCtx)
			})
			if err != nil {
				fmt.Println("授权列表同步功能开启失败", err)
			}
			// ============= add job end =============

			// Start the scheduler
			cronScheduler.Start()
		}(server, serverCtx)
	*/
}
