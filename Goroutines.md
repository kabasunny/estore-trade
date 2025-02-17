graph LR
    subgraph main.go
        A[main goroutine] --> B(EventStream goroutine)
        A --> C(AutoTradingUsecase goroutine)
        A --> D(HTTP Server goroutine)
    end

     subgraph HTTP Requests
        D -.->|Request| H[TradingHandler.HandleTrade]
    end

    subgraph tradingUsecase
      F[tradingUsecase]
    end

    H --> F

    style B fill:#ccf,stroke:#333,stroke-width:2px
    style C fill:#ccf,stroke:#333,stroke-width:2px
    style D fill:#ccf,stroke:#333,stroke-width:2px
    style H fill:#aaf,stroke:#333,stroke-width:2px
     style F fill:#fcf,stroke:#333,stroke-width:2px