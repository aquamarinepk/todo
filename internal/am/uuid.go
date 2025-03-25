package am

import (
	"database/sql"

	"github.com/google/uuid"
)

// ParseUUID safely parses a UUID string
func ParseUUID(s sql.NullString) uuid.UUID {
	if s.Valid {
		u, err := uuid.Parse(s.String)
		if err != nil {
			return uuid.Nil
		}
		return u
	}
	return uuid.Nil
}
