package health

type Health struct {
	Name  string
	Alive bool
}

type TransportHealth struct {
	Name  string `json:"name"`
	Alive bool   `json:"alive"`
}
type HealthResponse struct {
	Services []TransportHealth `json:"services"`
}
