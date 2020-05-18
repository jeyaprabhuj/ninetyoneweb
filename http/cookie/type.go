package cookie

type Cookie struct {
	secret string
}

func NewCookie(secret string) *Cookie {
	return &Cookie{secret: secret}
}
