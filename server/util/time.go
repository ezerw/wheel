package util

import "time"

// GetNextWorkingDay returns the next working day based in the day passed as parameter.
func GetNextWorkingDay() (*time.Time, error) {
	var next time.Time

	loc, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		return nil, err
	}

	today := time.Now().In(loc)

	year, month, day := today.Date()
	tomorrow := time.Date(year, month, day+1, 0, 0, 0, 0, loc)

	switch tomorrow.Weekday() {
	case time.Saturday:
		next = tomorrow.Add(48 * time.Hour)
	case time.Sunday:
		next = tomorrow.Add(24 * time.Hour)
	default:
		next = tomorrow
	}

	return &next, nil
}
