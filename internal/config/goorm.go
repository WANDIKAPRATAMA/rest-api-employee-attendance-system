package config

import (
	"employee-attendance-system/internal/entity/domain"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
postgres=# CREATE DATABASE db_employee_attendance_system;
CREATE DATABASE
postgres=# \c db_employee_attendance_system;
You are now connected to database "db_employee_attendance_system" as user "postgres".
db_employee_attendance_system=#
db_employee_attendance_system=# CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION
db_employee_attendance_system=#
*/
func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")
	sslmode := viper.GetString("database.sslmode")
	timezone := viper.GetString("database.timezone")
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	// Format DSN PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		host, username, password, database, port, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Connection pool setup
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}
	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))
	enumTypes := map[string][]string{
		"user_status": {
			"active",
			"inactive",
			"banned",
		},
		"app_role": {
			string(domain.Admin),
			string(domain.Employee),
		},

		"attendance_type": {
			string(domain.AttendanceTypeIn),
			string(domain.AttendanceTypeOut),
		},
	}
	for typeName, values := range enumTypes {
		var exists bool
		err := db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?)", typeName).Scan(&exists).Error
		if err != nil {
			log.Printf("gagal memeriksa tipe %s: %v", typeName, err)
		}
		if !exists {
			valuesStr := "'" + strings.Join(values, "', '") + "'"
			err = db.Exec(fmt.Sprintf("CREATE TYPE %s AS ENUM (%s)", typeName, valuesStr)).Error
			if err != nil {
				log.Printf("Gagal membuat tipe %s: %v", typeName, err)
			}
		}
	}

	// Auto migrate models (urutan penting: parent dulu)
	err = db.AutoMigrate(
		&domain.User{},
		&domain.UserProfile{},
		&domain.UserSecurity{},
		&domain.ApplicationRole{},
		&domain.RefreshToken{},
		&domain.Department{},
		&domain.Attendance{},
		&domain.AttendanceHistory{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	indexCommands := []struct {
		tableName string
		sql       string
	}{
		{
			tableName: "departments",
			sql: `
                ALTER TABLE departments ALTER COLUMN max_clock_in_time TYPE TIME USING max_clock_in_time::time,
				ALTER TABLE departments ALTER COLUMN max_clock_out_time TYPE TIME USING max_clock_out_time::time,

            `,
		},
	}

	for _, cmd := range indexCommands {
		log.Printf("Menjalankan indeks untuk tabel %s", cmd.tableName)
		if err := db.Exec(cmd.sql).Error; err != nil {
			log.Printf("Gagal membuat indeks untuk tabel %s: %v", cmd.tableName, err)
			log.Printf("Gagal menambahkan indeks: %v", err)
		}
		log.Printf("Berhasil membuat indeks untuk tabel %s", cmd.tableName)
	}
	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
