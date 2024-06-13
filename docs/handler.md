# handler

## crate_room

部屋の作成。  
リクエストしてきたユーザーのIDとナノ秒のUNIXタイムスタンプを結合後ハッシュ化して、N桁のRoomIDを生成。  
RoomIDが衝突していないことを確認後[^1]、部屋を作成しリクエストしてきたユーザーをGameMasterとして参加させる。  

[^1]: もし衝突していれば再生成を行う。もし10回再生成しても衝突していればエラーを返す。  

### Request

`/api/create_room` に `{id, name, lang}` を含んだPOSTリクエストを投げる。  

#### Example

```sh
curl -X POST http://localhost:8080/api/create_room \
                 -H "Content-Type: application/json" \
                 -d '{"id": "1", "name": "hoge", "lang": "ja"}'
```

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | Room created successfully | 処理が全て正常に終了 |
| 400 | Invalid request payload | 受信したデータが不正 |
| 400 | id is required | ID が空 |
| 400 | name is required | name が空 |
| 400 | lang is required | lang が空 |
| 405 | Method Not Allowed | HTTP MethodがPOST以外 |
| 500 | Failed to creation room | サーバー内部で部屋の作成に失敗 |

## join_room

作成済みの部屋への参加。  
リクエストしてきたユーザー情報を解析し、指定されたRoomIDに部屋が存在することを確認後、ユーザーをその部屋に参加させる。  

### Request

`/api/join_room?room_id=<room_id>` に `{id, name, lang}` を含んだPOSTリクエストを投げる。

#### Example

```sh
curl -X POST "http://localhost:8080/api/join_room?room_id=123456" \
                 -H "Content-Type: application/json" \
                 -d '{"id": "2", "name": "fuga", "lang": "en"}'
```

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | Successfully joined the room | 処理が全て正常に終了 |
| 400 | Invalid request payload | 受信したデータが不正 |
| 400 | id is required | ID が空 |
| 400 | name is required | name が空 |
| 400 | lang is required | lang が空 |
| 404 | Room does not exist | 指定されたRoomIDが存在しない |
| 405 | Method Not Allowed | HTTP MethodがPOST以外 |
| 500 | Failed to join room | サーバー内部で部屋への参加に失敗 |

## get_room_users

部屋内のユーザー情報の取得。  
指定されたRoomIDに基づいて部屋内のユーザー情報を取得し、JSON形式で返す。  

### Request

`/api/get_room_users?room_id=<room_id>` にGETリクエストを送る。

#### Example

```sh
curl -X GET "http://localhost:8080/api/get_room_users?room_id=123456"
```

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | JSON形式のユーザー情報 | 処理が全て正常に終了 |
| 404 | Room does not exist | 指定されたRoomIDが存在しない |
| 405 | Method Not Allowed | HTTP MethodがGET以外 |
| 500 | Failed to get users in the room | サーバー内部でユーザー情報の取得に失敗 |

## delete_room

部屋の削除。  
指定されたRoomIDに基づいて部屋を削除する。  

### Request

`/api/delete_room?room_id=<room_id>` にDELETEリクエストを送る。

#### Example

```sh
curl -X DELETE "http://localhost:8080/api/delete_room?room_id=123456"
```

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | Room deleted successfully | 処理が全て正常に終了 |
| 404 | Room does not exist | 指定されたRoomIDが存在しない |
| 405 | Method Not Allowed | HTTP MethodがDELETE以外 |
| 500 | Failed to delete room | サーバー内部で部屋の削除に失敗 |
