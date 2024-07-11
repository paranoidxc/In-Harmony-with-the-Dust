package cronjob

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"zero-zone/applet/api/internal/svc"
)

func AddTaskByFuncWithCompensation(
	svc *svc.ServiceContext,
	cacheKey string,
	cronName string,
	taskName string,
	spec string,
	second int,
	fun func(),
	option ...cron.Option) (err error) {

	taskTokenVal, err := svc.Redis.Get(cacheKey)
	if err != nil || taskTokenVal == "" {
		logx.Infow(fmt.Sprintf("任务[%s:%s]没有启动过不补偿", cronName, taskName))
		err = svc.Redis.Set(cacheKey, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			return
		}
		// 新建定时任务
		cronScheduler := cron.New(cron.WithSeconds())
		// ============= add job end =============
		_, err := cronScheduler.AddFunc(spec, fun)
		if err != nil {
			logx.Infow(fmt.Sprintf("任务[%s:%s]添加失败", cronName, taskName), logx.Field("err", err))
		} else {
			logx.Infow(fmt.Sprintf("任务[%s:%s]添加成功", cronName, taskName))
		}
		// ============= add job end =============
		cronScheduler.Start()
	} else {
		logx.Infow(fmt.Sprintf("任务[%s:%s]上次执行时间[%s]\n", cronName, taskName, taskTokenVal))
		t, err := time.Parse("2006-01-02 15:04:05", taskTokenVal)
		if err == nil {
			nowTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
			duration := nowTime.Sub(t)
			left := second - int(duration.Seconds())
			//fmt.Printf("================Every Seconds[%d], Left[%d]===================\n", second, left)
			if left <= 0 {
				left = 5
			}
			logx.Infow(fmt.Sprintf("任务[%s:%s]进行补偿[%d]s \n", cronName, taskName, left))
			go func() {
				c := cron.New(option...)
				tmpId, tmpErr := c.AddFunc(fmt.Sprintf("@every %ds", left), fun)
				if tmpErr == nil {
					go func() {
						ticker := time.NewTicker(time.Duration(left) * time.Second)
						<-ticker.C
						for {
							entry := c.Entry(tmpId)
							//fmt.Printf("%+v\n", entry)
							if entry.Prev.IsZero() {
								//fmt.Printf("sleep %s\n", cronName)
								time.Sleep(time.Second)
							} else {
								logx.Infow(fmt.Sprintf("移除补偿任务[%s:%s],taskId[%d]\n", cronName, taskName, tmpId))
								c.Remove(tmpId)
								// 补偿执行完毕后，新建定时任务
								go func() {
									cronScheduler := cron.New(cron.WithSeconds())
									// ============= add job end =============
									_, err := cronScheduler.AddFunc(spec, fun)
									if err != nil {
										logx.Infow(fmt.Sprintf("任务[%s:%s]添加失败", cronName, taskName), logx.Field("err", err))
									} else {
										logx.Infow(fmt.Sprintf("任务[%s:%s]添加成功", cronName, taskName))
									}
									// ============= add job end =============
									cronScheduler.Start()
								}()
								break
							}
						}
					}()
				}
				c.Start()
			}()
		}
	}

	return nil
}
