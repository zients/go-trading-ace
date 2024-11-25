package entities

import "time"

type Task struct {
	ID          int64      `db:"id"`          // SERIAL PRIMARY KEY
	Name        string     `db:"name"`        // VARCHAR(255) NOT NULL
	Description string     `db:"description"` // TEXT
	Points      float64    `db:"points"`      // BIGINT NOT NULL
	StartedAt   *time.Time `db:"started_at"`  // TIMESTAMP NULL
	EndAt       *time.Time `db:"end_at"`      // TIMESTAMP NULL
	Period      int        `db:"period"`      // INT DEFAULT 1
	CreatedAt   time.Time  `db:"created_at"`  // TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	UpdatedAt   time.Time  `db:"updated_at"`  // TIMESTAMP DEFAULT CURRENT_TIMESTAMP
}
