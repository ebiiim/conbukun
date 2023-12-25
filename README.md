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
- [よくある質問（Frequently Asked Questions）](#%E3%82%88%E3%81%8F%E3%81%82%E3%82%8B%E8%B3%AA%E5%95%8Ffrequently-asked-questions)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 使い方

`/help` で使い方を表示します。

<details><summary>詳しく見る（クリック）</summary>

> ## コマンド
> - `/help` このメッセージを表示します。
> - `/mule` ラバに関するヒントをランダムに投稿します（30秒後に自動削除）。
> - `/route-add` アバロンのルートを追加します。
> - `/route-mark` マップの色を変えたりメモを載せたりします。
> - `/route-list` いま持ってる情報を表示します（確認）。
> - `/route-print` アバロンのルートを画像で投稿します（共有）。
> - `/route-clear` アバロンのルートをリセットします。
> ## リアクション
> - `リアクション集計` 集計したいメッセージにリアクション（🤖）を行うとリマインダーを投稿します（2分後に自動削除）。
> ## おまけ
> - 呼びかけに反応したりお昼寝したりします。

</details>


## リリースノート

- 2023-12-26 v1.8.0
  - 機能: `ルートナビ` 確認用コマンド `/route-list` を追加
- 2023-12-24 v1.7.0
  - 改善: `ルートナビ` `/route-clear` 時にバックアップを取得＆マークは削除せずルートだけを削除するようにした
  - 機能: `ルートナビ` `/route-mark` にユーザ名を追加（情報提供目的）
  - 修正: `ルートナビ` 細かい改良や修正

<details><summary>続きを表示する（クリック）</summary>

- 2023-12-19 v1.6.0
  - 機能: `ルートナビ` 色指定とコメントを追加
  - 改善: `ルートナビ` ルートが多いときに画像生成に失敗する問題（タイムアウト値を増やした）
- 2023-12-05 v1.5.0
  - 改善: `ルートナビ` ルートが多くなると画像が横長になるので、ルートが32個以上あるときの描画形式を変更した
  - 修正: `ルートナビ` バージョンがあがるとセーブデータが引き継がれないバグ
- 2023-12-01 v1.4.1
  - 改善: タイムゾーンをJSTに変更（Kubernetesマニフェスト）
- 2023-12-01 v1.4.0
  - 修正: `ルートナビ` ルート数が多くなると3秒以内に応答できずDiscordの制約に引っかかる問題
  - 機能: `ルートナビ` ルートをクリアするコマンドを追加
  - 機能: `ルートナビ` 見た目を改善する（HO指定、色分け、etc）
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

### アップグレードガイド

> [!NOTE]
> 対応が必要なバージョンのみ記載しています（これらのバージョンを経由してください）。

- v1.5.x to v1.6.x
  - v1.6.xで `/route-mark` を実行するとセーブデータが移行されます。これをやらないとマーク情報が引き継がれず、さらにゴミが残ります。
- v1.4.x to v1.5.x
  - セーブデータの移行は手動です。各キーの値（ `jq keys[]` ）と各項目のname（ `jq .[].name` ）を `ギルド名#チャネル名` の形式に変更してください（後ろの ` (conbukun@v1.4.x)` を削除、重複は手動でマージ）。

</details>

### その他

- 既知のバグ
  - なし
- 今後の課題
  - 機能: 会話AIほしくない？
  - 改善: 情報やエラーメッセージをembedでキレイに表示したい
  - 内部: マルチインスタンス対応（{リアクション|メンション}ハンドラのGuild ID対応）
  - 機能: `/route-clear` にundoほしくない？（ほぼ使われていないので優先度低）

## ライセンス

- ソースコード: [MIT License](https://github.com/ebiiim/conbukun/blob/main/LICENSE)
- ライブラリおよびサブモジュール: それぞれのライセンスを参照のこと
- 画像: `assets/README.md` を参照のこと

---

> [!NOTE]
> 以降は稼働中のサービス（Discord Bot）に関する情報です。

---

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

## よくある質問（Frequently Asked Questions）

##### Q: こんぶくんって誰ですか？

##### A: 真のギルドマスターです。

##### Q: ルートナビ機能はゲームの利用規約に抵触しますか？

##### A: いいえ。理由は次のとおりです。

1. [こちらの告知](https://forum.albiononline.com/index.php/Thread/135576-Roads-of-Avalon-Mapping-GPS-Tools/)にある「bannable offence」に該当しない。
   - ユーザは**手動**で情報（ゾーン名、ポータルの種類、残り時間など）を入力する。
   - ゲーム内のリアルタイムの情報を自動的に取得**しない**（例：画面や通信の解析など）。
   - 一般的に、本機能がもたらす利益は紙とペンやスプレッドシートと同等であると考えられる。

2. [こちらの告知](https://forum.albiononline.com/index.php/Thread/124819-Regarding-3rd-Party-Software-and-Network-Traffic-aka-do-not-cheat-Update-16-45-U/)にあるチート行為に該当しない。
   - ゲームクライアントを改変**しない**。
   - オーバーレイ**ではない**。
   - 補助的に用いる情報は、[利用が問題ないと判断された](https://forum.albiononline.com/index.php/Thread/124819-Regarding-3rd-Party-Software-and-Network-Traffic-aka-do-not-cheat-Update-16-45-U/?postID=1001172#post1001172)「[Albion Data Project](https://www.albion-online-data.com/)」のみから得た。

※ ご不明な点がある場合は `mail@ebiiim.com` までお問い合わせください。

---

##### Q: Who is conbukun?

##### A: The TRUE guild master.

##### Q: Does the route navigation feature violate the game's ToS?

##### A: No. Here's why.

1. It does not fall under the "bannable offence" mentioned in [this announcement](https://forum.albiononline.com/index.php/Thread/135576-Roads-of-Avalon-Mapping-GPS-Tools/).
   - The user **manually** enters information (zone name, portal type, remaining time, etc.).
   - It does **not** automatically obtain real-time information from the game (e.g. screen or network traffic analysis).
   - In general, the benefits of this feature are considered to be equivalent to pen and paper or spreadsheets.

2. It does not fall under the "do not cheat" policy mentioned in [this announcement](https://forum.albiononline.com/index.php/Thread/124819-Regarding-3rd-Party-Software-and-Network-Traffic-aka-do-not-cheat-Update-16-45-U/).
    - It does **not** modify the game client.
    - It is **not** an overlay tool.
    - The information supplementary used is obtained only from the "[Albion Data Project](https://www.albion-online-data.com/)", which is [considered to be safe to use](https://forum.albiononline.com/index.php/Thread/124819-Regarding-3rd-Party-Software-and-Network-Traffic-aka-do-not-cheat-Update-16-45-U/?postID=1001172#post1001172).

※ Please contact `mail@ebiiim.com` if anything is unclear.
