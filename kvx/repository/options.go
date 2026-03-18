package repository

import (
	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/DaiYuANg/arcgo/kvx/mapping"
	"github.com/samber/mo"
)

type hashRepositoryConfig[T any] struct {
	keyBuilder *mapping.KeyBuilder
	tagParser  *mapping.TagParser
	codec      *mapping.HashCodec
	indexer    *Indexer[T]
	pipeline   mo.Option[pipelineProvider]
}

type jsonRepositoryConfig[T any] struct {
	keyBuilder *mapping.KeyBuilder
	tagParser  *mapping.TagParser
	serializer mapping.Serializer
	indexer    *Indexer[T]
	pipeline   mo.Option[pipelineProvider]
}

type HashRepositoryOption[T any] interface {
	applyHash(*hashRepositoryConfig[T])
}
type JSONRepositoryOption[T any] interface {
	applyJSON(*jsonRepositoryConfig[T])
}

type hashOptionFunc[T any] func(*hashRepositoryConfig[T])

func (f hashOptionFunc[T]) applyHash(cfg *hashRepositoryConfig[T]) { f(cfg) }

type jsonOptionFunc[T any] func(*jsonRepositoryConfig[T])

func (f jsonOptionFunc[T]) applyJSON(cfg *jsonRepositoryConfig[T]) { f(cfg) }

type dualOption[T any] struct {
	hash hashOptionFunc[T]
	json jsonOptionFunc[T]
}

func (o dualOption[T]) applyHash(cfg *hashRepositoryConfig[T]) {
	if o.hash != nil {
		o.hash(cfg)
	}
}
func (o dualOption[T]) applyJSON(cfg *jsonRepositoryConfig[T]) {
	if o.json != nil {
		o.json(cfg)
	}
}

// OptionSet lets callers reuse a standard repository configuration across entities.
type OptionSet[T any] struct {
	Hash []HashRepositoryOption[T]
	JSON []JSONRepositoryOption[T]
}

func (s OptionSet[T]) HashOptions(extra ...HashRepositoryOption[T]) []HashRepositoryOption[T] {
	result := append([]HashRepositoryOption[T]{}, s.Hash...)
	return append(result, extra...)
}

func (s OptionSet[T]) JSONOptions(extra ...JSONRepositoryOption[T]) []JSONRepositoryOption[T] {
	result := append([]JSONRepositoryOption[T]{}, s.JSON...)
	return append(result, extra...)
}

// Preset groups reusable repository options.
type Preset[T any] struct {
	Hash []HashRepositoryOption[T]
	JSON []JSONRepositoryOption[T]
}

func (p Preset[T]) HashOptions(extra ...HashRepositoryOption[T]) []HashRepositoryOption[T] {
	result := append([]HashRepositoryOption[T]{}, p.Hash...)
	return append(result, extra...)
}

func (p Preset[T]) JSONOptions(extra ...JSONRepositoryOption[T]) []JSONRepositoryOption[T] {
	result := append([]JSONRepositoryOption[T]{}, p.JSON...)
	return append(result, extra...)
}

func defaultHashConfig[T any](kv kvx.KV, keyPrefix string) hashRepositoryConfig[T] {
	return hashRepositoryConfig[T]{
		keyBuilder: mapping.NewKeyBuilder(keyPrefix),
		tagParser:  mapping.NewTagParser(),
		codec:      mapping.NewHashCodec(nil),
		indexer:    NewIndexer[T](kv, keyPrefix),
		pipeline:   mo.None[pipelineProvider](),
	}
}

func defaultJSONConfig[T any](kv kvx.KV, keyPrefix string) jsonRepositoryConfig[T] {
	return jsonRepositoryConfig[T]{
		keyBuilder: mapping.NewKeyBuilder(keyPrefix),
		tagParser:  mapping.NewTagParser(),
		serializer: mapping.NewJSONSerializer(),
		indexer:    NewIndexer[T](kv, keyPrefix),
		pipeline:   mo.None[pipelineProvider](),
	}
}

func applyHashOptions[T any](cfg *hashRepositoryConfig[T], options ...HashRepositoryOption[T]) {
	for _, option := range options {
		if option != nil {
			option.applyHash(cfg)
		}
	}
}

func applyJSONOptions[T any](cfg *jsonRepositoryConfig[T], options ...JSONRepositoryOption[T]) {
	for _, option := range options {
		if option != nil {
			option.applyJSON(cfg)
		}
	}
}

func WithPipeline[T any](provider pipelineProvider) dualOption[T] {
	return dualOption[T]{
		hash: func(cfg *hashRepositoryConfig[T]) { cfg.pipeline = mo.Some[pipelineProvider](provider) },
		json: func(cfg *jsonRepositoryConfig[T]) { cfg.pipeline = mo.Some[pipelineProvider](provider) },
	}
}

func WithKeyBuilder[T any](builder *mapping.KeyBuilder) dualOption[T] {
	return dualOption[T]{
		hash: func(cfg *hashRepositoryConfig[T]) {
			if builder != nil {
				cfg.keyBuilder = builder
			}
		},
		json: func(cfg *jsonRepositoryConfig[T]) {
			if builder != nil {
				cfg.keyBuilder = builder
			}
		},
	}
}

func WithTagParser[T any](parser *mapping.TagParser) dualOption[T] {
	return dualOption[T]{
		hash: func(cfg *hashRepositoryConfig[T]) {
			if parser != nil {
				cfg.tagParser = parser
			}
		},
		json: func(cfg *jsonRepositoryConfig[T]) {
			if parser != nil {
				cfg.tagParser = parser
			}
		},
	}
}

func WithIndexer[T any](indexer *Indexer[T]) dualOption[T] {
	return dualOption[T]{
		hash: func(cfg *hashRepositoryConfig[T]) {
			if indexer != nil {
				cfg.indexer = indexer
			}
		},
		json: func(cfg *jsonRepositoryConfig[T]) {
			if indexer != nil {
				cfg.indexer = indexer
			}
		},
	}
}

func WithHashCodec[T any](codec *mapping.HashCodec) HashRepositoryOption[T] {
	return hashOptionFunc[T](func(cfg *hashRepositoryConfig[T]) {
		if codec != nil {
			cfg.codec = codec
		}
	})
}

func WithSerializer[T any](serializer mapping.Serializer) JSONRepositoryOption[T] {
	return jsonOptionFunc[T](func(cfg *jsonRepositoryConfig[T]) {
		if serializer != nil {
			cfg.serializer = serializer
		}
	})
}

func NewPreset[T any](options ...dualOption[T]) Preset[T] {
	hashOptions := make([]HashRepositoryOption[T], 0, len(options))
	jsonOptions := make([]JSONRepositoryOption[T], 0, len(options))
	for _, option := range options {
		hashOptions = append(hashOptions, option)
		jsonOptions = append(jsonOptions, option)
	}
	return Preset[T]{Hash: hashOptions, JSON: jsonOptions}
}
