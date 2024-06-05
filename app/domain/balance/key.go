package domainbalance

const (
	memcacheKeyGetBalanceByUserId               = "domain:balance:user_id:%s"
	memcacheKeyGetLatestHistoryByUserId         = "domain:balance:history:user_id:%s"
	memcacheKeyGetHistorySummaryByUserIdAndType = "domain:balance:history_summary:user_id:%s:type:%d"
)

const (
	singleFlightKeyGetBalanceByUserId               = "sf:domain:balance:user_id:%s"
	singleFlightKeyGetLatestHistoryByUserId         = "sf:domain:balance:history:user_id:%s"
	singleFlightKeyGetHistorySummaryByUserIdAndType = "sf:domain:balance:history_summary:user_id:%s:type:%d"
)
