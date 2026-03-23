package model

import "time"

// TelemetryPayload representa o formato importante que vem do K6
// DICA: usar json:"" nas tags, senão o Marshal se confunde e não cruza os dados com o body
type TelemetryPayload struct {
	IdDispositivo   string    `json:"device_id"`
	HoraEData       time.Time `json:"timestamp"`
	TipoSensor      string    `json:"sensor_type"`
	TipoDeLeitura   string    `json:"reading_type"`
	ValorLido       float64   `json:"value"`
}
