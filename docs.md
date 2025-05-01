# Matchmaking Service API & WebSocket Routes

---

## **HTTP API Routes**

### **1. Authentication**
| Route                | Method | Description                          | Request Body (Example)               | Success Response (200)              |
|----------------------|--------|--------------------------------------|---------------------------------------|-------------------------------------|
| `/auth/register`     | POST   | Register a new player                | `{ "username": "Alice", "password": "secret" }` | `{ "playerId": "123", "token": "xyz" }` |
| `/auth/login`        | POST   | Authenticate player                  | `{ "username": "Alice", "password": "secret" }` | `{ "playerId": "123", "token": "xyz" }` |
| `/auth/logout`       | POST   | End session                          | `{ "playerId": "123" }`              | `{ "message": "Logged out" }`       |

---

### **2. Matchmaking**
| Route                  | Method | Description                          | Request Body (Example)               | Success Response                   |
|------------------------|--------|--------------------------------------|---------------------------------------|------------------------------------|
| `/matchmake/join`      | POST   | Enter matchmaking queue              | `{ "playerId": "123" }`              | `{ "queueStatus": "searching" }`   |
| `/matchmake/leave`     | POST   | Leave queue                          | `{ "playerId": "123" }`              | `{ "message": "Left queue" }`      |

---

### **3. Room Management**
| Route                  | Method | Description                          | Request Body (Example)               | Success Response                   |
|------------------------|--------|--------------------------------------|---------------------------------------|------------------------------------|
| `/rooms/create`        | POST   | Create a private room                | `{ "playerId": "123" }`              | `{ "roomId": "abc", "code": "ABCD12" }` |
| `/rooms/join`          | POST   | Join room by code                    | `{ "playerId": "456", "code": "ABCD12" }` | `{ "roomId": "abc", "players": ["123", "456"] }` |
| `/rooms/leave`         | POST   | Leave room                           | `{ "playerId": "123", "roomId": "abc" }` | `{ "message": "Left room" }`       |
| `/rooms/start`         | POST   | Start game (host-only)               | `{ "playerId": "123", "roomId": "abc" }` | `{ "gameServer": "ws://game-abc" }`|

---

## **WebSocket Events**
**Connection URL:** `wss://yourserver.com/socket`
*(All events include `event` and `data` fields)*

### **1. Connection Lifecycle**
| Event            | Sent When...                  | Data Example                          |
|------------------|-------------------------------|---------------------------------------|
| `connect`        | Player connects               | `{ "playerId": "123" }`               |
| `disconnect`     | Player disconnects            | `{ "playerId": "123" }`               |

### **2. Matchmaking Events (1v1)**
| Event            | Sent When...                  | Data Example                          |
|------------------|-------------------------------|---------------------------------------|
| `match_found`    | Auto-match succeeds           | `{ "roomId": "abc", "opponentId": "456" }` |
| `queue_update`   | Queue position changes        | `{ "waitTime": 15 }`                 |

### **3. Room Events**
| Event            | Sent To...                    | Data Example                          |
|------------------|-------------------------------|---------------------------------------|
| `player_joined`  | All room members              | `{ "roomId": "abc", "newPlayerId": "456" }` |
| `player_left`    | All room members              | `{ "roomId": "abc", "playerId": "123" }` |
| `game_starting`  | All room members              | `{ "roomId": "abc" }`                |

---

## **Session Data Structure (Server-Side)**
```javascript
{
  players: {
    "player123": {
      socketId: "socketABC",
      status: "searching", // "in_lobby", "in_game"
      roomId: "room456"    // Null if not in a room
    }
  },
  rooms: {
    "room456": {
      code: "ABCD12",
      players: ["player123", "player456"], // Always 2 players max
      hostId: "player123",
      status: "waiting"    // "in_game"
    }
  },
  matchmakingQueue: ["player123", "player789"] // Simple array for 1v1
}
