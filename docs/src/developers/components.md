# Components

## Observability

- Logging: I went with `logrus` for its ease of use, active contributor base, and extensive documentation.
- Metrics: I am defaulting to `runtime/metrics` for now
- Tracing: I am defaulting to `runtime/trace` for now

## Authentication

- AuthManager: a process that manages Leases & Storage of AuthTokens
    - Lease: An abstraction that allows extending what can be authenticated against
        - AuthFactory: The interface used to abstract how Credentials and Tokens are retrieved and managed
        - Credentials: The things used to request access
        - Tokens: The things used by the app to access a resource
    - Store: An abstraction that allows extending where/how tokens and Credentials can be stored

## User Interfaces

- Command Line Interface (CLI): This is a binary that can be used for simple semi-interactive sessions from a terminal
- Text User Interface (TUI): This is a fully interactive session that allows deeper and simpler visualization of the gosling data
- Application Programming Interface (API): This is intended to be a fully programmatic and machine focused interface, it can be fronted by some CSR UI Application or interacted with by scripts or other remote tools (including a CLI/TUI set to control a remote app).
- Headless: This is a binary that does not accept any form of interactive input, it's logic has been hardcoded and will only be looking for predefined env vars or configs it needs to authenticate to its endpoints. It will then request, process, and ship the data to where ever it has been configured too.
