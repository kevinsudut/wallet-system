package domainauth

const (
	cacheKeyGetUserById       = "domain:user:id:%s"
	cacheKeyGetUserByUsername = "domain:user:username:%s"
)

const (
	singleFlightKeyGetUserById       = "sf:domain:user:id:%s"
	singleFlightKeyGetUserByUsername = "sf:domain:user:username:%s"
)
