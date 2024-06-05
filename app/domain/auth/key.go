package domainauth

const (
	memcacheKeyGetUserById       = "domain:user:id:%s"
	memcacheKeyGetUserByUsername = "domain:user:username:%s"
)

const (
	singleFlightKeyGetUserById       = "sf:domain:user:id:%s"
	singleFlightKeyGetUserByUsername = "sf:domain:user:username:%s"
)
