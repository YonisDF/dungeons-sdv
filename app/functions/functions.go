package functions

import (
	"encoding/json"
	"math"
	"strings"

	"github.com/gofrs/uuid"
)

// NewUUID ...
func NewUUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

// RemoveDuplicate : Supression de doublons dans un tableau
func RemoveDuplicate(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		x = strings.TrimPrefix(x, "-")
		if !found[strings.ToLower(x)] {
			found[strings.ToLower(x)] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

// ConvertInputStructToDataStruct caution data must be an absolute pointer to work
func ConvertInputStructToDataStruct(input interface{}, data interface{}) error {
	// Sérialize Input in JSON
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	// Désérialize JSON in data
	err = json.Unmarshal(inputBytes, &data)
	if err != nil {
		return err
	}

	return nil
}

// Harvesine distance calculation
func HaversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0

	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	lat1Rad := toRad(lat1)
	lat2Rad := toRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
