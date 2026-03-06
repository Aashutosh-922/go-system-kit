flowchart LR

    subgraph Clients
        C1[Client / Postman]
        C2[Trading UI]
    end

    subgraph Gateway
        G[API Gateway]
    end

    subgraph Transaction Pipeline
        I[Ingestion Service]
        R[Risk Engine]
        A[AI Insights Service]
        TW[Timescale Writer]
    end

    subgraph Trading Pipeline
        O[Order Service]
        M[Matching Engine]
        PF[Portfolio Service]
        V[Volatility AI]
        MD[Market Data Service]
    end

    subgraph Databases
        P[(Postgres\nOperational DB)]
        TS[(TimescaleDB\nEvent Timeline)]
    end

    subgraph Streaming Backbone
        K[(Kafka Event Bus)]
    end

    C1 --> G
    C2 --> G

    G --> I
    I --> P
    I --> K

    K --> R
    R --> K

    K --> A
    A --> K

    K --> TW
    TW --> TS

    C1 --> O
    O --> K

    K --> M
    M --> K

    K --> PF
    PF --> P

    K --> V
    V --> K

    K --> MD
