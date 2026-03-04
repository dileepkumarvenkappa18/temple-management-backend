package donation

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type entityBankRepo struct {
	db *gorm.DB
}

func NewEntityBankRepo(db *gorm.DB) EntityRepository {
	return &entityBankRepo{db: db}
}

// row is the internal scan target for all query paths
type row struct {
	AccountHolderName string `gorm:"column:account_holder_name"`
	AccountNumber     string `gorm:"column:account_number"`
	IFSCCode          string `gorm:"column:ifsc_code"`
	UPIID             string `gorm:"column:upi_id"`
	BankName          string `gorm:"column:bank_name"`
	BranchName        string `gorm:"column:branch_name"`
	AccountType       string `gorm:"column:account_type"`
	RazorpayKeyID     string `gorm:"column:razorpay_key_id"`
	RazorpaySecret    string `gorm:"column:razorpay_secret"`
}

// GetBankDetailsByEntityID returns the TENANT's bank details including Razorpay credentials.
//
// Strategy (in order):
//  1. Direct: entities.created_by → tenant_bank_account_details.user_id
//  2. Fallback: first row in tenant_bank_account_details (ordered by user_id)
func (r *entityBankRepo) GetBankDetailsByEntityID(ctx context.Context, entityID uint) (*EntityBankDetails, error) {
	var result row

	// ── Path 1: entities.created_by → tenant_bank_account_details ────────
	// Note: No deleted_at filter — entities table does not have deleted_at column
	r.db.WithContext(ctx).Raw(`
		SELECT
			tba.account_holder_name,
			tba.account_number,
			tba.ifsc_code,
			COALESCE(tba.upi_id, '')              AS upi_id,
			COALESCE(tba.bank_name, '')           AS bank_name,
			COALESCE(tba.branch_name, '')         AS branch_name,
			COALESCE(tba.account_type, '')        AS account_type,
			COALESCE(tba.razorpay_key_id, '')     AS razorpay_key_id,
			COALESCE(tba.razorpay_secret, '')     AS razorpay_secret
		FROM entities e
		INNER JOIN tenant_bank_account_details tba ON tba.user_id = e.created_by
		WHERE e.id = ?
		LIMIT 1
	`, entityID).Scan(&result)

	if result.AccountHolderName != "" {
		log.Printf("✅ [path1] entity=%d holder=%s upi=%s razorpay_key=%s",
			entityID, result.AccountHolderName, result.UPIID, maskKey(result.RazorpayKeyID))
		return toDetails(result), nil
	}

	log.Printf("⚠️  [path1] no match for entity=%d, trying path2...", entityID)

	// ── Path 2: Last resort — first row in tenant_bank_account_details ───
	r.db.WithContext(ctx).Raw(`
		SELECT
			tba.account_holder_name,
			tba.account_number,
			tba.ifsc_code,
			COALESCE(tba.upi_id, '')              AS upi_id,
			COALESCE(tba.bank_name, '')           AS bank_name,
			COALESCE(tba.branch_name, '')         AS branch_name,
			COALESCE(tba.account_type, '')        AS account_type,
			COALESCE(tba.razorpay_key_id, '')     AS razorpay_key_id,
			COALESCE(tba.razorpay_secret, '')     AS razorpay_secret
		FROM tenant_bank_account_details tba
		ORDER BY tba.user_id ASC
		LIMIT 1
	`).Scan(&result)

	if result.AccountHolderName != "" {
		log.Printf("✅ [path2:first-row] entity=%d holder=%s upi=%s razorpay_key=%s",
			entityID, result.AccountHolderName, result.UPIID, maskKey(result.RazorpayKeyID))
		return toDetails(result), nil
	}

	log.Printf("❌ ALL paths failed for entity=%d — tenant_bank_account_details is empty?", entityID)
	return nil, fmt.Errorf("no tenant bank details found for entity %d", entityID)
}

// toDetails converts a scanned row to EntityBankDetails
func toDetails(r row) *EntityBankDetails {
	return &EntityBankDetails{
		AccountHolderName: r.AccountHolderName,
		AccountNumber:     r.AccountNumber,
		IFSCCode:          r.IFSCCode,
		UPIID:             r.UPIID,
		BankName:          r.BankName,
		BranchName:        r.BranchName,
		AccountType:       r.AccountType,
		RazorpayKeyID:     r.RazorpayKeyID,
		RazorpaySecret:    r.RazorpaySecret,
	}
}