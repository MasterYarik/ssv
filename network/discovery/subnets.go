package discovery

import (
	"github.com/bloxapp/ssv/utils/format"
	"strconv"
)

const (
	// SubnetsCount is the count of subnets in the network
	SubnetsCount = 128
)

var regPool = format.NewRegexpPool("\\w+:bloxstaking\\.ssv\\.(\\d+)")

// nsToSubnet converts the given topic to subnet
// TODO: return other value than zero upon failure?
func nsToSubnet(ns string) int64 {
	r, done := regPool.Get()
	defer done()
	found := r.FindStringSubmatch(ns)
	if len(found) != 2 {
		return -1
	}
	val, err := strconv.ParseUint(found[1], 10, 64)
	if err != nil {
		return -1
	}
	return int64(val)
}

// isSubnet checks if the given string is a subnet string
func isSubnet(ns string) bool {
	r, done := regPool.Get()
	defer done()
	return r.MatchString(ns)
}
