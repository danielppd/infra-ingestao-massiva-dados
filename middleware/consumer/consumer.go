package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ponderada/middleware/db"
	"github.com/ponderada/middleware/model"
)

// Consumer tem tudo o que precisamos para ouvir o Rabbit e jogar no Postgres
type Consumer struct {
	conexaoRabbit  *amqp.Connection
	canalRabbit    *amqp.Channel
	conexaoBD *db.Database
}

func NewConsumer(urlRabbit string, bancoConectado *db.Database) (*Consumer, error) {
	var conexao *amqp.Connection
	var erro error
	
	// Feito um loping pra dar retry caso a fila ainda não tiver ligado
	for tentativa := 0; tentativa < 10; tentativa++ {
		conexao, erro = amqp.Dial(urlRabbit)
		if erro == nil {
			break
		}
		log.Printf("Rabbit não ligou ainda, esperando 3s... (tentativa %d de 10)", tentativa+1)
		time.Sleep(3 * time.Second)
	}
	if erro != nil {
		return nil, fmt.Errorf("não deu pra conectar no rabbitmq: %w", erro)
	}

	canal, erro := conexao.Channel()
	if erro != nil {
		return nil, fmt.Errorf("falhou na hora de abrir o canal do rabbit: %w", erro)
	}

	_, erro = canal.QueueDeclare(
		"telemetry",
		true,
		false,
		false,
		false,
		nil,
	)
	if erro != nil {
		return nil, fmt.Errorf("nao declarei a fila! %w", erro)
	}

	return &Consumer{conexaoRabbit: conexao, canalRabbit: canal, conexaoBD: bancoConectado}, nil
}

// Start liga e começa a esvaziar a fila
func (c *Consumer) Start() error {
	mensagens, erro := c.canalRabbit.Consume(
		"telemetry",
		"",
		false, 
		false,
		false,
		false,
		nil,
	)
	if erro != nil {
		return fmt.Errorf("não consegui plugar o consumir ali, %w", erro)
	}

	log.Println("Middleware iniciou: Consumindo da filazinha 'telemetry'...")

	// Fica lendo as mensagens que chegam nesse canal do range
	for msg := range mensagens {
		var minhaMsgDaAleta model.TelemetryPayload
		
		// Desempacota o json
		if erro := json.Unmarshal(msg.Body, &minhaMsgDaAleta); erro != nil {
			log.Printf("Veio um json estragado, ignorando: %v", erro)
			msg.Nack(false, false) 
			continue
		}

		// Grava usando o DB que a gente fez
		if erro := c.conexaoBD.InsertReading(minhaMsgDaAleta); erro != nil {
			log.Printf("FALHA AO INSERIR CULPA DO BANCO DE DADOS: %v", erro)
			msg.Nack(false, true) // Deu erro de banco de dados! Re-enfilera essa pra tentarmos depois
			continue
		}

		// Inseriu, manda um sinal dizendo que pode deletar a mensagem
		msg.Ack(false)
		log.Printf("SUCESSO: Gravei do device %s | sensor: %s | valor medido: %v", minhaMsgDaAleta.IdDispositivo, minhaMsgDaAleta.TipoSensor, minhaMsgDaAleta.ValorLido)
	}

	return nil
}

// Fecha tudo se o apocalipse começar
func (c *Consumer) Close() {
	if c.canalRabbit != nil {
		c.canalRabbit.Close()
	}
	if c.conexaoRabbit != nil {
		c.conexaoRabbit.Close()
	}
}