# OrderbookSimulator

Real-time order matching simulator with a price-time priority matching engine, message-broker architecture, and a trading frontend.

## Overview
Users place limit or market orders to the orderbook via the frontend. That request then gets routed through an HTTP handler to the matching engine which matches them by price-time priority. The stream of events emitted from the engine and are then written to a Postgres database and routed live to the browser over websockets.

## Architecture

```mermaid
graph TD;
    A[HTTP handler] -->|orders exchange| B(orders exchange);
    B -->|engine queue| C{matching engine};
    C --> X(events exchange);
    X -->|writer queue| D[db-writer];
    X -->|ws queue| E[ws];
    E --> F((browser));
    D --> G[(Postgres)];
```
The HTTP handler wraps the request into an envelope which carries information about the order type (cancel or a place), the order itself, and the highest sequence number of envelopes placed. It then sends this payload into the orders exchange (direct) that sends the payload to an engine queue that the reads from. The matching engine subscribes to the engine queue after the orderbook is initialized, and then reads from the stream of envelopes and emits events based on how the order affected the orderbook. These events are fill, cancel, and rest. A separate publisher goroutine reads from a channel of these events and publishes them to the events exchange (fanout). From here two separate queues are streamed the events from the events exchange, one for the websocket and one for the database writer. A database writer goroutine subscribes to a writer queue and writes these events to the postgres database. A websocket goroutine subscribes to the ws queue which then streams the events to the frontend.

## System Design Decisions

### Engine
The engine is a single goroutine that is lock-free and deterministic. All concurrency is handled through channels and queues without requiring mutexes.
#### Engine implementation

Buy orders and sell orders are each represented in the order book in an ordered tree map that is keyed by price. Each node in the tree is a price node, which represents the price of the order. The price nodes contain head and tail pointers to order nodes (that have the same price) to create a linked list structure. Order nodes contain the order information and next and prev pointers to other order nodes. Additionally a hashmap that is keyed on a order id has a value of a ordernode, so that cancelling an order is a O(1) operation.

When an order comes in, it looks through the corresponding tree to try and find a match. If it finds a match, the head pointer of the price node gets priority (since it was the first order added at that price). If it doesn't find a match and is a limit order, it then rests in the tree of its type.

### Message-broker Architecture
Components communicate through RabbitMQ using the AMQP protocol. This was designed so that the matching engine, the database writer, and the websocket could be implemented as separate processes.

### [Idempotency](https://en.wikipedia.org/wiki/Idempotence)

Both the engine and the database writer processes implement idempotency. The engine checks the sequence number of the envelope, and if it's less than or equal to its highest sequence number, then it knows it's already processed that request and it is skipped. The same procedure takes place for the database writer, it gets the highest sequence number from the database and compares it to the sequence number from the event. If the event sequence number is less than or equal to the highest sequence number it returns an ack to signify that it's already processed that event.

### [Atomicity](https://en.wikipedia.org/wiki/Atomicity_(database_systems))
The database writer uses a transaction to write to the database. If any error occurs with writing to the database it sends a nack back to RabbitMQ to requeue the message. If there are no errors with the transaction the transaction is committed and the db writer then sends an ack.

### Replay ([Durability](https://en.wikipedia.org/wiki/Durability_(database_systems)))
The engine rebuilds the book from the database on startup, so that in the event that the engine ever crashes all open or partially filled orders still persist. Additionally the engine is fed the highest envelope sequence number from the database which is used to compare with incoming envelopes.

### Snapshot
Whenever a client on the frontend connects to the websocket they are sent a snapshot of the orderbook. The snapshot represents the price levels of the orderbook. Whenever a change occurs in the orderbook they are also sent the snapshot.

### Authentication
A user must create an account with their email and password. Once they log in they are given a JWT token that is used to authenticate the user prior to they viewing the orderbook or placing orders.

## Tech Stack

- Golang 
- PostgreSQL
- React
- RabbitMQ

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (1.22+)
- [Node.js](https://nodejs.org/) and npm
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [goose](https://github.com/pressly/goose) for database migrations

### Setup

1. **Clone the repository:**
```bash
git clone https://github.com/abolcerek/Trading-Matching-Simulator
cd Trading-Matching-Simulator
```

2. **Start RabbitMQ**:
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management
```
   The management UI is available at `http://localhost:15672`, the default login is guest for both username and password.

3. **Create the database:**
```bash
createdb exchange
```

4. **Run the migrations**:
```bash
cd sql/schema
goose postgres "postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:5432/exchange?sslmode=disable" up
cd ../..
```
   Replace the connection string with your local Postgres credentials. Connection string is different for mac and linux users; for mac it is `postgres://YOUR_USERNAME@localhost:5432/exchange?sslmode=disable`. To test the connection string:
```bash
psql "<connection_string>"
```

5. **Configure environment variables.** Create a `.env` file in the project root:
```dotenv
DB_URL="postgres://postgres:YOUR_PASSWORD@localhost:5432/exchange?sslmode=disable"
JWT_SECRET="YOUR_GENERATED_SECRET_TOKEN"
PORT="8080"
PLATFORM="dev"
```
   To generate a JWT token use:
```bash
openssl rand -base64 64
```

6. **Run the backend** (requires RabbitMQ and Postgres running):
```bash
go run ./cmd/api-gateway
```
   The API serves on `http://localhost:8080`.

7. **Run the frontend:**
```bash
cd frontend
npm install
npm run dev
```
   Open `http://localhost:5173` in your browser.


### Usage

Sign up for an account, log in, and you'll land on the trading view. Place buy or sell orders (limit or market) and watch the order book update live as orders match.

## Notes
This project coincided alongside my reading of [Designing Data-Intensive Applications](https://www.oreilly.com/library/view/designing-data-intensive-applications/9781098119058/) by Martin Kleppmann and Chris Riccomini which heavily inspired the bulk of design choices.

## License
Distributed under the MIT License. See `LICENSE` for more information.