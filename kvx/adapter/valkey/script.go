package valkey

import (
	"context"
	"github.com/valkey-io/valkey-go"
)

// ============== Script Interface ==============

// Load loads a script into the script cache.
func (a *Adapter) Load(ctx context.Context, script string) (string, error) {
	resp := a.client.Do(ctx, a.client.B().ScriptLoad().Script(script).Build())
	if resp.Error() != nil {
		return "", resp.Error()
	}
	return resp.ToString()
}

// Eval executes a script.
func (a *Adapter) Eval(ctx context.Context, script string, keys []string, args [][]byte) ([]byte, error) {
	// Build eval command
	argStrs := make([]string, len(args))
	for i, arg := range args {
		argStrs[i] = valkey.BinaryString(arg)
	}

	cmd := a.client.B().Eval().Script(script).Numkeys(int64(len(keys))).Key(keys...).Arg(argStrs...)

	resp := a.client.Do(ctx, cmd.Build())
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	return resp.AsBytes()
}

// EvalSHA executes a cached script by SHA.
func (a *Adapter) EvalSHA(ctx context.Context, sha string, keys []string, args [][]byte) ([]byte, error) {
	argStrs := make([]string, len(args))
	for i, arg := range args {
		argStrs[i] = valkey.BinaryString(arg)
	}

	cmd := a.client.B().Evalsha().Sha1(sha).Numkeys(int64(len(keys))).Key(keys...).Arg(argStrs...)

	resp := a.client.Do(ctx, cmd.Build())
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	return resp.AsBytes()
}
