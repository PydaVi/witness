# witness

`witness` e um reverse proxy HTTP/1.1 escrito em Go, do zero, sem frameworks.

Projeto de aprendizado com objetivo claro: entender como proxies funcionam **por dentro** — sockets, accept loops, parsing de protocolo, concorrência, resiliência — no nível que a maioria dos engenheiros de cloud nunca precisou se preocupar.

Não é um produto. É um laboratório. Mas tem uma ideia própria.

---

## Por que isso existe

Trabalho com infraestrutura e segurança em cloud há alguns anos. Opero firewalls, configuro ingress controllers, leio logs de proxy reverso todo dia. Em algum momento percebi que sabia *usar* essas ferramentas mas não sabia o que acontecia quando uma conexão chegava — o que o kernel fazia, por que um timeout se comportava diferente de outro, o que "backpressure" significa de verdade fora de uma slide de apresentação.

Esse projeto é a resposta para isso.

Mas tinha outra coisa que me incomodava: proxies e firewalls tomam decisões o tempo todo e ficam em silêncio sobre elas. Quando algo quebra, você vai caçar nos logs — e os logs estão cheios de ruído ou vazios de contexto. A ferramenta sabia o que estava acontecendo. Não te contou.

Daí veio a ideia principal desse projeto.

---

## Narrative Mode

A feature que diferencia o witness de um exercicio tecnico generico.

Em vez de so registrar fatos em logs e metricas, o witness emite **observacoes em linguagem proxima do humano** sobre o que esta vendo — antes que vire problema.

```
[narrative] backend-api está respondendo 40ms acima da média das últimas 2h.
            Ainda dentro do threshold, mas tendência é de piora desde 14:30.

[narrative] Taxa de erros 5xx em /api/checkout subiu de 0.2% para 1.8% nos últimos 5min.
            Padrão similar ocorreu antes da última janela de degradação (2025-03-01 02:14).

[narrative] Cliente 192.168.1.45 fez 312 requisições nos últimos 60s.
            Comportamento fora do padrão histórico desse IP. Rate limit não atingido ainda.
```

Nao e alerting tradicional (threshold -> alarme). E o witness comparando comportamento atual com historico, detectando desvios com contexto, e te contando o que esta notando — com voz.

Detalhes de design e implementação: [`docs/narrative-mode.md`](docs/narrative-mode.md)

---

## O que está sendo construído

**v0.1 — MVP**
Listener TCP, parsing HTTP/1.1 manual, roteamento basico, balanceamento round-robin, timeouts configuraveis e logs basicos. O objetivo aqui e ter o witness funcionando e entender cada linha do porque.

**v0.2 — Observabilidade + Segurança**
Logs estruturados com `slog`, métricas Prometheus, rate limiting, allow/deny list e sanitização de headers. A pergunta que guia essa fase: o que eu precisaria ver num incidente de segurança?

**v0.3 — Resiliência**
Health checks ativos, circuit breaker com máquina de estados (closed → open → half-open), retries com backoff e limites por backend. Sistemas distribuídos falham — essa fase é sobre falhar de forma controlada.

**v0.4 — TLS**
Terminação TLS e SNI. Entender o handshake no nível do código, não só configurar um certificado.

**v0.5 — Extensões**
Hot-reload de configuração, least-connections, roteamento avançado por path/header.

**v0.6 — Narrative Mode**
O witness ganha voz. Janelas de tempo deslizantes, deteccao de desvio por z-score, correlacao temporal e formatter narrativo. Sem ML, sem dependencias pesadas — estatistica basica sobre uma base de metricas solida.

CLI natural em ingles, mantendo a proposta de "testemunha do trafego":

```bash
witness start
witness logs
```

---

## Stack

Go puro, deliberadamente sem frameworks pesados.

- `net` — sockets e conexões TCP
- `net/http` — onde faz sentido não reinventar a roda
- `log/slog` — logs estruturados
- `gopkg.in/yaml.v3` — configuração
- Prometheus client — métricas

Ambiente de desenvolvimento: WSL2 + containers.

---

## Status

> Em desenvolvimento ativo. Atualmente na v0.1.

Acompanhe o progresso pelas branches e pelo histórico de commits — cada versão tem decisões documentadas de arquitetura em `docs/arquitetura.md`.

---

## Documentação

- [`docs/arquitetura.md`](docs/arquitetura.md) — decisões de design e tradeoffs
- [`docs/narrative-mode.md`](docs/narrative-mode.md) — a ideia central do projeto

---

## Conceitos explorados

À medida que o projeto avança, cada parte do código conecta a um conceito de sistemas ou redes:

- Como o kernel gerencia conexões TCP e o que acontece no `accept()`
- A estrutura real de um request HTTP/1.1 e por que parsear isso manualmente é instrutivo
- Goroutines, `sync.Mutex` e o modelo de concorrência do Go na prática
- O que é backpressure e quando um proxy deve recusar conexões
- Circuit breaker como máquina de estados — não como buzzword
- Detecção de anomalias sem ML: janelas deslizantes e z-score na prática

---

*Feito por alguém que prefere entender antes de operar — e que cansou de ferramentas que sabem mais do que dizem.*
