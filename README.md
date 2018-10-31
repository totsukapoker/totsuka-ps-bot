[![CircleCI](https://circleci.com/gh/totsukapoker/totsuka-ps-bot.svg?style=svg)](https://circleci.com/gh/totsukapoker/totsuka-ps-bot)

# totsuka-ps-bot

[戸塚ポーカースクール](https://totsukapoker.com)のゲームで使用する点数管理LINE Botです。

![LINEアカウントQRコード](/static/qrcode.png)

https://totsuka-ps-bot.herokuapp.com/

## Require

- Go
- MySQL

## Development

```shell
$ cd totsuka-ps-bot.git
$ go build
$ ./totsuka-ps-bot
$ open http://localhost:8000
```

直接 build して実行してください。MySQL が必要なのでローカルに起動して下さい。デフォルトで `root@localhost (パスワード無し)` に接続しにいきます。DB 名は `totsuka_ps_bot` です。
DB 接続先の変更や起動ポートの変更など、環境変数の設定が必要な場合は [.env.sample](/.env.sample) を `.env` として設置し、変更ができます。

## Production

- `master` に push されたものが Heroku に自動デプロイされます。
- [リリース](https://github.com/totsukapoker/totsuka-ps-bot/releases)は雰囲気でつけてるよ（雰囲気とは

## Misc

- [heroku/go-getting-started](https://github.com/heroku/go-getting-started) をベースに開発
