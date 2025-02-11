package config

import (
	"fmt"
	"go_base/internal/models"
	"go_base/internal/repositories"
	"go_base/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	dbHost := utils.GetEnv("DB_HOST", "postgres-base")
	dbUser := utils.GetEnv("DB_USER", "postgres")
	dbPassword := utils.GetEnv("DB_PASSWORD", "postgres")
	dbName := utils.GetEnv("DB_NAME", "postgres")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbSSL := utils.GetEnv("DB_SSL", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbSSL)
	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{}, &models.UserRole{})
	dbSeed(db)

	return db
}

func dbSeed(db *gorm.DB) {

	err := db.First(&models.UserRole{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			var userRoles []models.UserRole
			userRoles = append(userRoles, models.UserRole{
				Id:   1,
				Name: "Super Admin",
			})
			userRoles = append(userRoles, models.UserRole{
				Id:   2,
				Name: "HR",
			})

			userRoles = append(userRoles, models.UserRole{
				Id:   3,
				Name: "Employee",
			})
			db.Create(&userRoles)
		}
	}

	err = db.First(&models.User{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userRepo := repositories.NewUserRepository(db)
			var user []models.User

			hashedPassword, err := userRepo.HashPassword("rahasia")
			if err != nil {
				fmt.Println(err.Error())
			}

			user = append(user, models.User{
				Id:        uuid.NewString(),
				Name:      "Super Admin",
				Username:  "superadmin",
				Email:     "superadmin@gmail.com",
				Password:  hashedPassword,
				RoleId:    1,
				CreatedAt: time.Now(),
			})

			err = db.Create(&user).Error
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println("Created super admin")
		}
	}
}
