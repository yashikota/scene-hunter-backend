# handler

## generate_user_id

ユーザーIDの生成。  
発行したIDはJSONのmessageの値として返却。  

### Request

`/api/generate_user_id` にGETリクエストを投げる。  
オプションで期限を設定できる。デフォルトは86400秒(1日)。  

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | {UUID} | UUIDの生成に成功 |
| 500 | Failed to generate user ID | サーバー内部でIDの生成に失敗 |

## exist_user_id

ユーザーIDが存在するか確認。  

### Request

`/api/exist_user_id` にGETリクエストを投げる。  

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | User ID exists | ユーザーが存在する |
| 404 | User ID does not exist | ユーザーが存在しない |
| 500 | Failed to check user ID | サーバー内部でIDの確認に失敗 |

## crate_room

部屋の作成。  
リクエストしてきたユーザーのIDとナノ秒のUNIXタイムスタンプを結合後ハッシュ化して、N桁のRoomIDを生成。  
RoomIDが衝突していないことを確認後[^1]、部屋を作成しリクエストしてきたユーザーをGameMasterとして参加させる。  

[^1]: もし衝突していれば再生成を行う。もし10回再生成しても衝突していればエラーを返す。  

### Request

`/api/create_room` に `{id, name, lang}` を含んだPOSTリクエストを投げる。  

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

`/api/join_room?room_id={room_id}` に `{id, name, lang}` を含んだPOSTリクエストを投げる。

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

`/api/get_room_users?room_id={room_id}` にGETリクエストを送る。

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | JSON形式のユーザー情報 | 処理が全て正常に終了 |
| 404 | Room does not exist | 指定されたRoomIDが存在しない |
| 405 | Method Not Allowed | HTTP MethodがGET以外 |
| 500 | Failed to get users in the room | サーバー内部でユーザー情報の取得に失敗 |

## change_game_master

指定されたルームのゲームマスターを変更。  

### Request

`/api/change_game_master?room_id={room_id}` に `{id}` を含んだPUTリクエストを送る。

### Response

| Status Code | Response | cause |
| - | - | - |
| 200 | Successfully changed the game master | 処理が全て正常に終了 |
| 400 | room_id is required | RoomIDがリクエストに含まれていない |
| 404 | Room does not exist | 指定されたRoomIDが存在しない |
| 404 | User does not exist | ユーザーが存在しない |
| 405 | Method Not Allowed | HTTP MethodがPUT以外 |
| 500 | Failed to check if the user exists | ユーザーの存在確認に失敗 |
| 500 | Failed to change the game master | ゲームマスターの変更に失敗 |

## delete_room_users

RoomIDで指定された部屋にいるUserの削除。  

### Request

`/api/get_room_users?room_id={room_id}` にDELETEリクエストを送る。

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | Room deleted successfully | 処理が全て正常に終了 |
| 404 | Room does not exist | 指定されたRoomIDが存在しない |
| 405 | Method Not Allowed | HTTP MethodがDELETE以外 |
| 500 | Failed to delete room | サーバー内部でユーザーの削除に失敗 |

## delete_room

部屋の削除。  
指定されたRoomIDに基づいて部屋を削除する。  

### Request

`/api/delete_room?room_id={room_id}` にDELETEリクエストを送る。

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | Room deleted successfully | 処理が全て正常に終了 |
| 404 | Room does not exist | 指定されたRoomIDが存在しない |
| 405 | Method Not Allowed | HTTP MethodがDELETE以外 |
| 500 | Failed to delete room | サーバー内部で部屋の削除に失敗 |
