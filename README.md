# こんぶくん🤖

<img align="right" src="https://raw.githubusercontent.com/ebiiim/conbukun/main/assets/icon/conbu.jpg" alt="conbukun" width="" height="100" />

<!-- https://discordapi.com/permissions.html#60416 -->
[![Static Badge](https://img.shields.io/badge/add%20to%20Discord-7289DA?logo=discord&labelColor=FFFFFF)](https://discord.com/oauth2/authorize?client_id=1151028506470404096&scope=bot&permissions=60416)
[![Static Badge](https://img.shields.io/badge/devs%20only-7289DA?logo=discord&labelColor=FFFFFF)
](https://discord.com/oauth2/authorize?client_id=1151570933543342101&scope=bot&permissions=60416)
[![Release (GitHub)](https://img.shields.io/github/v/release/ebiiim/conbukun)](https://github.com/ebiiim/conbukun/releases/latest)

<a href="https://www.buymeacoffee.com/ebiiim" target="_blank"><img src="https://raw.githubusercontent.com/ebiiim/conbukun/main/assets/doc/buymeacoffee.png" alt="Buy Me A Coffee" width="" height="30" /></a>

Albion Onlineのギルド [Dog The Boston](https://twitter.com/DogTheBoston) 用のDiscord Botです。


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [使い方](#%E4%BD%BF%E3%81%84%E6%96%B9)
- [リリースノート](#%E3%83%AA%E3%83%AA%E3%83%BC%E3%82%B9%E3%83%8E%E3%83%BC%E3%83%88)
- [ライセンス](#%E3%83%A9%E3%82%A4%E3%82%BB%E3%83%B3%E3%82%B9)
- [利用規約（Terms of Service）](#%E5%88%A9%E7%94%A8%E8%A6%8F%E7%B4%84terms-of-service)
- [プライバシーポリシー（Privacy Policy）](#%E3%83%97%E3%83%A9%E3%82%A4%E3%83%90%E3%82%B7%E3%83%BC%E3%83%9D%E3%83%AA%E3%82%B7%E3%83%BCprivacy-policy)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 使い方

`/help` で使い方を表示します。

<details>

<summary>詳しく見る（クリック）</summary>

> ## コマンド
> - `/help` このメッセージを表示します。
> - `/mule` ラバに関するヒントをランダムに投稿します（30秒後に自動削除）。
> - `/route-add` アバロンのルートを追加します。
> - `/route-print` アバロンのルートを画像で投稿します。
> - `/route-clear` アバロンのルートをリセットします。
> ## リアクション
> - `リアクション集計` 集計したいメッセージにリアクション（🤖）を行うとリマインダーを投稿します（2分後に自動削除）。
> ## おまけ
> - 呼びかけに反応したりお昼寝したりします。

</details>


## リリースノート

- 今後の課題
  - 機能: 会話（？）
  - 内部: マルチインスタンス対応（{リアクション|メンション}ハンドラのGuild ID対応）
  - 改善: `ルートナビ` エラーメッセージをembedでキレイに出力する
  - 機能: `ルートナビ` 見た目を改善する（HO指定、色分け、etc）
- 2023-xx-xx v1.4.0
  - 修正: `ルートナビ` ルート数が多くなると3秒以内に応答できずDiscordの制約に引っかかる問題
  - 機能: `ルートナビ` ルートをクリアするコマンドを追加
  - 機能: 保存機能を追加（起動時にデータをロード＆終了時にデータをセーブ） 
- 2023-11-28 v1.3.0
  - 機能: `ルートナビ` こんぶくんがみんなのためにアバロンのルートを覚えてくれるようになった
- 2023-10-28 v1.2.0
  - 機能: `おたのしみ` ラバコマンドのレスポンスの種類が増えた（ハロウィン＆シェイプシフター）
- 2023-09-21 v1.1.0
  - 機能: `おたのしみ` こんぶくんがリプライに反応するようになった
  - 機能: `おたのしみ` こんぶくんがAlbion Onlineをプレイするようになった（プレゼンス）
- 2023-09-19 v1.0.0
  - 安定しているので正式リリース
  - サービス: 利用規約とプライバシーポリシーを作成
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

## ライセンス

- ソースコード: [MIT License](https://github.com/ebiiim/conbukun/blob/main/LICENSE)
- ライブラリおよびサブモジュール: それぞれのライセンスを参照のこと
- 画像: `assets/README.md` を参照のこと

## 利用規約（Terms of Service）

発効日: 2023年9月19日<br>
最終更新日: 2023年9月19日

conbukun（本サービス）を利用する場合、ユーザーは次の利用規約に同意したものとします。

- サービスを利用する権利: ユーザーは本サービスを[Discordの利用規約](https://discord.com/terms)および参加しているサーバーのルールに違反しない目的においてのみ利用できます。
- 個人情報の扱い: ユーザーは、[プライバシーポリシー](#%E3%83%97%E3%83%A9%E3%82%A4%E3%83%90%E3%82%B7%E3%83%BC%E3%83%9D%E3%83%AA%E3%82%B7%E3%83%BCprivacy-policy)をよく読んで理解した上で同意する必要があります。

---

Effective date: September 19, 2023<br>
Last updated: September 19, 2023

By using conbukun（the Service） you automatically agree to the Terms of Service below.

- Rights to use the service: You have the right to use the Service as long as you don't use it in any way that would break [Discord ToS](https://discord.com/terms) or the rules of the guild ("server") you are in.
- Handling of personal data: You have to read, understand, and agree to the [Privacy Policy](#%E3%83%97%E3%83%A9%E3%82%A4%E3%83%90%E3%82%B7%E3%83%BC%E3%83%9D%E3%83%AA%E3%82%B7%E3%83%BCprivacy-policy).


## プライバシーポリシー（Privacy Policy）

発効日: 2023年9月19日<br>
最終更新日: 2023年9月19日

conbukun（本サービス）のプライバシーポリシーは次のとおりです。

- データ閲覧: サービス提供のために、本サービスは本サービスのボットへのダイレクトメッセージおよび読み取り権限が付与されたチャンネルに限り、メッセージおよびリアクションを処理します。
- データ収集: 分析のために、本サービスは本サービスのボットへのダイレクトメッセージおよび読み取り権限が付与されたチャンネルに限り、メッセージおよびリアクションを匿名化した上で収集します。

---

Effective date: September 19, 2023<br>
Last updated: September 19, 2023

The Privacy Policy of conbukun (the Service) is as follows.

- View of data: In order to provide the service, the Service processes messages and reactions, limited to direct messages to the Service's bots and channels where read permissions have been granted.
- Collection of data: For analytic purposes, the Service anonymously collects messages and reactions, limited to direct messages to the Service's bots and channels where read permissions have been granted.
