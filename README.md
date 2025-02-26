# go-expert-rate-limiter

Implementação de um servidor simples HTTP com apenas um handler para demonstração de desenvolvimento de um rate limiter utilizando Redis.

É necessário ter `make`, `docker`, `compose`, `apache benchmark` e `curl` instalados para executar os passos abaixo.

Para executar a aplicação:

```bash
make run
```

A configuração do rate limiter se encontra no arquivo [./cmd/server/.env](./cmd/server/.env)

- MAX_REQUESTS_PER_SECOND_PER_IP: define a quantidade de requests por segundo por IP.
- MAX_REQUESTS_PER_SECOND_PER_API_TOKEN: define a quantidade de requests por segundo por API KEY.
- BAN_DURATION: define a quantidade de tempo que o IP ou a API KEY ficarão banidos caso atinjam os limites acima.

Para testes utilize os seguintes comandos:

Teste IP:

```bash
ab -n 1000 -H "X-Real-IP: 127.0.0.1" http://localhost:8080/
```

Teste API KEY:

```bash
ab -n 1000 -H "API_KEY: ABCD" http://localhost:8080/
```

Teste API KEY com prioridade ao IP:

```bash
ab -n 1000 -H "X-Real-IP: 127.0.0.1" -H "API_KEY: ABCD" http://localhost:8080/
```

Para visualizar as mensagens de erro, utilize os comandos abaixo em conjunto com os comandos acima.

```bash
curl --header "X-Real-IP: 127.0.0.1" http://localhost:8080/
```

```bash
curl --header "API_KEY: ABCD" http://localhost:8080/
```

```bash
curl --header "X-Real-IP: 127.0.0.1" --header "API_KEY: ABCD" http://localhost:8080/
```

Personalizei a mensagem de erro para mostrar qual item foi banido(IP ou API KEY) e também até quando ficará banida, exemplo:

```log
The API KEY ABCD has reached the maximum number of requests or actions allowed within a certain time frame, wait until 2025-02-26T22:21:12Z
```
