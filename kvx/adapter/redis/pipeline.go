package redis

import (
	"context"
	"errors"
	"github.com/DaiYuANg/archgo/kvx"
	"github.com/redis/go-redis/v9"
)

// ============== Pipeline Interface ==============

// Pipeline creates a new pipeline.
func (a *Adapter) Pipeline() kvx.Pipeline {
	return &redisPipeline{
		pipe: a.client.Pipeline(),
	}
}

type redisPipeline struct {
	pipe redis.Pipeliner
}

// Enqueue adds a command to the pipeline.
func (p *redisPipeline) Enqueue(command string, args ...[]byte) {
	// Convert args to interface{}
	ifaceArgs := make([]interface{}, len(args)+1)
	ifaceArgs[0] = command
	for i, v := range args {
		ifaceArgs[i+1] = v
	}
	p.pipe.Do(context.Background(), ifaceArgs...)
}

// Exec executes all queued commands.
func (p *redisPipeline) Exec(ctx context.Context) ([][]byte, error) {
	cmders, err := p.pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	results := make([][]byte, len(cmders))
	for i, cmd := range cmders {
		val, err := cmd.(*redis.Cmd).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			results[i] = nil
			continue
		}
		results[i], _ = valueToBytes(val)
	}
	return results, nil
}

// Close closes the pipeline.
func (p *redisPipeline) Close() error {
	// Pipeline doesn't need explicit close in go-redis
	return nil
}
