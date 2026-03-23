package main

import (
	"log"
	"os"
	"time"

	"github.com/ponderada/middleware/consumer"
	"github.com/ponderada/middleware/db"
)

func main() {
	// Pega os endereços dos serviços do ambiente
	urlRabbit := os.Getenv("RABBITMQ_URL")
	if urlRabbit == "" {
		urlRabbit = "amqp://guest:guest@localhost:5672/"
	}

	urlBancoDeDados := os.Getenv("DATABASE_URL")
	if urlBancoDeDados == "" {
		urlBancoDeDados = "postgres://postgres:postgres@localhost:5432/telemetry?sslmode=disable"
	}

	// Essa parte serve pra aguardar o Postgres ligar e não crashar tudo
	var bancoDoPostgres *db.Database
	var erro error
	
	for tentativa := 0; tentativa < 10; tentativa++ {
		bancoDoPostgres, erro = db.NewDatabase(urlBancoDeDados)
		if erro == nil {
			break
		}
		log.Printf("banco de dados ainda não ligou, tentando de novo em 3 segundos... (tentativa %d de 10)", tentativa+1)
		time.Sleep(3 * time.Second)
	}
	
	// Se mesmo depois de tentar 10 vezes deu erro, aí não tem jeito
	if erro != nil {
		log.Fatalf("Não consegui logar no banco de maneira alguma: %v", erro)
	}
	// Fecha a conexão na hora que a main() acabar
	defer bancoDoPostgres.Close()

	// Agora vou criar o nosso consumidor passando o RabbitMQ e a conexão com o banco
	consumidor, erro := consumer.NewConsumer(urlRabbit, bancoDoPostgres)
	if erro != nil {
		log.Fatalf("Erro mortal ao criar o Consumer: %v", erro)
	}
	defer consumidor.Close()

	// Dá o start nele e se der erro explode na tela
	if erro := consumidor.Start(); erro != nil {
		log.Fatalf("Erro ao iniciar o consumidor: %v", erro)
	}
}