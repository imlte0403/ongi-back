# WebSocket 실시간 채팅 API v3 - 간단 가이드

## WebSocket 연결

### 엔드포인트
```
ws://localhost:3000/ws/chat/:roomId?user_id=USER_ID
```

### 파라미터
- `roomId` (path): 채팅방 ID
- `user_id` (query): 사용자 ID

### 연결 예시
```
ws://localhost:3000/ws/chat/1?user_id=123
```

---

## 수신 메시지 타입

WebSocket으로부터 받게 되는 메시지 형식입니다.

### 1. 새 메시지 (message)

**발생 시점**: 누군가가 메시지를 보냈을 때

```json
{
  "type": "message",
  "room_id": 1,
  "user_id": 2,
  "data": {
    "id": 15,
    "chat_room_id": 1,
    "user_id": 2,
    "message": "안녕하세요!",
    "message_type": "text",
    "file_url": null,
    "is_read": false,
    "created_at": "2024-11-13T16:30:00Z",
    "updated_at": "2024-11-13T16:30:00Z",
    "user": {
      "id": 2,
      "email": "user2@example.com",
      "name": "김철수"
    }
  }
}
```

### 2. 읽음 처리 (read)

**발생 시점**: 누군가가 메시지를 읽음 처리했을 때

```json
{
  "type": "read",
  "room_id": 1,
  "user_id": 3,
  "data": {
    "user_id": 3,
    "last_read_at": "2024-11-13T16:35:00Z"
  }
}
```

### 3. 멤버 추가 (member_join)

**발생 시점**: 새로운 멤버가 채팅방에 추가되었을 때

```json
{
  "type": "member_join",
  "room_id": 1,
  "user_id": 5,
  "data": {
    "id": 10,
    "chat_room_id": 1,
    "user_id": 5,
    "role": "member",
    "joined_at": "2024-11-13T16:40:00Z",
    "last_read_at": null,
    "unread_count": 0,
    "created_at": "2024-11-13T16:40:00Z",
    "user": {
      "id": 5,
      "email": "user5@example.com",
      "name": "이영희"
    }
  }
}
```

### 4. 멤버 제거 (member_leave)

**발생 시점**: 멤버가 채팅방에서 제거되었을 때

```json
{
  "type": "member_leave",
  "room_id": 1,
  "user_id": 5,
  "data": {
    "user_id": 5
  }
}
```

### 5. 멤버 접속 (member_online)

**발생 시점**: 멤버가 WebSocket에 연결되었을 때

```json
{
  "type": "member_online",
  "room_id": 1,
  "user_id": 4,
  "data": {
    "user_id": 4,
    "status": "online"
  }
}
```

### 6. 멤버 접속 종료 (member_offline)

**발생 시점**: 멤버가 WebSocket 연결을 종료했을 때

```json
{
  "type": "member_offline",
  "room_id": 1,
  "user_id": 4,
  "data": {
    "user_id": 4,
    "status": "offline"
  }
}
```

---

## HTTP API 연동

WebSocket은 **수신 전용**입니다. 메시지 전송, 읽음 처리, 멤버 추가/제거는 **HTTP API**를 사용하세요.

### 1. 메시지 전송

**요청**
```http
POST /api/v1/chat/rooms/:roomId/messages
Content-Type: application/json

{
  "user_id": 123,
  "message": "안녕하세요!",
  "message_type": "text",
  "file_url": null
}
```

**응답**
```json
{
  "success": true,
  "message": "Message sent successfully",
  "data": {
    "id": 15,
    "chat_room_id": 1,
    "user_id": 123,
    "message": "안녕하세요!",
    "message_type": "text",
    "file_url": null,
    "is_read": false,
    "created_at": "2024-11-13T16:30:00Z",
    "updated_at": "2024-11-13T16:30:00Z",
    "user": {
      "id": 123,
      "email": "user@example.com",
      "name": "홍길동"
    }
  }
}
```

**자동 브로드캐스트**: 모든 WebSocket 연결에 `type: "message"` 전송

---

### 2. 읽음 처리

**요청**
```http
POST /api/v1/chat/rooms/:roomId/read
Content-Type: application/json

{
  "user_id": 123
}
```

**응답**
```json
{
  "success": true,
  "message": "Messages marked as read"
}
```

**자동 브로드캐스트**: 모든 WebSocket 연결에 `type: "read"` 전송

---

### 3. 멤버 추가

**요청**
```http
POST /api/v1/chat/rooms/:roomId/members
Content-Type: application/json

{
  "user_id": 456
}
```

**응답**
```json
{
  "success": true,
  "message": "Member added successfully",
  "data": {
    "id": 10,
    "chat_room_id": 1,
    "user_id": 456,
    "role": "member",
    "joined_at": "2024-11-13T16:40:00Z",
    "last_read_at": null,
    "unread_count": 0,
    "created_at": "2024-11-13T16:40:00Z",
    "user": {
      "id": 456,
      "email": "user456@example.com",
      "name": "이영희"
    }
  }
}
```

**자동 브로드캐스트**: 모든 WebSocket 연결에 `type: "member_join"` 전송

---

### 4. 멤버 제거

**요청**
```http
DELETE /api/v1/chat/rooms/:roomId/members/:userId
```

**응답**
```json
{
  "success": true,
  "message": "Member removed successfully"
}
```

**자동 브로드캐스트**: 모든 WebSocket 연결에 `type: "member_leave"` 전송

---

## 플로우 예시

### 메시지 전송 플로우

```
1. 클라이언트 A → HTTP POST /api/v1/chat/rooms/1/messages
   Body: { "user_id": 123, "message": "안녕!" }

2. 서버 → DB에 메시지 저장

3. 서버 → WebSocket Hub로 브로드캐스트

4. WebSocket → 모든 연결된 클라이언트 (A, B, C)에게 전송
   {
     "type": "message",
     "room_id": 1,
     "user_id": 123,
     "data": { ... }
   }

5. 클라이언트 A, B, C → 실시간으로 메시지 수신
```

### 읽음 처리 플로우

```
1. 클라이언트 B → HTTP POST /api/v1/chat/rooms/1/read
   Body: { "user_id": 456 }

2. 서버 → DB에 읽음 시간 업데이트

3. 서버 → WebSocket Hub로 브로드캐스트

4. WebSocket → 모든 연결된 클라이언트에게 전송
   {
     "type": "read",
     "room_id": 1,
     "user_id": 456,
     "data": { "user_id": 456, "last_read_at": "..." }
   }

5. 클라이언트 A, C → 실시간으로 읽음 표시 업데이트
```

### 멤버 추가 플로우

```
1. 클라이언트 A → HTTP POST /api/v1/chat/rooms/1/members
   Body: { "user_id": 789 }

2. 서버 → DB에 멤버 추가

3. 서버 → WebSocket Hub로 브로드캐스트

4. WebSocket → 모든 연결된 클라이언트에게 전송
   {
     "type": "member_join",
     "room_id": 1,
     "user_id": 789,
     "data": { ... }
   }

5. 클라이언트들 → 실시간으로 멤버 목록 업데이트
```

### 멤버 제거 플로우

```
1. 클라이언트 A → HTTP DELETE /api/v1/chat/rooms/1/members/789

2. 서버 → DB에서 멤버 제거

3. 서버 → WebSocket Hub로 브로드캐스트

4. WebSocket → 모든 연결된 클라이언트에게 전송
   {
     "type": "member_leave",
     "room_id": 1,
     "user_id": 789,
     "data": { "user_id": 789 }
   }

5. 클라이언트들 → 실시간으로 멤버 목록 업데이트
```

---

## 요약

### WebSocket 역할
- **실시간 메시지 수신만** 담당
- 연결만 유지하면 자동으로 모든 이벤트 수신

### HTTP API 역할
- 메시지 전송
- 읽음 처리
- 멤버 추가/제거
- **모든 쓰기 작업**은 HTTP API 사용

### 메시지 타입
| 타입 | 설명 | 트리거 |
|------|------|--------|
| `message` | 새 메시지 | POST /messages |
| `read` | 읽음 처리 | POST /read |
| `member_join` | 멤버 추가 | POST /members |
| `member_leave` | 멤버 제거 | DELETE /members/:userId |
| `member_online` | 접속 | WebSocket 연결 시 |
| `member_offline` | 접속 종료 | WebSocket 종료 시 |

---

## 클라이언트 구현 예시 (간단)

### Flutter
```dart
// WebSocket 연결
final channel = WebSocketChannel.connect(
  Uri.parse('ws://localhost:3000/ws/chat/1?user_id=123')
);

// 메시지 수신
channel.stream.listen((message) {
  final data = jsonDecode(message);
  print('Type: ${data['type']}');
  print('Data: ${data['data']}');
});

// 메시지 전송 (HTTP)
await http.post(
  Uri.parse('http://localhost:3000/api/v1/chat/rooms/1/messages'),
  body: jsonEncode({'user_id': 123, 'message': '안녕!'}),
);
```

### JavaScript
```javascript
// WebSocket 연결
const ws = new WebSocket('ws://localhost:3000/ws/chat/1?user_id=123');

// 메시지 수신
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Type:', data.type);
  console.log('Data:', data.data);
};

// 메시지 전송 (HTTP)
fetch('http://localhost:3000/api/v1/chat/rooms/1/messages', {
  method: 'POST',
  body: JSON.stringify({ user_id: 123, message: '안녕!' }),
});
```

### Python
```python
import websocket
import requests
import json

# WebSocket 연결
def on_message(ws, message):
    data = json.loads(message)
    print(f"Type: {data['type']}")
    print(f"Data: {data['data']}")

ws = websocket.WebSocketApp(
    "ws://localhost:3000/ws/chat/1?user_id=123",
    on_message=on_message
)
ws.run_forever()

# 메시지 전송 (HTTP)
requests.post(
    'http://localhost:3000/api/v1/chat/rooms/1/messages',
    json={'user_id': 123, 'message': '안녕!'}
)
```

---

끝!
