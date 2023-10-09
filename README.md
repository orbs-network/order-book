# Order Book

WIP

## Folder structure

1. REST API is exposed at the `transport` layer
2. `transport` translates it into a known type, handles authentication, calls `service`
3. `service` is where main business logic takes place, calls `data` layer to fetch and persist data
4. `data` layer has specific implementations for interacting with blockchain / database / external services
