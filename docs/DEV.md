# Architecture

* Reverse proxy server (Nginx)
* API server (Golang/Echo)
* App server (TypeScript/Remix)
* Database (PostgreSQL)
* Worker (Golang/Echo + Swift/SwiftWasm + WebAssembly/Wasmtime)

# Dependencies

* Docker
* Docker Compose
* Node.js 20.0.0 or later
* Npm
* Go 1.22.3 or later

# Run

1. Clone the repository.
1. `cd path/to/the/repo`
1. `make init`
1. `make up`
1. Access to http://localhost:5173.
    * User `a`, `b` and `c` can log in with `pass` password.
    * User `a` and `b` are players.
    * User `c` is an administrator.
