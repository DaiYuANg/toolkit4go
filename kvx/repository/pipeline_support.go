package repository

import (
	"strconv"
	"time"

	"github.com/DaiYuANg/archgo/kvx"
)

type pipelineProvider interface {
	Pipeline() kvx.Pipeline
}

func enqueueExpire(pipe kvx.Pipeline, key string, expiration time.Duration) {
	if expiration <= 0 {
		return
	}

	pipe.Enqueue("PEXPIRE", []byte(key), []byte(strconv.FormatInt(expiration.Milliseconds(), 10)))
}
