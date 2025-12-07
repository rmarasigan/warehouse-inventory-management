package log

import "log/slog"

type Map map[string]any

type KeyValue struct {
	Key   string
	Value any
}

func (kv KeyValue) Attr() slog.Attr {
	return slog.Any(kv.Key, kv.Value)
}

func KV(key string, value any) KeyValue {
	return KeyValue{Key: key, Value: value}
}

func KVs(kvs map[string]any) []KeyValue {
	var merged []KeyValue

	for key, value := range kvs {
		merged = append(merged,
			KeyValue{
				Key:   key,
				Value: value,
			},
		)
	}

	return merged
}
