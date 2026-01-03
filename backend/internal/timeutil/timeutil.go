package timeutil

import "time"

var (
	locationTokyo *time.Location
)

func init() {
	locationTokyo = mustLoadLocation("Asia/Tokyo")
}

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

// Tokyo returns the Asia/Tokyo location.
func Tokyo() *time.Location {
	return locationTokyo
}

// NowTokyo returns current time in Asia/Tokyo.
func NowTokyo() time.Time {
	return time.Now().In(locationTokyo)
}
