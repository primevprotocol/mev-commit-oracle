package transactor

func SetAllowedPendingTxnCount(count int) func() {
	original := allowedPendingTxnCount
	allowedPendingTxnCount = count
	return func() {
		allowedPendingTxnCount = original
	}
}
