package databases

import (
	"fmt"

	"github.com/thienhi/fusionstart/internal/configs"
	"github.com/thienhi/fusionstart/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) bool {
	err := db.AutoMigrate(
		&models.BaseModel{},
		&models.Booking{},
		&models.User{},
		&models.Event{},
	)

	if err != nil {
		fmt.Println("Database migration error.", err)
		return false
	}
	fmt.Println("Database migration completed.")
	return true
}

func ConnectDatabase(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.Port,
		config.Database.SSLMode,
	)
	fmt.Println("dsn", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if migrated := MigrateDatabase(db); migrated != true {
		return nil, fmt.Errorf("failed to migrate database")
	}
	return db, nil
}

func CloseConnectionDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
