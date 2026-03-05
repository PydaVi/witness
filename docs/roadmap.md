# Roadmap Detalhado

Este documento detalha as entregas planejadas, exemplos de uso e o potencial de aprendizado em cada etapa. O objetivo e guiar a implementacao incremental do witness, mantendo foco em fundamentos de redes e sistemas.

## V0.1 - MVP HTTP/1.1 (Base Funcional)

### Objetivo
Construir o witness como um reverse proxy minimo, funcional e observavel o suficiente para depurar conexoes e comportamento HTTP.

### Entregas
- Listener TCP com backlog configuravel.
- Parsing HTTP/1.1 (linha de request + headers).
- Roteamento simples (host + path prefix).
- Load balancing round-robin.
- Timeouts basicos (read/write/connect).
- Logs basicos com latencia e status.
- Docker Compose com 2 backends HTTP simples.

### Exemplos de uso
- Iniciar o witness apontando para dois backends.
- Enviar requisicoes com `curl` e observar o round-robin.

Exemplo de configuracao (YAML):
```yaml
listener:
  addr: "0.0.0.0:8080"

upstreams:
  - name: app
    targets:
      - "127.0.0.1:9001"
      - "127.0.0.1:9002"

routes:
  - host: "example.local"
    path_prefix: "/"
    upstream: "app"
```

### Potencial de aprendizado
- Como sockets TCP funcionam em Linux.
- Nocoes de backlog, accept loop e read/write com timeouts.
- Fundamentos do protocolo HTTP/1.1.
- Distribuicao de carga basica.

## V0.2 - Observabilidade + Seguranca Basica

### Objetivo
Entender como proxies produzem sinais operacionais e como aplicam protecoes iniciais.

### Entregas
- Logs estruturados (JSON) com request-id.
- Metricas Prometheus: QPS, latencia, status codes.
- Rate limiting por IP (token bucket).
- Allow/deny por CIDR.
- Sanitizacao de headers e limite de body.

### Exemplos de uso
- Configurar rate limit por IP e testar burst com `wrk`.
- Bloquear um IP e confirmar resposta 403.

Exemplo (YAML):
```yaml
filters:
  max_body_bytes: 1048576
  rate_limit:
    requests: 20
    per_seconds: 1
  allow_cidrs:
    - "10.0.0.0/24"
```

### Potencial de aprendizado
- Design de logs e metricas para diagnostico.
- Implementacao de token bucket.
- Validacao de headers e defesa inicial contra abuso.

## V0.3 - Resiliencia

### Objetivo
Aprender como proxies protegem backends e mantem disponibilidade.

### Entregas
- Health checks ativos (HTTP ping periodico).
- Circuit breaker por upstream.
- Retries controlados (apenas para metodos idempotentes).
- Limites de conexoes por backend.

### Exemplos de uso
- Derrubar um backend e observar failover.
- Simular falhas e ver circuito abrir/fechar.

### Potencial de aprendizado
- Teoria e pratica de circuit breaker.
- Estrategias de retry e riscos de tempestade.
- Backpressure e protecao de recursos.

## V0.4 - TLS Termination

### Objetivo
Entender criptografia aplicada em proxies e implicacoes de terminacao TLS.

### Entregas
- TLS offload (cert local).
- Suporte a SNI para roteamento por host.
- Hardening basico (versoes/ciphers).

### Exemplos de uso
- Subir o witness com certificado autoassinado.
- Testar com `curl -k https://...`.

### Potencial de aprendizado
- Handshake TLS e SNI.
- Como proxies lidam com certificados.
- Diferencas entre TLS passthrough e termination.

## V0.5 - Extensoes e Refinos

### Objetivo
Expandir recursos e consolidar o projeto para portfolio.

### Entregas
- Hot-reload de configuracao.
- Match avancado (headers, query params).
- Least-connections.
- Tracing opcional (span por request).

### Exemplos de uso
- Alterar config sem reiniciar o witness.
- Roteamento por header especifico.

### Potencial de aprendizado
- Reconfiguracao segura em runtime.
- Estruturas de dados e sincronizacao.
- Observabilidade avancada.

## V0.6 - Narrative Mode

### Objetivo
Dar voz ao witness por meio de observacoes narrativas sobre o comportamento do trafego.

### Entregas
- Janelas deslizantes de tempo por backend.
- Baseline historico (media e desvio padrao) por janela.
- Deteccao de desvio por z-score.
- Correlacao temporal com anomalias anteriores.
- Formatter narrativo para mensagens humanas.

### Exemplos de uso
- Emitir observacao quando latencia sobe acima do baseline.
- Detectar aumento gradual de 5xx antes do threshold.

### Potencial de aprendizado
- Estatistica basica aplicada a observabilidade.
- Tradeoffs entre precisao e custo de memoria.
- Diferenca entre alerting reativo e observacao proativa.

## Expansoes Opcionais
- HTTP/2 (suporte parcial no witness).
- Cache simples por rota.
- WAF basico com regras customizaveis.
- Suporte a WebSockets.

## Como Explorar Conceitos
- Criar pequenas provas de conceito por feature antes de integrar.
- Registrar observacoes no `docs/` para portfolio.
- Comparar comportamento com proxies reais (nginx/Envoy) para validar entendimento.
