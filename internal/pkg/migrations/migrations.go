package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/rs/zerolog/log"
)

func StartMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Info().
			Err(err).
			Msg("Ошибка создания драйвера миграции ")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к миграциям
		"postgres",          // Имя базы данных
		driver,
	)
	if err != nil {
		log.Info().
			Err(err).
			Msg("Ошибка инициализации миграции ")
	}

	// Выполняем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Info().
			Err(err).
			Msg("Ошибка выполнения миграции ")
	} else {
		fmt.Println("Миграции применены успешно!")
	}
}
