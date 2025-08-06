package superadmin

type TenantApprovalCount struct {
	Approved int64 `json:"approved"`
	Pending  int64 `json:"pending"`
	Rejected int64 `json:"rejected"`
}

// ================ TEMPLE APPROVAL COUNTS ================

type TempleApprovalCount struct {
	PendingApproval int64 `json:"pending_approval"`
	ActiveTemples   int64 `json:"active_temples"`
	Rejected        int64 `json:"rejected"`
	TotalDevotees   int64 `json:"total_devotees"`
}