# Arquitetura do Witness (Didatico)

## Visao Geral
O witness implementa um reverse proxy HTTP/1.1 para estudo profundo de networking, seguranca e comportamento de sistemas no Linux. O foco e construir os conceitos na unha: parsing de requests, roteamento, balanceamento, filtros, controle de trafego e resiliencia.

## Requisitos de Design
- Baixo nivel: sockets, timeouts, concorrencia, backpressure.
- Transparencia: cada decisao deve ser explicada e documentada.
- Didatico: codigo e docs devem facilitar o aprendizado e o portfolio.
- Sem frameworks pesados ou abstracoes altas.

## Stack Tecnologica
- Linguagem: Go (stdlib)
- Rede/HTTP: `net`, `net/http`
- Logs: `log/slog`
- Config: YAML (`gopkg.in/yaml.v3`)
- Metricas: Prometheus client
- Ambiente: WSL2 + Docker Compose

## Componentes Principais
- `core/`: ciclo de vida do witness e orquestracao do accept loop.
- `netx/`: wrappers didaticos para `net.Listener`, timeouts e limites de conexoes.
- `http1/`: parsing HTTP/1.1, validacao minima e normalizacao.
- `routing/`: match de host/path/method.
- `balancer/`: algoritmos de balanceamento (round-robin, least-connections).
- `filters/`: allow/deny, rate limit, max body, header sanitation.
- `resilience/`: timeouts, retries, circuit breaker, backpressure.
- `observability/`: logs estruturados e metricas.
- `config/`: parsing YAML, validacao e defaults.
- `admin/`: health/metrics/status do witness.

## Fluxo Interno do Witness
1. Aceita conexao TCP.
2. Aplica limites globais de conexoes e timeouts de handshake.
3. Faz parse do request HTTP/1.1.
4. Valida request e aplica filtros de seguranca.
5. Resolve rota e seleciona upstream (balancer).
6. Conecta ao upstream e faz proxy de request/response.
7. Aplica politicas de resiliencia.
8. Emite logs e metricas.

## Configuracao (YAML)
- `listeners`: endereco e porta de entrada.
- `routes`: regras por host/path/method.
- `upstreams`: pools de backends com pesos e timeouts.
- `filters`: rate limit, allow/deny, header rules.
- `observability`: log level e metrics endpoint.

## Roadmap
### v0.1 MVP
- Listener TCP
- Parser HTTP/1.1
- Roteamento simples
- Round-robin
- Timeouts basicos
- Logs simples

### v0.2 Observabilidade + Seguranca
- Logs estruturados (request-id)
- Metricas Prometheus
- Rate limiting por IP
- Allow/deny
- Sanitizacao de headers

### v0.3 Resiliencia
- Health checks ativos
- Circuit breaker
- Retries controlados
- Limites por backend

### v0.4 TLS
- Terminacao TLS
- Suporte a SNI
- Hardening basico

### v0.5 Extensoes
- Hot-reload de config
- Match avancado
- Least-connections
- Tracing opcional

### v0.6 Narrative Mode
- Janelas deslizantes de tempo
- Deteccao de desvio (z-score)
- Baseline por backend
- Correlacao temporal
- Formatter narrativo

## Critérios de Sucesso
- Entendimento claro de cada componente e tradeoff.
- Witness funcional em ambientes locais (WSL2 + containers).
- Documentacao suficiente para portfolio tecnico.
