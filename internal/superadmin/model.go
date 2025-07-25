package superadmin

type DashboardMetrics struct {
	PendingApprovals int64 `json:"pending_approvals"`
	ActiveTemples    int64 `json:"active_temples"`
	TotalUsers       int64 `json:"total_users"`
	RejectedCount    int64 `json:"rejected_count"`
}