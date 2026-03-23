package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ guarda as conexões que o nosso backend precisa
type RabbitMQ struct {
	conexao *amqp.Connection
	canal   *amqp.Channel
}

// Essa função serve pra gente inicializar a conexão no Rabbit
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	var conexao *amqp.Connection
	var erro error
	
	// Fiz um laço de repetição pra tentar conectar 10 vezes com atraso.
	for tentativa := 0; tentativa < 10; tentativa++ {
		conexao, erro = amqp.Dial(url)
		if erro == nil {
			break 
		}
		log.Printf("O RabbitMQ ainda não acordou, tentando de novo em 3 segundos... (tentativa %d de 10)", tentativa+1)
		time.Sleep(3 * time.Second)
	}
	
	if erro != nil {
		return nil, fmt.Errorf("deu ruim de vez para conectar no Rabbit, desisto: %w", erro)
	}

	// Abre um canal de comunicação 
	canal, erro := conexao.Channel()
	if erro != nil {
		return nil, fmt.Errorf("falha ao abrir o canal de comunicação: %w", erro)
	}

	// Aqui a gente cria a fila "telemetry" caso ela ainda não exista
	_, erro = canal.QueueDeclare(
		"telemetry", // nome da nossa fila	
		true,        // vai ser durável para não sumir no reset
		false,       // auto-delete falso
		false,       // exclusive falso
		false,       // no-wait falso
		nil,         // sem argumentos extras
	)
	if erro != nil {
		return nil, fmt.Errorf("não consegui criar ou achar a fila: %w", erro)
	}

	// Retorna o objeto pronto pra uso
	return &RabbitMQ{conexao: conexao, canal: canal}, nil
}

// Publish pega nosso dado genérico, transforma num JSON e arremessa na fila
func (r *RabbitMQ) Publish(dadoGenerico any) error {
	corpoDaMensagem, erro := json.Marshal(dadoGenerico)
	if erro != nil {
		return fmt.Errorf("deu erro na hora de transformar a struct pra JSON: %w", erro)
	}

	// Define um timeout de 5 segundinhos só por precaução
	ctx, cancela := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancela()

	return r.canal.PublishWithContext(ctx,
		"",          // vazio porque a gente tá mandando direto pra fila
		"telemetry", // routing key é simplesmente o nome da nossa fila
		false,       // mandatory falso
		false,       // immediate falso
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         corpoDaMensagem,
			DeliveryMode: amqp.Persistent, 
		},
	)
}

// Close é importante para fechar as conexões
func (r *RabbitMQ) Close() {
	if r.canal != nil {
		r.canal.Close()
	}
	if r.conexao != nil {
		r.conexao.Close()
	}
}