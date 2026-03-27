package dbx

import "time"

func registerBuiltinCodecs(registry *codecRegistry) {
	if registry == nil {
		return
	}
	registry.mustRegister(jsonCodec{})
	registry.mustRegister(textCodec{})
	registry.mustRegister(newTimeStringCodec("rfc3339_time", time.RFC3339))
	registry.mustRegister(newTimeStringCodec("rfc3339nano_time", time.RFC3339Nano))
	registry.mustRegister(newUnixTimeCodec("unix_time", unixSeconds))
	registry.mustRegister(newUnixTimeCodec("unix_milli_time", unixMillis))
	registry.mustRegister(newUnixTimeCodec("unix_nano_time", unixNanos))
}
