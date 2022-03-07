# Go-LINEBot

リクエスト集中講義(Go入門)で作成した成果物の Wiki です。構築手順を残すために、最初の導入から記録しておきます。

# 目次

# 概要

タスクを管理してくれる、LINEBot です。

LINE が一番見ることが多いため、今回は LINE をプラットフォームにします。

「課題」と入力すると課題の一覧を返し、「予定」と入力すると登録した予定を返し、「カレンダー」と入力すると Google カレンダーからの情報を返します。

登録したのを解除しない限り、登録情報の前日にリマインドをしてくれます。

## 開発環境

- Windows11
- Go 1.17.8
- Heroku

## **LINE BOTの仕組み**

実装を始める前に、大まかなしくみと手順を説明します。今回は、勉強も兼ねてHerokuを使ってLINE BOTを作成します。

【LINE BOTの仕組み】

![Untitled](https://user-images.githubusercontent.com/59988255/157020763-64911d60-5614-4eaa-b634-caf45c5f734c.png)

1. ユーザーがボットにメッセージを送ると、Webhookを利用し、Messaging APIを通して、Herokuへリクエストを送信します。
2. Herokuが受け取ったリクエストを、Goで実装したソースコードが処理し、値をHerokuへ返します。
3. Herokuが受け取った応答リクエストを、Messaging APIへ送ります。(応答リクエストは、JSON形式でHTTPSを使って送信されます。)
4. LINEが受け取り、データが表示されます。

という流れで、データの授受が行われています。

# LINE APIの登録＆設定

まずは、**LINE APIの登録と設定**を行います。

LINE Botを作成する上では、LINE APIの使用は欠かせません。

普段使っている**LINEアカウントから簡単に登録**を進めていくことができます。

## **LINE Developers 登録**

**LINE Developers**とは、LINEを使って開発を行うためのものでこれを登録すると、APIで様々なことができるようになります。

LINEアカウントで簡単に登録できるので、[LINE Developersページ](https://developers.line.me/ja/)から**ログイン***します*。

- 右上の [ログイン]
- [LINE Business ID]
    - [LINE アカウントでログイン]
- [LINE]
    - 自分のアカウントでログインしてください

### 開発者名を登録

開発者名の登録は、単純に**自分の実名**または、**ニックネーム**などを登録すれば大丈夫です。

- [Hi, H.M.! Welcome to the LINE Developers Console.]
    - [Developer name] :  任意の名前（例：H.M.）
    - [Your email] :  自身のメールアドレス
    - [*]  [I have read and agreed to the LINE Developers Agreement]
    - [Create my account]

### 新規プロバイダーを作成

プロバイダーとは、サービス提供者（企業・個人）の名前のことですが、**名称はなんでもいい**ので基本的に個人開発の場合は開発者名と同じで大丈夫です。

- [Welcome to the LINE Developers Console!]
    - [Create a new provider]
- [Create a new provider]
    - [Provider name] :  任意の名前（例：H.M.）
    - [Create]

### チャネル作成

チャネルとは、botに対応するもので、**１つのチャネルが1つのBot**を表しています。

任意のプロバイダーから新規チャネル作成を行います。

- [This provider doesn't have any channels yet]
    - [Create a Messaging API channel]
- [Create a channel]
    - [Channel type] :  [Messaging API]
    - [Provider] :  先ほど作成した開発者名
    - [Channel icon] :  任意で設定してください（後から変えられるかも？？）
    - [Channel name] :  任意の名前（例：タスク管理Bot）
    - [Channel description] :  適当に記述（例：課題や予定を管理してくれるBotです！）
    - [Category] :  なんでもいいです（例：個人）
    - [Subcategory] :  なんでもいいです（例：個人（学生））
    - [Email address] :  自分のメールアドレス
    - [Privacy policy URL] : 任意です（空で大丈夫です）
    - [Terms of use URL] :  任意です（空で大丈夫です）
    - [*]  [I have read and agree to the LINE Official Account Terms of Use]
    - [*]  [I have read and agree to the LINE Official Account API Terms of Use]
    - [Create]
- [**Create a Messaging API channel with the following details?]**
    - [OK]
- [**情報利用に関する同意について]**
    - [同意する]

## **LINE Messaging APIの登録＆設定**

LINE Messaging APIとは、開発するLINE BotとユーザーがLINEのアカウントを通じて**相互コミュニケーションを実現**するAPIです。

ここで、いくつか先ほど作成したL**INE Messaging API のチャネルの設定**を変更しておきます。

（２つのページを行き来するため、分かるように先頭に「｛D｝」(LINE Developers)と「｛O｝」(LINE Official Account Manager) をつけておきます。）

- ｛D｝先ほど作成したチャンネルの [Messaging API] タブ
    - [LINE Official Account features]
        - [Allow bot to join group chats] > [Edit]
- ｛O｝[アカウント設定]
    - [機能の利用]
        - [チャットへの参加] :  [*] [グループ・複数人チャットへの参加を許可する]
- [設定を変更]
    - [変更]

リンク前の [LINE Developers] のタブに戻ってください。

- ｛D｝[LINE Official Account features]
    - [Auto-reply messages] > [Edit]
- ｛O｝[応答設定]
    - [基本設定]
        - [あいさつメッセージ] :  [*]  [オフ]
        - これはどっちでもいいです。友達追加した時のメッセージです。
    - [詳細設定]
        - [応答メッセージ] :  [*]  [オフ]
        - この部分を Python で実装するためオフにします。

リンク前の [LINE Developers] のタブに戻ってください。

### アクセストークン（ロングターム）

**アクセストークン**とは、**APIを使用する上で必要なトークン**のことであとでHerokuに環境変数として設定します。

ここでは、まだ何も表示されていないと思うので、発行のボタンを押してアクセストークンを発行してください。

**失効までの時間は0のままで大丈夫**です。

- ｛D｝[Channel access token]
    - [Channel access token (long-lived)] > [issue]

### Webhook送信

Webhook送信では、友だち追加やユーザーからの**メッセージ送信などのイベントが発生した際**に、Webhook URLで**リクエストを受信するか否か**を設定するので、必ず「**利用する**」を選択します。

これが利用しないとなっていると、リクエストを受信できずになんの反応もない**既読スルーなBot**になってしまいます。

- ｛D｝[LINE Official Account features]
    - [Auto-reply messages] > [Edit]
- ｛O｝[応答設定]
    - [詳細設定]
        - [Webhook] :  [*]  [オン]

リンク前の [LINE Developers] のタブに戻ってください。

### [ToDo] Webhook URL

**Webhook URL**には、LINE Platformからの**リクエストを受信するURL**を設定します。

つまり、あとで設定する**HerokuのURL**をここに入力します。そのため、ここはまだ空欄で構いません。

- ｛D｝[Webhook settings]
    - [Webhook URL] > [Edit]
    - ここに後ほど URL を入力します。

# **Herokuの登録＆設定**

**Heroku**（ヘロク）とは、**PaaS**（パース）と呼ばれるサービスで、無料で**とにかく簡単にWebサービスを公開する**ことができます。

PaaS（パース）とは、「**Platform as a Service**」の略で、Webサービスを公開するために必要なものをあらかじめ用意してくれるサービスを指します。

- Webサーバー
- OS
- データベース
- プログラミング言語の実行環境

などが具体的には含まれています。

今回は、LINE BotをLINE Platformからの**リクエストを受信するURL**（Webhook URL）を作るためにHerokuを採用しています。

## **Herokuにログイン**

まずは、[Heroku](https://id.heroku.com/login)にアクセスし、アカウント登録します。

- [Log in to your account]
    - [New to Heroku? Sign Up]
- [Sign up for free and experience Heroku today]
    - [First name] :  任意の名前（例：Masaki）
    - [Last name] :  任意の名前（例：Hibino）
    - [Email address] :  自分のメールアドレス
    - [Company name] :  任意です（空で大丈夫です）
    - [Role] :  なんでもいいです（例：Student）
    - [Country] :  Japan
    - [Primary development language] :  なんでもいいです（例：Python）
    - [*] [I'm not a robot]
    - [CREAT FREE ACCOUNT]

登録したメールアドレスにメールが届きます。

その、URLをクリックしてパスワードの設定を行ってください。

[Welcome to Heroku] と表示されれば完了です。

## **Heroku CLIインストール**

下記のサイトで、Windows版をインストールするとコマンドラインが使えるようになります。（違うOSの場合はそれに合わせてください。）

[https://devcenter.heroku.com/articles/heroku-cli#download-and-install](https://devcenter.heroku.com/articles/heroku-cli#download-and-install)

[64-bit Installer] をクリックして、ダウンロードします。ダウンロードしたファイルを、ダブルクリックしてインストールを進めます。

- [Choose Components]
    - 全てにチェック
    - [Next]
- [Choose Install Location]
    - 任意の場所にインストールしてください
    - [Install]

終わったら [Close] から終了してください。

## **Herokuにログイン**

インストールが完了しましたら、Herokuにログインします。

Git bash を起動し、今インストールしたホームディレクトリに行き、以下のコマンドを実行します。

```bash
$ cd /e/hosei_pc/Heroku/
$ heroku login
# 何か押して、エンター
```

- [Log in to the Heroku CLI]
    - [Log In]
- [Log in to your account]
    - ログインしてください

Git Bash に戻って、[Logged in as ********@gmail.com] と出てれば、完了です。

## **アプリケーション登録**

次に、アプリケーションを登録します。(あとの手順で、https://<アプリケーション名>.herokuapp.comをWebhookで使います。)

以下のコマンドが自動的に、リモートリポジトリ”heroku”という名前でgit@heroku.com:<アプリケーション名>.git を登録してくれます。

```bash
$ heroku create <アプリケーション名>  # 例：line-tasks-bot
# 前後の<>はみやすいためであり、実際には書きません

Did you mean create? [y/n]: y
Creating line-tasks-bot... done
https://line-tasks-bot.herokuapp.com/ | https://git.heroku.com/line-tasks-bot.git

# アプリケーション名が既に使われているとエラーになります
```

## **環境変数の設定**

次に、環境変数の設定をします。

Messaging API にアクセスする為に必要な、Channel Secretとアクセストークンを設定します。

いつもの LINE Messaging API のページにあります。

- [Messaging API]
    - [Channel access token]
        - [Channel access token (long-lived)]
        - アクセストークンをコピー
    - [Basic settings]
        - [Channel secret]
        - シークレットキーをコピー

```bash
$ heroku config:set YOUR_CHANNEL_ACCESS_TOKEN="アクセストークンの欄の文字列" --app <アプリケーション名>
$ heroku config:set YOUR_CHANNEL_SECRET="Channel Secretの欄の文字列" --app <アプリケーション名>
```

# **設定ファイル＆Pythonファイルの作成**

ここまでくれば、後は**必要なファイルを作成してgit push（デプロイ）を行えば完成**です。
機能的な面を考慮して、section05 を参考に機能実装をしていきます。

## Go **ファイル（main.py）の作成**

**~~main.py**には、以下のようなコードを記載しました。~~

各ソースコードを見てください。

# **デプロイ**

git コマンドを用いて、今作成したプログラムをデプロイします。

今作成したフォルダに移動し、いつもの git の流れを行ってください。

```bash

$ cd /e/hosei_pc/Heroku/lineBot

// gitの初期ファイルを作成
$ git init

// ローカルリポジトリに結びつくリモートリポジトリを設定
$ heroku git:remote -a <アプリ名>

// 変更したファイルをインデックスに登録
$ git add .

// 変更したファイルをリポジトリに書き込む
$ git commit -m "new commit"

// herokuにローカルで作成したファイルをpush
$ git push heroku main
```

もし、git push でパーミッションエラーになったら、鍵ペアを作り直すことで、解決します。

```bash
$ heroku keys:add
```

# **Webhookの設定**

アプリケーションの更新情報を、他のアプリケーションへリアルタイム提供する仕組みや概念のことです。イベント（リポジトリにプッシュなど）発生時、指定したURLにPOSTリクエストします。

LINE BOTのイベントを、リアルタイムで通知する為に使います。

今回は、LINEにメッセージが届いたときにオウム返しをさせたいので、LINE Developers上で以下の様に設定します。

- [Messaging API]
    - [Webhook settings]
        - [Webhook URL] > [Edit]
        - https://<アプリケーション名>.herokuapp.com/callback
        - [Use webhook] > オンにする

# 動作テスト

LINE Developers 上の QR コードから友達追加してください。

![Untitled (1)](https://user-images.githubusercontent.com/59988255/157020920-d41cad36-13e1-4d9b-9839-c41d7f3d63a0.png)
