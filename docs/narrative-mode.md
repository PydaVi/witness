# Narrative Mode

## O problema

Proxies e firewalls tomam decisões o tempo todo e ficam em silêncio sobre elas.

Quando algo quebra, você abre os logs — e encontra ruído, ou contexto insuficiente para entender o que levou àquele momento. A ferramenta sabia o que estava acontecendo. Não te contou.

O alerting tradicional não resolve isso. Threshold atingido → alarme é reativo por definição. Você já perdeu.

## A ideia

Um modo de operacao onde o witness **narra o proprio comportamento em tempo real** — emitindo observacoes sobre o que esta notando antes que vire problema.

Nao e um dashboard. Nao e um alarme. E o witness tendo opiniao sobre o proprio trafego e te contando, com contexto suficiente para voce entender por que aquilo importa.

Exemplos do que Narrative Mode emite:

```
[narrative] backend-api está respondendo 40ms acima da média das últimas 2h.
            Ainda dentro do threshold, mas tendência é de piora desde 14:30.
            Últimas 3 anomalias nesse backend: hoje 14:30, ontem 09:12, sex 23:47.

[narrative] Taxa de erros 5xx em /api/checkout subiu de 0.2% para 1.8% nos últimos 5min.
            Padrão similar ocorreu antes da última janela de degradação (2025-03-01 02:14).

[narrative] Cliente 192.168.1.45 fez 312 requisições nos últimos 60s.
            Comportamento fora do padrão histórico desse IP. Rate limit não atingido ainda.
```

## O que diferencia do que já existe

Nginx, Envoy, Caddy têm métricas e logs. Registram fatos.

Narrative Mode tem **voz**. Compara comportamento atual com histórico, detecta desvios antes dos thresholds, e entrega contexto — não só números.

## Implementação (sem dependências pesadas)

A base é estatística simples sobre janelas de tempo deslizantes:

- **Baseline por backend**: média e desvio padrão de latência e error rate por janela (1h, 24h, 7d)
- **Detecção de desvio**: comparação do comportamento atual vs baseline — sem ML, z-score básico resolve
- **Correlacao temporal**: o witness lembra quando viu padroes similares e menciona
- **Formatter narrativo**: transforma métricas em texto legível com contexto

Tudo em memória no MVP. Persistência opcional em v0.6+.

## Posicionamento no roadmap

Narrative Mode é a **v0.6** — depois que observabilidade (v0.2) e resiliência (v0.3) estiverem sólidas.

Voce precisa de metricas de qualidade e historico de comportamento antes de detectar desvios com significado. A v0.2 constroi a fundacao. A v0.6 usa essa fundacao para dar ao witness uma voz.

## Por que isso existe

Trabalho com firewalls e proxies todo dia. A pergunta que me incomodava era simples: a ferramenta sabe o que está acontecendo — por que só me conta quando já é tarde?

Narrative Mode é a resposta.
