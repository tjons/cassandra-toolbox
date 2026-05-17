package qb

type paramType string

const (
	ttlUpdateParam       paramType = "TTL"
	timestampUpdateParam paramType = "TIMESTAMP"
)

type updateParam struct {
	T   paramType
	Arg int64
}

func updateParamTTL(ttl int64) updateParam {
	return updateParam{
		T:   ttlUpdateParam,
		Arg: ttl,
	}
}

func updateParamTimestamp(timestamp int64) updateParam {
	return updateParam{
		T:   timestampUpdateParam,
		Arg: timestamp,
	}
}

var (
	// TTL is a helper function for setting the TTL of an update or insert query.
	TTL = updateParamTTL

	// Timestamp is a helper function for setting the timestamp of an update, insert or delete query.
	Timestamp = updateParamTimestamp
)
