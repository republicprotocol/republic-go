# Architecture

The official Go implementation of Republic Protocol uses two core design principles.

**[Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)** encourages a layered design and a uni-direction flow of dependencies that is decoupled from the flow of control.

**[Interface Driven Development](https://en.wikipedia.org/wiki/Interface-based_programming)** encourages components to depend on abstract interfaces, instead of concrete implementations.

## Layers and Dependencies

The official Go implementation of Republic Protocol is composed of four layers. Dependencies can only flow in one direction — from the most concrete layers, to the most abstract layers — without exception.

1. **Driver**
   The **driver** layer is the most concrete layer, defining infrastructure and external interfaces (e.g. the command-line, http, grpc, and leveldb). Packages and components in the driver layer can depend on any other layer; depending on the adapter layer for data formats, and the service layer for interfaces. Drivers typically implement interfaces defined in the service layer since this allows services to hand control to the drivers without introducing a dependency from the service layer to the driver layer.
2. **Adapter**
   The **adapter** layer defines adapters that consume these data formats used by the drivers, and produce data formats used by services. Data formats defined in this layer are used directly by drivers. Packages and components in the adapter layer can depend on the service and domain layers, but are forbidden from depending on the driver layer.
3. **Service**
   The **service** layer defines different applications that implement the core business logic of Republic Protocol (e.g. the secure order matching engine). Packages and components in the service layer can depend on the domain layer, but are forbidden from depending on the driver and adapter layers.
4. **Domain**
   The **domain** layer is the most abstract layer, defining components that are used heavily throughout Republic Protocol and are intrinsic to its domain — the Darkpool (e.g. orders, order fragments, addresses, and cryptography). This layer is the most abstract because without services, no core functionality is expressed.

### Drivers

#### Smart Contracts

#### gRPC

#### LevelDB

#### HTTP

### Services

#### Darknodes

#### Secure Multiparty Computations (sMPC)

#### Swarm

#### Brokers

### Domains

#### Orders

#### Addresses

#### Keystores