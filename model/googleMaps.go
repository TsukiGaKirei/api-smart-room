package model

// DistanceMatrixResponse represents a Distance Matrix API response.
type DistanceMatrixResponse struct {
	DestinationAddresses []string                    `json:"destination_addresses"`
	OriginAddresses      []string                    `json:"origin_addresses"`
	Rows                 []DistanceMatrixElementsRow `json:"rows"`
	Status               string                      `json:"status"`
}

// DistanceMatrixElementsRow is a row of distance elements.
type DistanceMatrixElementsRow struct {
	Elements []*DistanceMatrixElement `json:"elements"`
}

// DistanceMatrixElement is the travel distance and time for a pair of origin
// and destination.
type DistanceMatrixElement struct {
	Status string `json:"status"`
	// Duration is the length of time it takes to travel this route.
	Duration Duration `json:"duration"`
	Distance Distance `json:"distance"`
}

type Duration struct {
	Text  string `json:"string"`
	Value int64  `json:"value"`
}

type Distance struct {
	// HumanReadable is the human friendly distance. This is rounded and in an
	// appropriate unit for the request. The units can be overriden with a request
	// parameter.
	HumanReadable string `json:"text"`
	// Meters is the numeric distance, always in meters. This is intended to be used
	// only in algorithmic situations, e.g. sorting results by some user specified
	// metric.
	Meters int `json:"value"`
}
