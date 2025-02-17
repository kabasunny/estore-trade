graph LR
    subgraph Tachibana API
        E[EVENT I/F]
        G[Tachibana API]
    end

    subgraph tradingUsecase
        F[tradingUsecase]
    end

    B(EventStream goroutine) -- "eventCh" --> F
    %% OrderEvent sent via eventCh
    C(AutoTradingUsecase goroutine) -- "eventCh" --> F
    %% OrderEvent sent via eventCh
    E -.->|Event Data| B
    F -->|TachibanaClient| G


    style B fill:#ccf,stroke:#333,stroke-width:2px
    style C fill:#ccf,stroke:#333,stroke-width:2px
    style F fill:#fcf,stroke:#333,stroke-width:2px