package domainbalance

const (
	cacheKeyGetBalanceByUserId               = "domain:balance:user_id:%s"
	cacheKeyGetLatestHistoryByUserId         = "domain:balance:history:user_id:%s"
	cacheKeyGetHistorySummaryByUserIdAndType = "domain:balance:history_summary:user_id:%s:type:%d"
)

const (
	singleFlightKeyGetBalanceByUserId               = "sf:domain:balance:user_id:%s"
	singleFlightKeyGetLatestHistoryByUserId         = "sf:domain:balance:history:user_id:%s"
	singleFlightKeyGetHistorySummaryByUserIdAndType = "sf:domain:balance:history_summary:user_id:%s:type:%d"
)
