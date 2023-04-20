package access

import (
	"time"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

const stmtOAuth = `
SELECT access_token,
    is_active,
    expires_in,
    created_utc
FROM oauth.access
WHERE access_token = UNHEX(?)
LIMIT 1`

type OAuth struct {
	Token     string      `db:"access_token"`
	Active    bool        `db:"is_active"`
	ExpiresIn null.Int    `db:"expires_in"` // seconds
	CreatedAt chrono.Time `db:"created_utc"`
}

func (o OAuth) Expired() bool {

	if o.ExpiresIn.IsZero() {
		return false
	}

	expireAt := o.CreatedAt.Add(time.Second * time.Duration(o.ExpiresIn.Int64))

	if expireAt.Before(time.Now()) {
		return true
	}

	return false
}
