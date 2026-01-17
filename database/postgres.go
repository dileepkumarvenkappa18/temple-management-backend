package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/donation"
	"github.com/sharath018/temple-management-backend/internal/entity"
	"github.com/sharath018/temple-management-backend/internal/event"
	"github.com/sharath018/temple-management-backend/internal/seva"
	"github.com/sharath018/temple-management-backend/internal/userprofile"
	"github.com/sharath018/temple-management-backend/internal/notification"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
)

var DB *gorm.DB

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	var err error
	for i := 1; i <= 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // 🆕 Enable SQL logging
		})
		if err == nil {
			log.Println("✅ Connected to database")
			break
		}
		log.Printf("⚠️  DB connection attempt %d/5 failed: %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("❌ Could not connect to database: %v", err)
	}

	// 🆕 Log current database name
	var currentDB string
	DB.Raw("SELECT current_database()").Scan(&currentDB)
	log.Printf("🔍 Connected to database: %s", currentDB)

	// ✅ Migrate all required models in correct order
	log.Println("🔄 Starting database migration...")
	
	if err := DB.AutoMigrate(
		// Auth models (order matters due to foreign keys)
		&auth.UserRole{},
		&auth.User{},
		&auth.ApprovalRequest{},
		&auth.TenantDetails{},
		&auth.BankAccountDetails{}, // 🆕 Bank details
		&auth.TenantUserAssignment{},
		
		// Other models
		&seva.Seva{},
		&seva.SevaBooking{},
		&entity.Entity{},
		&event.Event{},
		&donation.Donation{},
		&notification.NotificationTemplate{},
		&notification.NotificationLog{},
		&userprofile.DevoteeProfile{},
		&userprofile.Child{},
		&userprofile.EmergencyContact{},
		&userprofile.UserEntityMembership{},
		&auditlog.AuditLog{},
	); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	// 🆕 Verify critical tables exist
	verifyTables(DB)

	log.Println("✅ Database schema migrated successfully")

	// 🌱 Seed default roles
	if err := auth.SeedUserRoles(DB); err != nil {
		log.Fatalf("❌ Seeding roles failed: %v", err)
	}

	return DB
}

func verifyTables(db *gorm.DB) {
	tables := []string{
		"user_roles",
		"users",
		"tenant_details",
		"bank_account_details", // 🔍 Check this specifically
		"approval_requests",
	}

	log.Println("🔍 Verifying critical tables...")
	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			log.Printf("  ✅ Table '%s' exists", table)
		} else {
			log.Printf("  ❌ Table '%s' MISSING!", table)
		}
	}
}