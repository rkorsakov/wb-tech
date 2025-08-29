package ntpclient

import (
	"github.com/beevik/ntp"
	"time"
)

const NTPServer = "0.beevik-ntp.pool.ntp.org"

func GetCurrentTime() (time.Time, error) {
	response, err := ntp.Query(NTPServer)
	if err != nil {
		return time.Time{}, err
	}
	err = response.Validate()
	if err != nil {
		return time.Time{}, err
	}
	exactTime := time.Now().Add(response.ClockOffset)
	return exactTime, nil
}
