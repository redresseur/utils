package utime

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

// CreateUtcTimestamp returns a google/protobuf/Timestamp in UTC
func CreateUtcTimestamp() *timestamp.Timestamp {
	now := time.Now().UTC()
	secs := now.Unix()
	nanos := int32(now.UnixNano() - (secs * 1000000000))
	return &(timestamp.Timestamp{Seconds: secs, Nanos: nanos})
}

func GenerateTimestamp(nowTime time.Time) *timestamp.Timestamp {
	secs := nowTime.Unix()
	nanos := int32(nowTime.UnixNano() - (secs * 1000000000))
	return &(timestamp.Timestamp{Seconds: secs, Nanos: nanos})
}