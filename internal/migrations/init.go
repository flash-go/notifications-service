package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migration_notifications_init() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "notifications_init",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS email_folders (
					id SERIAL PRIMARY KEY,
					parent_id INTEGER REFERENCES email_folders(id) ON UPDATE CASCADE ON DELETE CASCADE,
					name TEXT NOT NULL,
					description TEXT NOT NULL,
					system_flag BOOLEAN NOT NULL,
					updated TIMESTAMPTZ NOT NULL,
					created TIMESTAMPTZ NOT NULL
				);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS emails (
					id SERIAL PRIMARY KEY,
					folder_id INTEGER REFERENCES email_folders(id) ON UPDATE CASCADE ON DELETE CASCADE,
					from_email TEXT NOT NULL,
					from_name TEXT NOT NULL,
					subject TEXT NOT NULL,
					html TEXT NOT NULL,
					text TEXT NOT NULL,
					description TEXT NOT NULL,
					system_flag BOOLEAN NOT NULL,
					updated TIMESTAMPTZ NOT NULL,
					created TIMESTAMPTZ NOT NULL
				);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`
				CREATE TABLE IF NOT EXISTS email_logs (
					id SERIAL PRIMARY KEY,
					from_email TEXT NOT NULL,
					from_name TEXT NOT NULL,
					subject TEXT NOT NULL,
					to_email TEXT NOT NULL,
					html TEXT NOT NULL,
					text TEXT NOT NULL,
					status TEXT NOT NULL,
					message_id TEXT,
					errors TEXT,
					created TIMESTAMPTZ NOT NULL
				);
			`).Error; err != nil {
				return err
			}

			if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_emails_folder_id ON emails(folder_id);`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP TABLE IF EXISTS emails;`).Error; err != nil {
				return err
			}
			return nil
		},
	}
}
