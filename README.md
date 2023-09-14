# こんぶくん🤖

<img align="right" src="https://raw.githubusercontent.com/ebiiim/conbukun/main/assets/conbu.jpg" alt="conbukun" width="" height="100" />

[![Static Badge](https://img.shields.io/badge/add%20to%20Discord-7289DA?logo=discord&labelColor=FFFFFF)](https://discord.com/oauth2/authorize?client_id=1151028506470404096&scope=bot&permissions=11264) [![Static Badge](https://img.shields.io/badge/add%20to%20Discord%20(dev)-7289DA?logo=discord&labelColor=FFFFFF)
](https://discord.com/oauth2/authorize?client_id=1151570933543342101&scope=bot&permissions=11264)

Albion Onlineのギルド [Dog The Boston](https://twitter.com/DogTheBoston) 用のDiscord Botです。

## 使い方

`/help` で使い方を表示します。

<details>

<summary>詳しく見る（クリック）</summary>

> ## コマンド
> - `/help` このメッセージを表示します。
> - `/mule` ラバに関するヒントをランダムに投稿します（30秒後に自動削除）。
> ## リアクション
> - `リアクション集計` 集計したいメッセージにリアクション（🤖）を行うとリマインダーを投稿します（2分後に自動削除）。
> ## おまけ
> - 呼びかけに反応したりお昼寝したりします。

</details>


## リリースノート

- 今後の課題
  - 機能: 会話（？）
  - 内部: マルチインスタンス対応（{リアクション|メンション}ハンドラのGuild ID対応）
- 2023-xx-xx v0.x.x
  - 機能: `おたのしみ` プレゼンスぐるぐる
  - 機能: `おたのしみ` 反応いろいろ強化
  - 改善: `リアクション集計` どの投稿に対する集計かがわからなくなるのでリプライにした
  - 改善: `リアクション集計` 名前の順番を固定した（表示名昇順）
- 2023-09-14 v0.2.1
  - 改善: `リアクション集計` 誰もメンションされていない場合は反応しない
  - 改善: `リアクション集計` 投稿後にemojiをリセットする
- 2023-09-14 v0.2.0
  - 機能: `リアクション集計` 未反応の人をリストする
  - 機能: `おたのしみ` 話しかけられたら反応する
  - 修正: 一部のemojiがうまく表示できないバグ
- 2023-09-14 v0.1.1
  - 修正: ログが徐々に長くなる恐ろしいバグ
- 2023-09-13 v0.1.0
  - 記念すべき初回リリース
  - 機能: `/help` `/mule` `リアクション集計`
