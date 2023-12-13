# Order Book

WIP

## Folder structure

```
Transport
  |
  +---> Service
          |
          +---> Data
```

### Transport

Handles the communication with the outside world (like HTTP requests, gRPC, message queues) and translates those requests into actions or queries in the service layer.

### Service

Contains business logic, orchestrates data flow between the data layer and the transport layer, and makes decisions based on business rules.

### Data

Responsible for data persistence, retrieval, and direct interactions with the data storage mechanisms (databases, file systems, external APIs, including blockchain nodes).
