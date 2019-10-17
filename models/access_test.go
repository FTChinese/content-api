package models

import (
	"strings"
	"testing"
)

func TestExtractFromHeader(t *testing.T) {
	auth := "Bearer cc41e4afc7e74c42e5d522d75bd0ec3f984ba2b3"

	split := strings.SplitN(auth, " ", 2)
	t.Logf("%d", len(split))

	split = strings.Split(auth, " ")
	t.Logf("%d", len(split))
}
