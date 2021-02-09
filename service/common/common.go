package common

import "encoding/base64"

// BasicAuth will add basic auth to connection
func BasicAuth(user, pass string) string {
	auth := user + ":" + pass
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
