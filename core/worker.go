package core

import (
	"github.com/gawwo/fake115-go/config"
	"github.com/gawwo/fake115-go/dir"
	"go.uber.org/zap"
)

type Task struct {
	Dir  *dir.Dir
	File *NetFile
}

var WorkerChannel = make(chan Task, config.WorkerNum*config.WorkerNumRate)

func Worker() {
	// WorkerChannel关闭前一直工作，直到生产者枯竭
	for task := range WorkerChannel {
		// 有recover，保证这里不会panic，能让任务持续进行

		result := task.File.Export()
		if result == "" {
			config.Logger.Warn("export failed", zap.String("name", task.File.Name))
			continue
		}
		// 扫尾工作，添加记录到dir对象，累加文件总大小
		lock.Lock()
		task.Dir.Files = append(task.Dir.Files, result)
		config.TotalSize += task.File.Size
		lock.Unlock()
	}
	config.WaitGroup.Done()
}
