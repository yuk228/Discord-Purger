# Discord Purger
Discordで自分のメッセージを消去するbot

アカウント消す時に使うかもしれません。

## Setup
Dockerが必要です。

1. `.env.example`を`.env`にrename

2. `docker compose up --build`

- `PREFIX`: Command Prefix
- `TOKEN`: Discord Token
- `OWNER_IDS`: あなたのDiscord ID

## Commands

- `purge [channel_id]`: 指定されたチャンネルのメッセージを削除します
- `purge2 [channel_id]`: Discordのsearch apiを使用してメッセージ取得/削除します

`purge`はチャンネル内のメッセージを100件ずつ全て取得していきます。

`purge2`は現在DM, グループでのみ使用可能です。

いずれもチャンネル内のメッセージ全てが対象です。
