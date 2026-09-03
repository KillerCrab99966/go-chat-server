# Go Chat Server

## Idea

A real-time communication server using WebSockets. Clients can connect via a web interface (hosted by the same server) or a dedicated terminal client.

- Support for multiple rooms
    - Room creation can be disabled and a fixed number set
- No accounts
- Code-based authentication
    - Similar to `croc`
    - E.g. `jolt-wreck-scorn`
- Configurable with TOML like usual

## Implementation Plan

1. ~~A single room with no auth~~
    1. ~~Accepts connections~~
    2. ~~Clients can send in messages, server broadcasts them~~
    3. ~~Connect/Disconnect messages~~
    4. ~~Client tracking~~
	5. ~~Usernames~~
2. Authentication
    1. ~~Random code generation~~
		- [EEF Short Wordlist](https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt)
		- Use 4 words for ~2.8 trillion combinations
    2. ~~`/api/room/join` endpoint~~
    3. ~~Single-use token generation and validation~~
		- Use `crypto/rand`
	4. Send JSON instead of plain text
	5. Rate-limiting/Token-stealing
3. Multiple rooms
    1. `/api/room/create` endpoint
4. Message storage
    1. Store sent messages in an embedded database e.g. BadgerDB
    2. Send last chunk (e.g. 100 msgs) to clients on connect and more on request
    3. Remove with room
5. User Interfaces
    1. Terminal
    2. Web
6. TOML Configuration

## Connection

### Two-Step Token Exchange

1. `POST /api/room/join`: The client sends an HTTP POST request with the body containing a room and code.
2. **Issue One-Time Token**: The server verifies the credentials. If valid, it returns a temporary single-use token valid for a short time (e.g. 1 min).
3. `GET /ws?token=...`: The client initiates the WebSocket connection using the temporary token in the param. The server validates and invalidates the token upon upgrading.

#### Benefits
- Secure auth
- Direct codes not sent in url

## API

### `GET /ws`

WebSocket connections and communications

#### Params

- `token`: The one-time token from `/api/room/join`
- `username?`: Displayed name for client — will get a random 4-digit number added to avoid conflicts
    - **Default:** `User`
    - E.g. `Eric` → `Eric1234`

### `POST /api/room/join`

Join a room

#### Body

- `code`: The auth code for the room

#### Successful Response

```jsonc
{
    // Name of room associated with code
    "room": "",

    // A one-time token used to connect to the room.
    // The server knows what room is associated with it.
    "token": "",
}
```

#### Token (OTP)

- A 6-digit pseudo-random number
- 30s TTL
- **Associated metadata**
	- Room code (used as the ID of the room)
	- Connecting IP (to avoid token-stealing)

As soon as the client sends a request to `/ws`, the token is invalidated, before the connection is upgraded.

### `GET /api/room/create`

Create a new room (if allowed by server configuration)

#### Params

- `name?`: The desired name of the room
    - **Default:** `Room123456`
    - Rooms are stored via there codes, meaning there will not be naming conflicts
	- If a generated code is already in use, it is regenerated

#### Successful Response

```jsonc
{
    // Same as requested; verification
    "name": "",

    // The code for the room
    "code": "",
}
```

Client automatically gets a prompt to join the room (asked for optional username)

## Client Interfaces (Frontend)

### Terminal

### Web

## Configuration Options

- `address`: The address for the server to listen on

- `allow_new_rooms`: Allow arbitrary room creation

- `delete_empty_rooms`: If empty rooms should be removed

- `room_ttl`: How long a room is kept for once it is empty. `expire_rooms` must be `true`.
