package serializer

import (
	"github.com/bytedance/sonic"
)

type Config struct {
	SonicConfig sonic.Config
}

func DefaultConfig() Config {
	return Config{
		SonicConfig: sonic.Config{
			// The string is very likely to be used in someplace else.
			// This package is going to be used on a backend, and on the backend variables
			// can be stored for a long time. This means the string and its JSON buffer
			// might not be deallocated for a long time. On the backend, we never want to
			// excessively and redundantly consume memory. So, this is enabled.
			CopyString: true,
			// HTTP response needs to be kept as little as possible.
			// In this context, it is better to sacrifice CPU speed for network speed/latency.
			CompactMarshaler: true,
			// Security. Definitely enable this.
			EscapeHTML: true,
			// Redundant, and the last thing needed on this Earth.
			SortMapKeys: false,
		},
	}
}
