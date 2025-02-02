package events

type Meta struct {
	URI       string `json:"uri"`
	RequestID string `json:"request_id"`
	ID        string `json:"id"`
	Dt        string `json:"dt"`
	Domain    string `json:"domain"`
	Stream    string `json:"stream"`
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Offset    int    `json:"offset"`
}

type EventData struct {
	Schema string `json:"$schema"`
}
