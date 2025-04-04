# Um programa para popular um banco de dados com milhares de filmes

![](https://raw.githubusercontent.com/devgymbr/files/main/devgymblack.png). 

Esse é um desafio Devgym, encontre a descrição [aqui](https://app.devgym.com.br/challenges/ec36e7e2-6a2d-4406-98e1-3029f843b5c3). 

## Instruções

* rode `source env.sh` para exportar variávies de ambientes importantes.
* rode `docker-compose up --force-recreate --build` para subir o banco de dados.
* rode `go run ./cmd/main.go` para executar o programa com todos os valores default. Adicione `--help` para ver todas as opcões do programa.
* você pode verificar o banco de dados numa interface gráfica acessando `http://localhost:8079`.

## Rodar tests

```bash
docker-compose up test
```