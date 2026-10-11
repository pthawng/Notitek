// Package idgen creates identifiers for Notitek resources.
package idgen

import (
	"fmt"

	"github.com/google/uuid"
)

// New returns a time-ordered UUIDv7 string. Time ordering keeps B-tree
// indexes on primary keys compact and makes ids roughly sortable by creation.
func New() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("idgen: generate uuidv7: %w", err)
	}
	return id.String(), nil
}
