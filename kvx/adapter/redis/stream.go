package redis

import (
	"context"
	"errors"
	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/redis/go-redis/v9"
	"time"
)

// ============== Stream Interface ==============

// XAdd adds an entry to a stream.
func (a *Adapter) XAdd(ctx context.Context, key string, id string, values map[string][]byte) (string, error) {
	// Convert map[string][]byte to map[string]interface{}
	ifaceValues := make(map[string]interface{}, len(values))
	for k, v := range values {
		ifaceValues[k] = v
	}

	args := &redis.XAddArgs{
		Stream: key,
		Values: ifaceValues,
	}
	if id != "*" {
		args.ID = id
	}

	return a.client.XAdd(ctx, args).Result()
}

// XRead reads entries from a stream.
func (a *Adapter) XRead(ctx context.Context, key string, start string, count int64) ([]kvx.StreamEntry, error) {
	streams := []string{key, start}

	result, err := a.client.XRead(ctx, &redis.XReadArgs{
		Streams: streams,
		Count:   count,
		Block:   0,
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	entries := make([]kvx.StreamEntry, len(result[0].Messages))
	for i, msg := range result[0].Messages {
		entries[i] = kvx.StreamEntry{
			ID:     msg.ID,
			Values: convertInterfaceMapToBytes(msg.Values),
		}
	}
	return entries, nil
}

// XReadMultiple reads entries from multiple streams.
func (a *Adapter) XReadMultiple(ctx context.Context, streams map[string]string, count int64, block time.Duration) (map[string][]kvx.StreamEntry, error) {
	streamKeys := make([]string, 0, len(streams)*2)
	for key, start := range streams {
		streamKeys = append(streamKeys, key, start)
	}

	result, err := a.client.XRead(ctx, &redis.XReadArgs{
		Streams: streamKeys,
		Count:   count,
		Block:   block,
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return make(map[string][]kvx.StreamEntry), nil
		}
		return nil, err
	}

	entries := make(map[string][]kvx.StreamEntry)
	for _, stream := range result {
		streamEntries := make([]kvx.StreamEntry, len(stream.Messages))
		for i, msg := range stream.Messages {
			streamEntries[i] = kvx.StreamEntry{
				ID:     msg.ID,
				Values: convertInterfaceMapToBytes(msg.Values),
			}
		}
		entries[stream.Stream] = streamEntries
	}
	return entries, nil
}

// XRange reads entries in a range.
func (a *Adapter) XRange(ctx context.Context, key string, start, stop string) ([]kvx.StreamEntry, error) {
	result, err := a.client.XRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]kvx.StreamEntry, len(result))
	for i, msg := range result {
		entries[i] = kvx.StreamEntry{
			ID:     msg.ID,
			Values: convertInterfaceMapToBytes(msg.Values),
		}
	}
	return entries, nil
}

// XRevRange reads entries in reverse order.
func (a *Adapter) XRevRange(ctx context.Context, key string, start, stop string) ([]kvx.StreamEntry, error) {
	result, err := a.client.XRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]kvx.StreamEntry, len(result))
	for i, msg := range result {
		entries[i] = kvx.StreamEntry{
			ID:     msg.ID,
			Values: convertInterfaceMapToBytes(msg.Values),
		}
	}
	return entries, nil
}

// XLen gets the number of entries in a stream.
func (a *Adapter) XLen(ctx context.Context, key string) (int64, error) {
	return a.client.XLen(ctx, key).Result()
}

// XTrim trims the stream to approximately maxLen entries.
func (a *Adapter) XTrim(ctx context.Context, key string, maxLen int64) error {
	return a.client.XTrimMaxLen(ctx, key, maxLen).Err()
}

// XDel deletes specific entries from a stream.
func (a *Adapter) XDel(ctx context.Context, key string, ids []string) error {
	return a.client.XDel(ctx, key, ids...).Err()
}

// XGroupCreate creates a consumer group.
func (a *Adapter) XGroupCreate(ctx context.Context, key string, group string, startID string) error {
	return a.client.XGroupCreate(ctx, key, group, startID).Err()
}

// XGroupDestroy destroys a consumer group.
func (a *Adapter) XGroupDestroy(ctx context.Context, key string, group string) error {
	return a.client.XGroupDestroy(ctx, key, group).Err()
}

// XGroupCreateConsumer creates a consumer in a group.
func (a *Adapter) XGroupCreateConsumer(ctx context.Context, key string, group string, consumer string) error {
	return a.client.XGroupCreateConsumer(ctx, key, group, consumer).Err()
}

// XGroupDelConsumer deletes a consumer from a group.
func (a *Adapter) XGroupDelConsumer(ctx context.Context, key string, group string, consumer string) error {
	return a.client.XGroupDelConsumer(ctx, key, group, consumer).Err()
}

// XReadGroup reads entries as part of a consumer group.
func (a *Adapter) XReadGroup(ctx context.Context, group string, consumer string, streams map[string]string, count int64, block time.Duration) (map[string][]kvx.StreamEntry, error) {
	streamKeys := make([]string, 0, len(streams)*2)
	for key, start := range streams {
		streamKeys = append(streamKeys, key, start)
	}

	result, err := a.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  streamKeys,
		Count:    count,
		Block:    block,
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return make(map[string][]kvx.StreamEntry), nil
		}
		return nil, err
	}

	entries := make(map[string][]kvx.StreamEntry)
	for _, stream := range result {
		streamEntries := make([]kvx.StreamEntry, len(stream.Messages))
		for i, msg := range stream.Messages {
			streamEntries[i] = kvx.StreamEntry{
				ID:     msg.ID,
				Values: convertInterfaceMapToBytes(msg.Values),
			}
		}
		entries[stream.Stream] = streamEntries
	}
	return entries, nil
}

// XAck acknowledges processing of stream entries.
func (a *Adapter) XAck(ctx context.Context, key string, group string, ids []string) error {
	return a.client.XAck(ctx, key, group, ids...).Err()
}

// XPending gets pending entries information.
func (a *Adapter) XPending(ctx context.Context, key string, group string) (*kvx.PendingInfo, error) {
	result, err := a.client.XPending(ctx, key, group).Result()
	if err != nil {
		return nil, err
	}

	return &kvx.PendingInfo{
		Count:     result.Count,
		StartID:   result.Lower,
		EndID:     result.Higher,
		Consumers: result.Consumers,
	}, nil
}

// XPendingRange gets pending entries in a range.
func (a *Adapter) XPendingRange(ctx context.Context, key string, group string, start string, stop string, count int64) ([]kvx.PendingEntry, error) {
	result, err := a.client.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: key,
		Group:  group,
		Start:  start,
		End:    stop,
		Count:  count,
	}).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]kvx.PendingEntry, len(result))
	for i, p := range result {
		entries[i] = kvx.PendingEntry{
			ID:         p.ID,
			Consumer:   p.Consumer,
			IdleTime:   p.Idle,
			Deliveries: p.RetryCount,
		}
	}
	return entries, nil
}

// XClaim claims pending entries for a consumer.
func (a *Adapter) XClaim(ctx context.Context, key string, group string, consumer string, minIdleTime time.Duration, ids []string) ([]kvx.StreamEntry, error) {
	result, err := a.client.XClaim(ctx, &redis.XClaimArgs{
		Stream:   key,
		Group:    group,
		Consumer: consumer,
		MinIdle:  minIdleTime,
		Messages: ids,
	}).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]kvx.StreamEntry, len(result))
	for i, msg := range result {
		entries[i] = kvx.StreamEntry{
			ID:     msg.ID,
			Values: convertInterfaceMapToBytes(msg.Values),
		}
	}
	return entries, nil
}

// XAutoClaim auto-claims pending entries.
func (a *Adapter) XAutoClaim(ctx context.Context, key string, group string, consumer string, minIdleTime time.Duration, start string, count int64) (string, []kvx.StreamEntry, error) {
	messages, next, err := a.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   key,
		Group:    group,
		Consumer: consumer,
		MinIdle:  minIdleTime,
		Start:    start,
		Count:    count,
	}).Result()
	if err != nil {
		return "", nil, err
	}

	entries := make([]kvx.StreamEntry, len(messages))
	for i, msg := range messages {
		entries[i] = kvx.StreamEntry{
			ID:     msg.ID,
			Values: convertInterfaceMapToBytes(msg.Values),
		}
	}
	return next, entries, nil
}

// XInfoGroups gets info about consumer groups.
func (a *Adapter) XInfoGroups(ctx context.Context, key string) ([]kvx.GroupInfo, error) {
	result, err := a.client.XInfoGroups(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	groups := make([]kvx.GroupInfo, len(result))
	for i, g := range result {
		groups[i] = kvx.GroupInfo{
			Name:            g.Name,
			Consumers:       g.Consumers,
			Pending:         g.Pending,
			LastDeliveredID: g.LastDeliveredID,
		}
	}
	return groups, nil
}

// XInfoConsumers gets info about consumers in a group.
func (a *Adapter) XInfoConsumers(ctx context.Context, key string, group string) ([]kvx.ConsumerInfo, error) {
	result, err := a.client.XInfoConsumers(ctx, key, group).Result()
	if err != nil {
		return nil, err
	}

	consumers := make([]kvx.ConsumerInfo, len(result))
	for i, c := range result {
		consumers[i] = kvx.ConsumerInfo{
			Name:    c.Name,
			Pending: c.Pending,
			Idle:    c.Idle,
		}
	}
	return consumers, nil
}

// XInfoStream gets info about a stream.
func (a *Adapter) XInfoStream(ctx context.Context, key string) (*kvx.StreamInfo, error) {
	result, err := a.client.XInfoStream(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	info := &kvx.StreamInfo{
		Length:          result.Length,
		RadixTreeKeys:   result.RadixTreeKeys,
		RadixTreeNodes:  result.RadixTreeNodes,
		Groups:          result.Groups,
		LastGeneratedID: result.LastGeneratedID,
	}

	if result.FirstEntry.ID != "" {
		info.FirstEntry = &kvx.StreamEntry{
			ID:     result.FirstEntry.ID,
			Values: convertInterfaceMapToBytes(result.FirstEntry.Values),
		}
	}

	if result.LastEntry.ID != "" {
		info.LastEntry = &kvx.StreamEntry{
			ID:     result.LastEntry.ID,
			Values: convertInterfaceMapToBytes(result.LastEntry.Values),
		}
	}

	return info, nil
}
