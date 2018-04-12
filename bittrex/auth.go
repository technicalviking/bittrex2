package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

func (c *Client) sign(challenge string) string {
	//hasher := hmac.New(sha512.New, []byte())
	hasher := hmac.New(sha512.New, []byte(c.apiSecret))
	hasher.Write([]byte(challenge))

	return hex.EncodeToString(hasher.Sum(nil))
}
