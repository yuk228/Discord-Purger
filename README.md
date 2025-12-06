# Discord Purger

Discord で自分のメッセージを消去する bot

アカウント消す時に使うかもしれません。

## Setup

Docker が必要です。

1. `.env.example`を`.env`に rename

2. `docker compose up --build`

- `PREFIX`: Command Prefix
- `TOKEN`: Discord Token
- `OWNER_IDS`: あなたの Discord ID

## Commands

- `purge [channel_id]`: 指定されたチャンネルのメッセージを削除します
- `purge2 [guild_id]`: Discord の search api を使用してメッセージ取得/削除します

`purge`はチャンネル内のメッセージを 100 件ずつ全て取得し、削除していきます。

`purge2`は指定した guild_id のメッセージを全て削除しますが、**なぜか一回じゃ消しきれません**

`Deleted 0 messages`になるまで何回か実行して下さい。

また、削除について RateLimit を考慮していませんが、messages 関連の api が 429 返ってくるだけで特にペナルティはありません。少し待てば直ります。

## 注意事項

- RateLimit に引っかかってアカウントが制限を受けるかもしれません
- selfbot として使用した場合、Discord の規約に抵触します
- いかなる場合も、責任は実行者にあります
