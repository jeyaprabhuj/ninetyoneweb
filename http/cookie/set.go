package cookie

import (
	"github.com/jeyaprabhuj/ninetyoneweb/util/encoder"
	"github.com/jeyaprabhuj/ninetyoneweb/util/protection"
	"net/http"
)

func (c *Cookie) Set(name string, value string, w http.ResponseWriter, r *http.Request) {

	cryptKey := []byte(c.secret)
	crypt, _ := protection.Encrypt([]byte(value), cryptKey)

	newCookie := http.Cookie{
		Name:     name,
		Value:    encoder.Base64EncodeToString(crypt),
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &newCookie)
}
