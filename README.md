# The Answer Protocol

## English

A multiplayer text adventure in Go. Players explore rooms inspired by Greek mythology, collect items, talk to characters, fight, and complete quests in a shared world.

A server manages game state, while clients provide controls. There is a terminal client (CLI) for typed commands and a graphical client (GUI) with maps and illustrations. You can try both server and client on one computer using separate terminals.

### Features

- Room navigation, items, dialogue, combat, and quests.
- Chat, groups, and connection/disconnection events.
- English and Japanese text, with player state saved and restored as JSON.

### Build and run

Requires Go 1.25 or newer. Run from the project directory.

```sh
make                 # Build the server and CLI
make run-server      # Terminal 1: listen on port 4242
make run-client      # Terminal 2: connect using the CLI
# Set a different server address
make run-client ADDR=127.0.0.1:4242
```

Enter these commands in the CLI. `LANG en` selects English, `CONNECT alice` joins under that name, `LOOK` describes the room, `MOVE east` moves east, `STATUS` shows health, and `QUIT` saves and disconnects.

```text
LANG en
CONNECT alice
LOOK
MOVE east
STATUS
QUIT
```

The GUI uses Fyne. On Linux it also needs a C compiler and OpenGL/X11 development libraries.

```sh
make build-gui
make run-client-gui
```

World data is in `data/world.json`; saved state goes in `saves/`. Saves and assignment materials are excluded from Git.

### Tests

```sh
make test
make lint
```

Checks connection handling, duplicate names, invalid commands, split/coalesced TCP input, room inspection, chat, state restoration, and save failures. Server and CLI builds and tests do not need GUI libraries.

### Implementation

- Uses a goroutine per connection and mutexes to protect shared player state from concurrent updates.
- Queues outgoing messages per client to limit the impact of slow connections.
- The GUI includes maps, dialogue, combat views, and illustrations.

## 日本語

複数のプレイヤーが同じ世界に接続して遊ぶ、Go製のテキストアドベンチャーです。ギリシャ神話を題材にした部屋を移動し、アイテムを集め、登場人物と会話して戦闘やクエストを進めます。

ゲームの状態を管理するサーバーと、操作するためのクライアントを分けています。クライアントは、文字を入力する端末版（CLI）と、地図やイラストを表示する画面版（GUI）の2種類です。1台のPCでも、別のターミナルから接続して試せます。

### 機能

- 部屋の移動、アイテムの取得、NPCとの会話、戦闘、クエスト。
- チャット、グループ、接続・退出などのイベント通知。
- 日本語・英語の切り替えと、JSONによるプレイヤー状態の保存・復元。

### ビルドと実行

Go 1.25以降を使います。プロジェクトのディレクトリで実行してください。

```sh
make                 # サーバーとCLIをビルド
make run-server      # ターミナル1：ポート4242で待ち受け
make run-client      # ターミナル2：CLIで接続
# 接続先を変更する場合
make run-client ADDR=127.0.0.1:4242
```

CLIで次のように入力します。`LANG ja` で日本語を選び、`CONNECT alice` で名前を登録して参加します。`LOOK` で現在の部屋、`STATUS` で体力を確認し、`MOVE east` で東へ進み、`QUIT` で保存して退出します。

```text
LANG ja
CONNECT alice
LOOK
MOVE east
STATUS
QUIT
```

GUIはFyneを使用します。LinuxではCコンパイラとOpenGL・X11の開発ライブラリが必要です。

```sh
make build-gui
make run-client-gui
```

ワールドデータは `data/world.json`、セーブデータは `saves/` に置きます。セーブデータと課題資料はGitの管理対象外です。

### テスト

```sh
make test
make lint
```

基本テストでは接続、名前の重複、不正コマンド、TCPデータの分割・連結、部屋の取得、チャット、保存・復元と保存失敗を確認します。サーバーとCLIのビルド・テストにはGUI用ライブラリは不要です。

### 実装メモ

- 接続ごとの処理をgoroutineで並行に動かし、複数人の共有データはMutexで同時更新を防ぎます。
- 各クライアントへの送信をキューに分け、遅い接続の影響を抑えます。
- GUIにはマップ、会話、戦闘などの表示とイラストを用意しています。
