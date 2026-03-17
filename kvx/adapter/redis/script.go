package redis

import (
	"context"
)

// ============== Script Interface ==============

// Load loads a script into the script cache.
func (a *Adapter) Load(ctx context.Context, script string) (string, error) {
	return a.client.ScriptLoad(ctx, script).Result()
}

// Eval executes a script.
func (a *Adapter) Eval(ctx context.Context, script string, keys []string, args [][]byte) ([]byte, error) {
	ifaceArgs := make([]interface{}, len(args))
	for i, v := range args {
		ifaceArgs[i] = v
	}

	val, err := a.client.Eval(ctx, script, keys, ifaceArgs...).Result()
	if err != nil {
		return nil, err
	}

	return valueToBytes(val)
}

// EvalSHA executes a cached script by SHA.
func (a *Adapter) EvalSHA(ctx context.Context, sha string, keys []string, args [][]byte) ([]byte, error) {
	ifaceArgs := make([]interface{}, len(args))
	for i, v := range args {
		ifaceArgs[i] = v
	}

	val, err := a.client.EvalSha(ctx, sha, keys, ifaceArgs...).Result()
	if err != nil {
		return nil, err
	}

	return valueToBytes(val)
}
