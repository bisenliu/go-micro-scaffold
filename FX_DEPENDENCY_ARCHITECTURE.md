# Go 微服务架构依赖关系图

## 概述

本图展示了 Go 微服务架构中的组件依赖关系，包括配置、数据库、HTTP 服务、业务逻辑、验证器等模块之间的调用关系。通过此图可以清晰地了解系统各部分的交互方式和依赖结构。

## 依赖关系图

```mermaid
graph TD
    A[main] --> B[go.uber.org/fx]
    B --> C[fx.DotGraph]
    C --> D[*App.dotGraph-fm]
    D --> E[go.uber.org/fx]
    E --> F[fx.Shutdowner]
    F --> G[*App.shutdowner-fm]
    G --> H[go.uber.org/fx]
    H --> I[fx.Lifecycle]
    I --> J[New.func1]

    J --> K[common/config]
    K --> L[*config.Config]
    L --> M[NewConfig]

    M --> N[common/logger]
    N --> O[*zap.Logger]
    O --> P[NewLogger]

    M --> Q[common/databases/mysql]
    Q --> R[*mysql.Manager]
    R --> S[NewManager]

    M --> T[common/databases/redis]
    T --> U[*redis.RedisClient]
    U --> V[NewRedisClient]

    M --> W[common/pkg/idgen]
    W --> X[ulgen.Generator]
    X --> Y[NewGenerator]

    M --> Z[common/http]
    Z --> AA[*gin.Engine]
    AA --> AB[NewServer]

    AB --> AC[common/http]
    AC --> AD[*http.Server]
    AD --> AE[init.func1]

    AE --> AF[common/databases/mysql]
    AF --> AG[mysql.ManagerInterface]
    AG --> AH[init.func1]

    AH --> AI[services/internal/infrastructure/persistence/ent]
    AI --> AJ[*gen.Client]
    AJ --> AK[init.func1]

    AK --> AL[services/internal/infrastructure/messaging]
    AL --> AM[NewRedisEventPublisher]
    AM --> AN[messaging.EventPublisher]

    AN --> AO[services/internal/interfaces/http/handler]
    AO --> AP[*handler.HealthHandler]
    AP --> AQ[NewHealthHandler]

    AQ --> AR[services/internal/application/query/handler]
    AR --> AS[*queryhandler.UserQueryHandler]
    AS --> AT[NewUserQueryHandler]

    AT --> AU[services/internal/domain/user/validator]
    AU --> AV[validator.UserValidator]
    AV --> AW[NewUserValidator]

    AW --> AX[services/internal/domain/user/service]
    AX --> AY[*service.UserDomainService]
    AY --> AZ[NewUserDomainService]

    AZ --> BA[services/internal/application/command/handler]
    BA --> BB[*commandhandler.UserCommandHandler]
    BB --> BC[NewUserCommandHandler]

    BC --> BD[services/internal/interfaces/http/handler]
    BD --> BE[*handler.UserHandler]
    BE --> BF[NewUserHandler]

    BF --> BG[common/pkg/validation]
    BG --> BH[*validation.Validator]
    BH --> BI[NewValidator]

    BI --> BJ[common/pkg/json]
    BJ --> BK[*jwt.JWT]
    BK --> BL[NewJWTService]

    BL --> BM[reflect]
    BM --> BN[mustFuncStub]
    BN --> BO[repository.UserRepository]

    BO --> BP[services/internal/infrastructure/persistence/ent]
    BP --> BQ[*gen.Client]
    BQ --> BR[init.func1]
```



## 

![graph](./assets/graph.png)



## 依赖关系.dto文件

 [dependency-graph.dot](assets/dependency-graph.dot) 