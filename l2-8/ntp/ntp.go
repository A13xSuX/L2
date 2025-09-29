package ntp

import (
	"github.com/beevik/ntp"
)

func GetCurrentTime() (string, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return "", err
	}
	return time.String(), nil
}
