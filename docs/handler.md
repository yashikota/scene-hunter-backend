# handler

## crate_room

部屋の作成。  
リクエストしてきたユーザーのIDとナノ秒のUNIXタイムスタンプ結合後ハッシュ化して、N桁のRoomIDを生成。  
RoomIDが衝突していないことを確認後、部屋を作成しリクエストしてきたユーザーをGameMasterとして参加させる。  

### Response

| Status Code | Response | cause |
| ----------- | -------- | ----- |
| 200 | Room created successfully | 処理が全て正常に終了 |
| 400 | Invalid request payload | 受診したデータが不正 |
| 400 | id is required | ID が空 |
| 400 | name is required | name が空 |
| 400 | lang is required | lang が空 |
| 500 | Failed to creation room | サーバー内部で部屋の作成に失敗 |
