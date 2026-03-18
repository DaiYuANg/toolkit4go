package valkey

import (
	"context"
	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/valkey-io/valkey-go"
)

// ============== Pipeline Interface ==============

// Pipeline creates a new pipeline.
func (a *Adapter) Pipeline() kvx.Pipeline {
	return &valkeyPipeline{
		client: a.client,
	}
}

type valkeyPipeline struct {
	client valkey.Client
	cmds   []valkey.Completed
}

// Enqueue adds a command to the pipeline.
func (p *valkeyPipeline) Enqueue(command string, args ...[]byte) {
	argStrs := make([]string, len(args))
	for i, v := range args {
		argStrs[i] = valkey.BinaryString(v)
	}
	cmd := p.client.B().Arbitrary(command).Args(argStrs...).Build()
	p.cmds = append(p.cmds, cmd)
}

// Exec executes all queued commands.
func (p *valkeyPipeline) Exec(ctx context.Context) ([][]byte, error) {
	if len(p.cmds) == 0 {
		return nil, nil
	}

	// Use DoMulti for pipeline execution
	resps := p.client.DoMulti(ctx, p.cmds...)

	results := make([][]byte, len(resps))
	for i, resp := range resps {
		if resp.Error() != nil && !valkey.IsValkeyNil(resp.Error()) {
			results[i] = nil
			continue
		}
		b, _ := resp.AsBytes()
		results[i] = b
	}
	return results, nil
}

// Close closes the pipeline.
func (p *valkeyPipeline) Close() error {
	// No explicit close needed
	return nil
}
