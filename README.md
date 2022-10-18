## Golang Simple Line Bot Practice( gin+mongoDB )

### 步驟
1. 修改配置，config/config.yaml
* MongoDB 連線設定
* Linebot SDK 設定(從line developer後台取得)
2. 使用ngrok，讓local取得https網址
   ![](https://i.imgur.com/wV65g63.png)

3. 至line developer後台，在Messaging API將網址設定到Webhook URL
   ![](https://i.imgur.com/oW2rZQL.png)

4.
```
go run main.go
```
5. Api網址
* 發送文字訊息api
  http://localhost/api/v1/sendMessage/：userId?message=測試發送2
* 取得user api
  http://localhost/api/v1/users
* 取得webhook line message api
  http://localhost/api/v1/lineMessages