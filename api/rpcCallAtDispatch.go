package api

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
)

type TaskAtFutureMs struct {
	ServiceName  string
	TimeAtUnixNs int64
}

var TasksAtFutureList = []*TaskAtFutureMs{}
var mut sync.Mutex = sync.Mutex{}

// put parameter to redis ,make it persistent
func rpcCallAtTaskAddOne(serviceName string, timeAtStr string, bytesValue string) {
	var (
		rds *redis.Client = config.RdsDefaultClient()
		err error
	)
	task := &TaskAtFutureMs{ServiceName: serviceName}
	if task.TimeAtUnixNs, err = strconv.ParseInt(timeAtStr, 10, 64); err != nil {
		log.Info().Err(err).Send()
		return
	}
	if cmd := rds.HSet(context.Background(), serviceName+":delay", timeAtStr, bytesValue); cmd.Err() != nil {
		log.Info().Err(cmd.Err()).Send()
		return
	}
	index := sort.Search(len(TasksAtFutureList), func(i int) bool { return TasksAtFutureList[i].TimeAtUnixNs < task.TimeAtUnixNs })

	// Insert the new task into the TasksAtFuture at the found index.
	mut.Lock()
	TasksAtFutureList = append(TasksAtFutureList[:index], append([]*TaskAtFutureMs{task}, TasksAtFutureList[index:]...)...)
	mut.Unlock()
}
func rpcCallAtRoutine() {
	var (
		bytes                 []byte
		TaskAtFutureNs, nowNs int64
		err                   error
		cmd                   []redis.Cmder
		rds                   *redis.Client = config.RdsDefaultClient()
	)
	for {
		if len(TasksAtFutureList) == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		mut.Lock()
		nowNs, TaskAtFutureNs = time.Now().UnixNano(), TasksAtFutureList[0].TimeAtUnixNs
		task := TasksAtFutureList[0]
		TasksAtFutureList = TasksAtFutureList[1:]
		mut.Unlock()

		if timeSpan := TaskAtFutureNs - nowNs; timeSpan > 0 {
			if timeSpan > 100*1000*1000 {
				timeSpan = 100 * 1000 * 1000
			}
			time.Sleep(time.Duration(timeSpan))
			continue
		}
		pipeline := rds.Pipeline()
		pipeline.HGet(context.Background(), task.ServiceName+":delay", strconv.FormatInt(TaskAtFutureNs, 10))
		pipeline.HDel(context.Background(), task.ServiceName+":delay", strconv.FormatInt(TaskAtFutureNs, 10))
		if cmd, err = pipeline.Exec(context.Background()); err != nil {
			log.Info().Err(err).Send()
			continue
		}
		if bytes, err = cmd[0].(*redis.StringCmd).Bytes(); err != nil || len(bytes) == 0 {
			continue
		}
		CallApiLocallyAndSendBackResult(task.ServiceName, strconv.FormatInt(TaskAtFutureNs, 10), bytes)
	}
}

func rpcCallAtTasksLoad() {
	var (
		services   = apiServiceNames()
		timeAtStrs []string
		cmd        []redis.Cmder
		err        error
		rds        *redis.Client = config.RdsDefaultClient()
	)
	log.Info().Msg("rpcCallAtTasksLoading started")
	pipeline := rds.Pipeline()
	for _, service := range services {
		pipeline.HKeys(context.Background(), service+":delay")
	}
	if cmd, err = pipeline.Exec(context.Background()); err != nil {
		log.Info().AnErr("err LoadDelayApiTask, ", err).Send()
		return
	}
	var _TasksAtFutureList = []*TaskAtFutureMs{}
	for i, service := range services {
		if err = cmd[i].(*redis.StringSliceCmd).Err(); err != nil {
			continue
		}
		timeAtStrs = cmd[i].(*redis.StringSliceCmd).Val()
		for _, timeAtStr := range timeAtStrs {
			if timeAt, err := strconv.ParseInt(timeAtStr, 10, 64); err == nil {
				_TasksAtFutureList = append(_TasksAtFutureList, &TaskAtFutureMs{ServiceName: service, TimeAtUnixNs: timeAt})
			}
		}
	}
	sort.Slice(_TasksAtFutureList, func(i, j int) bool {
		return _TasksAtFutureList[i].TimeAtUnixNs < _TasksAtFutureList[j].TimeAtUnixNs
	})
	mut.Lock()
	TasksAtFutureList = _TasksAtFutureList
	mut.Unlock()
	log.Info().Msg("rpcCallAtTasksLoading completed")
}
func init() {
	go func() {
		rpcCallAtTasksLoad()

		rpcCallAtRoutine()
	}()

}
