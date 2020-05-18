package cookie

import (
	"github.com/jeyaprabhuj/ninetyoneweb/util/encoder"
	"github.com/jeyaprabhuj/ninetyoneweb/util/protection"
	"net/http"
)

func (c *Cookie) GetValue(name string, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, _ := r.Cookie(name)
	cryptKey := []byte(c.secret)
	cipher, err := encoder.Base64DecodeToByte(cookie.Value)
	if err != nil {
		return "", err
	}
	byteValue, error := protection.Decrypt(cipher, cryptKey)

	return string(byteValue), error
}
