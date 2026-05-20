# OpenWatchIt — Architecture

## Component Overview

Visão completa dos componentes e do fluxo de dados de um query de logs via Loki.

```mermaid
flowchart TB
    classDef cli     fill:#1e3a5f,stroke:#4a90d9,color:#e8f4fd,rx:6
    classDef engine  fill:#1a3a2a,stroke:#4caf50,color:#e8f5e9,rx:6
    classDef plugin  fill:#3a1a2a,stroke:#e91e63,color:#fce4ec,rx:6
    classDef ext     fill:#2d2d2d,stroke:#ff9800,color:#fff3e0,rx:6
    classDef proto   fill:#2a1a3a,stroke:#9c27b0,color:#f3e5f5,rx:6
    classDef user    fill:#0d0d0d,stroke:#888,color:#fff,rx:20

    USER(["User — Terminal"]):::user

    subgraph CLI ["CLI  ·  owit binary"]
        direction LR
        OWL["OWL Parser\n─────────────\nflags + args → AST"]:::cli
        REN["Renderer\n─────────────\ntable · json · stream"]:::cli
    end

    subgraph CP ["Control Panel  ·  Engine Core"]
        direction TB
        QP["Query Planner\n─────────────────────────────\nresolver de fan-out scope\nvalidação de signal type\ndetecção de cross-signal join"]:::engine
        PM["Plugin Manager\n─────────────────────────────\ndescobre plugins instalados\nspawn de processos filhos\nmantem conexões gRPC\nhealth check · timeout"]:::engine
        DI["Dispatcher\n─────────────────────────────\nfan-out paralelo\nnormalização de resultados\nbackends lentos não bloqueiam"]:::engine
        RM["Result Merger\n─────────────────────────────\nconsolida streams de N backends\nordena por timestamp"]:::engine
    end

    subgraph PL ["Plugin Process  ·  gRPC Server"]
        direction TB
        CAP["Capabilities\n──────────────────\ndeclara signal types\ne versão suportados"]:::plugin
        EXEC["Execute  AST  →  stream\n──────────────────────────────\n1. walk AST\n2. emite query nativa  ex: LogQL\n3. chama vendor API\n4. normaliza → NormalizedResultRow"]:::plugin
    end

    subgraph EXT ["Externos"]
        direction TB
        LOKI[("Loki HTTP API\n─────────────\n/loki/api/v1/\nquery_range")]:::ext
        TOML[["TOML Config\n─────────────\n~/.owit/config.toml\nbackends · credentials"]]:::ext
        DIR[["Plugin Dir\n─────────────\n~/.owit/plugins/\n.owit/plugins/"]]:::ext
    end

    PROTO["proto  ·  Contrato gRPC\n────────────────────\nAST Schema\nNormalizedResult\nPlugin RPC interface"]:::proto

    %% main data flow
    USER     -->|"owit logs level=error service=checkout --last 30m --limit 100"| OWL
    OWL      -->|"AST { signal:LOG, filters:[...], last:30m, limit:100 }"| QP
    TOML     -->|"backends + credentials"| QP
    QP       -->|"plano de execução"| PM
    QP       -->|"AST + backends alvo"| DI
    PM      -.->|"discover + spawn"| DIR
    PM       -->|"gRPC conn pronta"| DI
    DI       -->|"Execute(AST, config)  gRPC streaming"| EXEC
    EXEC     -->|"LogQL query"| LOKI
    LOKI     -->|"log entries  JSON"| EXEC
    EXEC    -->|"NormalizedResultRow  stream"| DI
    DI       -->|"rows stream"| RM
    RM       -->|"rows consolidados + sorted"| REN
    REN      -->|"output formatado"| USER
    CAP     --->|"signal types suportados"| PM

    %% contract
    PROTO   -. "define AST" .-> OWL
    PROTO   -. "define RPC + NormalizedResult" .-> EXEC
```

---

## Modos de Execução

### Embedded Mode — uso solo, sem infraestrutura

```mermaid
flowchart LR
    classDef cli    fill:#1e3a5f,stroke:#4a90d9,color:#e8f4fd
    classDef engine fill:#1a3a2a,stroke:#4caf50,color:#e8f5e9
    classDef plugin fill:#3a1a2a,stroke:#e91e63,color:#fce4ec

    U(["User"])
    CLI["owit binary\n────────────────\nOWL Parser\nControl Panel  in-process\nRenderer"]:::cli
    PL["Plugin Process\n(filho do CLI)"]:::plugin
    API[("Vendor API\nex: Loki")]

    U -->|"owit logs ..."| CLI
    CLI -->|"gRPC"| PL
    PL -->|"HTTP"| API
    API -->|"entries"| PL
    PL -->|"NormalizedResultRow"| CLI
    CLI -->|"table"| U

    note1["Sem servidor.\nPlugins sobem lazy,\nmorrem após 30s idle.\nCredenciais no TOML local."]
```

### Remote Mode — uso em time com config centralizada

```mermaid
flowchart LR
    classDef cli    fill:#1e3a5f,stroke:#4a90d9,color:#e8f4fd
    classDef server fill:#1a3a2a,stroke:#4caf50,color:#e8f5e9
    classDef plugin fill:#3a1a2a,stroke:#e91e63,color:#fce4ec

    U(["User"])
    CLI2["owit binary\n──────────────\nclient thin\nOWL Parser\nRenderer"]:::cli
    SRV["owit server\n──────────────────\nControl Panel  remoto\nPlugin Manager\nDispatcher\nResult Merger\nHTTP / gRPC API"]:::server
    PL2["Plugin Process\n(filho do server)"]:::plugin
    API2[("Vendor API")]

    U -->|"owit logs ..."| CLI2
    CLI2 -->|"OWL query  HTTP/gRPC"| SRV
    SRV -->|"gRPC"| PL2
    PL2 -->|"HTTP"| API2
    API2 -->|"entries"| PL2
    PL2 -->|"NormalizedResultRow"| SRV
    SRV -->|"resultado"| CLI2
    CLI2 -->|"table"| U

    note2["Plugins são processos\nlong-running do servidor.\nCredenciais centralizadas.\nTimes compartilham backends."]
```

---

## Contrato gRPC do Plugin

Definido em `.proto` — a fronteira estável entre o engine e qualquer plugin.

```mermaid
flowchart LR
    classDef engine fill:#1a3a2a,stroke:#4caf50,color:#e8f5e9
    classDef plugin fill:#3a1a2a,stroke:#e91e63,color:#fce4ec
    classDef msg    fill:#1a1a2e,stroke:#7986cb,color:#e8eaf6

    ENG["Engine\nPlugin Manager + Dispatcher"]:::engine

    subgraph RPC ["Plugin gRPC Service"]
        C["Capabilities\n─────────────────────\nreq:  empty\nres:  signal_types, version"]:::plugin
        E["Execute\n─────────────────────\nreq:  ExecuteRequest\n      ast + backend_config\nres:  stream NormalizedResultRow"]:::plugin
    end

    subgraph MSG ["Mensagens proto"]
        AST["AST\n──────────────────\nsignal_type: SignalType\nfilters:     Filter list\ntime_range:  Duration\nlimit:       int32"]:::msg
        NRR["NormalizedResultRow\n──────────────────\ntimestamp  int64\nlevel      string\nservice    string\nmessage    string\n_source    string\npassthrough map string string"]:::msg
    end

    ENG -->|"Capabilities()"| C
    C   -->|"CapabilitiesResponse"| ENG
    ENG -->|"Execute(ExecuteRequest)"| E
    E   -->|"stream NormalizedResultRow"| ENG

    AST -. "compõe ExecuteRequest" .-> E
    E   -. "emite" .-> NRR
```

---

## Ordem de Implementação — Phase B First

```mermaid
gantt
    title Roadmap de Implementação v0.1
    dateFormat  X
    axisFormat  Passo %s

    section Contrato
    .proto  AST + Plugin RPC + NormalizedResult  :crit, done, p1, 0, 1

    section Plugin Loki  isolado
    gRPC server stub                             :crit, p2, 1, 2
    Tradução AST → LogQL                         :crit, p3, 2, 3
    Chamada Loki HTTP API                        :crit, p4, 3, 4
    Testes de integração  AST hardcoded          :crit, p5, 4, 5

    section Engine Core
    OWL Parser  flags → AST                      :p6, 5, 6
    Plugin Manager  spawn + gRPC lifecycle       :p7, 6, 7
    Dispatcher + Result Merger                   :p8, 7, 8

    section CLI
    owit binary  Embedded Mode                   :p9, 8, 9
    Renderer  table + json                       :p10, 9, 10
```
