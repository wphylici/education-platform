package database

import (
	"fmt"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	config *Config
	db     *gorm.DB
}

func NewPostgreSQL(config *Config) *PostgreSQL {
	return &PostgreSQL{
		config: config,
	}
}

func (p *PostgreSQL) Connect() error {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		p.config.Host,
		p.config.User,
		p.config.Password,
		p.config.Database,
		p.config.Port,
		p.config.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	p.db = db

	return nil
}

func (p *PostgreSQL) InitDB() error {

	tx := p.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if tx.Error != nil {
		return tx.Error
	}

	err := p.db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
