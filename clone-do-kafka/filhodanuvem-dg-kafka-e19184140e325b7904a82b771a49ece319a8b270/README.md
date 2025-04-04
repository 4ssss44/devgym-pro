# Clone do Kafka

![](https://raw.githubusercontent.com/devgymbr/files/main/devgymblack.png). 

Esse é um desafio Devgym, encontre a descrição [aqui](https://app.devgym.com.br/challenges/1ccb06b2-ce93-4450-a17f-9f2479664cff). 

### Instruções

* Rode os testes de integração com `go test -v ./...`. 
* `go run cmd/server/main.go` para rodar o servidor. 
* `go run cmd/cli/main.go -c -n <consumer> -t <topico>` para começar um consumer. 
* `go run cmd/cli/main.go -p -t <topic> -m <message>` para publicar um mensagem num tópico. 
