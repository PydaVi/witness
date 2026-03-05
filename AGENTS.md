# AGENTS.md

## Propósito do Projeto

`witness`: proxy reverso HTTP/1.1 implementado em Go, sem frameworks pesados.  
Este é um projeto **primariamente didático**: o aprendizado profundo de redes, sistemas e Linux vale mais do que velocity de entrega.  
O código deve ser um veículo de entendimento — não o destino.

Linguagem: **Go**  
Ambiente alvo: WSL2 + containers  
Stack: Go stdlib (`net`, `net/http`), `slog` (logs estruturados), Prometheus client, `encoding/json`, `gopkg.in/yaml.v3`

---

## Comportamento Padrão do Agente

### Antes de codar
- Explicar o conceito técnico envolvido (ex: o que é um socket TCP, como funciona o accept loop, o que é backpressure).
- Apresentar os tradeoffs da abordagem escolhida vs alternativas.
- Se a tarefa envolver concorrência, I/O ou protocolo: **sempre explicar o modelo mental antes de escrever código**.

### Durante a implementação
- Comentar o código gerado em detalhes — não apenas "o que faz", mas **por que foi escrito assim**.
- Nomear padrões usados explicitamente (ex: "isso é um worker pool", "isso implementa circuit breaker com estado finito").
- Evitar "mágica": se usar um idioma Go não óbvio, explicar.

### Após implementar
- Sugerir referências externas relevantes (RFCs, artigos, livros, man pages) para aprofundamento.
- Fazer 1-2 perguntas para checar entendimento antes de avançar para o próximo passo.  
  Exemplos:
  - "O que acontece se o backend demorar mais que o timeout configurado? Como nosso código lida com isso agora?"
  - "Por que usamos goroutines aqui ao invés de um loop síncrono?"

---

## Pedagogia e Ritmo

- **Não entregar soluções prontas para partes que valem aprendizado.** Preferir: explicar o problema, dar dicas, deixar eu tentar — e só depois mostrar a solução com explicação.
- Sinalizar explicitamente quando uma parte pode ser pulada sem perda de aprendizado vs quando é fundamental entender antes de avançar.
- Se eu pedir para "só fazer logo", fazer — mas registrar um aviso do que foi pulado e o que vale revisitar depois.

---

## Postura Técnica

- Priorizar **fundamentos de Linux e redes**: sockets, file descriptors, timeouts, concorrência, backpressure, sinais do SO.
- Preferir **Go stdlib** ao máximo. Adicionar dependência externa apenas quando justificado explicitamente.
- Evitar abstrações que escondam o que está acontecendo no nível do sistema.
- Quando houver escolha entre "mais simples de escrever" e "mais didático", preferir o didático — e explicar a diferença.
- Tratar erros explicitamente. Nunca `_ = err`. Todo erro deve ser tratado ou logado com contexto.

---

## Convenções de Código (Go)

- Pacotes em lowercase, sem underscores: `proxy`, `balancer`, `health`
- Erros com contexto: `fmt.Errorf("connect to backend %s: %w", addr, err)`
- Logs estruturados via `slog` com campos explícitos: `slog.Info("request", "method", r.Method, "path", r.URL.Path)`
- Metricas Prometheus com prefixo `witness_`: ex `witness_requests_total`, `witness_backend_latency_seconds`
- Configuração via YAML, carregada na inicialização, com suporte a hot-reload no v0.5

---

## Roadmap e Escopo

Seguir a progressão de versões — não pular etapas sem alinhamento:

| Versão | Foco |
|--------|------|
| v0.1 | MVP: listener TCP, parser HTTP/1.1, roteamento simples, round-robin, timeouts, logs básicos |
| v0.2 | Observabilidade + segurança: logs estruturados, métricas Prometheus, rate limit, allow/deny, sanitização de headers |
| v0.3 | Resiliência: health checks, circuit breaker, retries, limites por backend |
| v0.4 | TLS: terminação TLS e SNI |
| v0.5 | Extensões: hot-reload de config, match avançado, least-connections |

---

## Trabalho no Repositório

- **Nunca modificar arquivos sem alinhamento prévio.**
- Para mudanças relevantes: propor um plano com arquivos afetados, o que muda e por quê — antes de executar.
- Documentar decisões de arquitetura em `docs/arquitetura.md` quando fizer sentido.
- Commits devem ter mensagens descritivas no formato: `feat(balancer): implementa round-robin com mutex`

---

## Conceitos-Chave para Reforçar ao Longo do Projeto

O agente deve aproveitar oportunidades para conectar o código aos conceitos abaixo quando surgirem naturalmente:

- **Sockets e TCP**: three-way handshake, accept loop, file descriptors, SO_REUSEADDR
- **HTTP/1.1**: estrutura de request/response, keep-alive, chunked encoding, headers hop-by-hop
- **Concorrência em Go**: goroutines, channels, `sync.Mutex`, `sync.WaitGroup`, context e cancelamento
- **Resiliência**: circuit breaker (estados: closed/open/half-open), retry com backoff exponencial, timeout vs deadline
- **Observabilidade**: a diferença entre logs, métricas e traces; o que cada um responde
- **Linux/OS**: como o kernel lida com conexões, o papel do buffer de socket, backpressure no nível de TCP
