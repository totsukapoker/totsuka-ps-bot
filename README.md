[![CircleCI](https://circleci.com/gh/totsukapoker/totsuka-ps-bot.svg?style=svg)](https://circleci.com/gh/totsukapoker/totsuka-ps-bot)

# totsuka-ps-bot

[戸塚ポーカースクール](https://totsukapoker.com)のゲームで使用する点数管理LINE Botです。

![LINEアカウントQRコード](/static/qrcode.png)

https://totsuka-ps-bot.herokuapp.com/

## Require

- Go
- MySQL 5.7

## Development

### Docker Compose

```shell
$ git clone git@github.com:totsukapoker/totsuka-ps-bot.git
$ cd totsuka-ps-bot
$ docker-compose up -d
$ open http://localhost:8000
```

- MySQL がポート `53306` に立ち上がります (`root` パスワード無し)
  - 各種 GUI クライアント ([Sequel Pro](https://www.sequelpro.com/) など) から接続する際にどうぞ
- MySQL に database `totsuka_ps_bot` が必要なので初回時には手動で作ってください
  - `$ docker-compose exec mysql echo 'CREATE DATABASE totsuka_ps_bot;' | mysql`

## Production

- `master` に push されたものが Heroku に自動デプロイされます。
- [リリース](https://github.com/totsukapoker/totsuka-ps-bot/releases) は雰囲気でつけてるよ（雰囲気とは

## Misc

- [heroku/go-getting-started](https://github.com/heroku/go-getting-started) をベースに開発
