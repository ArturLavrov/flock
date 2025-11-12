package shortID

import (
	"github.com/google/uuid"
	"strings"
)

// Generate returns a short, unique string
func Generate() string {
	id := uuid.New()

	return strings.Split(id.String(), "-")[0]
}
