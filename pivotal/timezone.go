package pivotal

type TimeZone struct {
	Kind      string `json:"kind"`
	OlsonName string `json:"olson_name"`
	Offset    string `json:"offset"`
}
