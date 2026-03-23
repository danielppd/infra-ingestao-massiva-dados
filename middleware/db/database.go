package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ponderada/middleware/model"
)

// Database só engloba a nossa conexão crua do pacote sql
type Database struct {
	conexao *sql.DB
}

// NewDatabase tenta abrir a conexão usando o DSN sinistro do postgres
func NewDatabase(urlConexao string) (*Database, error) {
	banco, erro := sql.Open("postgres", urlConexao)
	if erro != nil {
		return nil, fmt.Errorf("não deu pra instanciar o sql: %w", erro)
	}

	// Manda um ping pra testar se ouviu
	if erro := banco.Ping(); erro != nil {
		return nil, fmt.Errorf("não deu pra conectar no banco: %w", erro)
	}

	return &Database{conexao: banco}, nil
}

// InsertReading faz o papel de enviar as paradinhas do struct pro banco
func (d *Database) InsertReading(dadoDaVez model.TelemetryPayload) error {
	
	// A tabela é 'sensor_readings' conforme tá lá no nosso script init.sql
	querySql := `INSERT INTO sensor_readings (device_id, timestamp, sensor_type, reading_type, value)
			  VALUES ($1, $2, $3, $4, $5)`
	
	_, erro := d.conexao.Exec(querySql, dadoDaVez.IdDispositivo, dadoDaVez.HoraEData, dadoDaVez.TipoSensor, dadoDaVez.TipoDeLeitura, dadoDaVez.ValorLido)
	return erro
}

// Fecha
func (d *Database) Close() {
	if d.conexao != nil {
		d.conexao.Close()
	}
}
