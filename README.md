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

### Command reference

Type one command per line; no `/` prefix is needed. Command keywords are case-insensitive. Choose the language before `CONNECT`; gameplay commands require a connection.

| Command | Action |
| --- | --- |
| `LANG <language>` | Select `en` or `ja` before joining |
| `CONNECT <name>` | Join with a unique player name containing no spaces |
| `LOOK` | Show the room, available exits, players, items, and NPCs |
| `MOVE <direction>` | Follow an exit shown by `LOOK` |
| `WHO` | Show the number of connected players; also available before joining |
| `CHAT <scope> <message>` | Send a message to `GLOBAL`, `ROOM`, or your `GROUP` |
| `TAKE <item>` | Pick up an available item in the room |
| `DROP <item>` | Drop an item from your inventory |
| `INVENTORY` | Show carried items |
| `TALK <npc>` | Talk to an NPC in the room |
| `QUEST <npc>` | Accept a quest or review its details with its quest giver |
| `QUESTS` | Show your quest progress |
| `ATTACK <npc>` | Attack an NPC in the room |
| `DEFEND` | Brace during combat |
| `FLEE` | Attempt to escape the current encounter |
| `STATUS` | Show your current health and status |
| `STATE` | Show crew, online player names, group membership, and invitations as JSON |
| `GROUP CREATE` | Create a group |
| `GROUP INVITE <name>` | Invite an online player to your group |
| `GROUP JOIN <leader>` | Accept an invitation using the group leader's player name |
| `GROUP LEAVE` | Leave your group |
| `QUIT` | Save and disconnect; also available before joining |

For `<item>` and `<npc>`, use an ID or the displayed name in the selected language. Displayed names and chat messages may contain spaces; do not add quotes. Use the player's `CONNECT` name for group invitations, and the leader's name rather than the group ID for `GROUP JOIN`.

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

Checks connection handling, duplicate names, invalid commands, split/coalesced TCP input, room inspection, chat, state restoration, and save failures. GUI tests use Fyne's headless `ci` driver to verify protocol messages, JSON responses, world data, embedded images, and status-bar rendering without X11 or OpenGL. Opening and interacting with a real desktop window is a separate manual check.

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

### コマンド一覧

1行に1コマンドを入力します。先頭の `/` は不要で、コマンド名の大文字・小文字は区別しません。`LANG` は `CONNECT` の前に指定し、ゲーム内の操作は参加後に実行します。

| コマンド | 動作 |
| --- | --- |
| `LANG <language>` | 参加前に `en` または `ja` を選ぶ |
| `CONNECT <name>` | 空白を含まない、他のプレイヤーと重複しない名前で参加する |
| `LOOK` | 部屋・移動可能な出口・プレイヤー・アイテム・NPCを確認する |
| `MOVE <direction>` | `LOOK` に表示された出口へ移動する |
| `WHO` | 接続中のプレイヤー数を確認する。参加前も使用可能 |
| `CHAT <scope> <message>` | `GLOBAL`・`ROOM`・`GROUP` のいずれかへ発言する |
| `TAKE <item>` | 部屋にある取得可能なアイテムを拾う |
| `DROP <item>` | 持っているアイテムを置く |
| `INVENTORY` | 持ち物を確認する |
| `TALK <npc>` | 同じ部屋のNPCと話す |
| `QUEST <npc>` | 依頼主からクエストを受注する、または内容を確認する |
| `QUESTS` | 自分のクエストの進捗を確認する |
| `ATTACK <npc>` | 同じ部屋のNPCを攻撃する |
| `DEFEND` | 戦闘中に防御する |
| `FLEE` | 現在の戦闘からの離脱を試みる |
| `STATUS` | 体力などの現在の状態を確認する |
| `STATE` | 船員数・接続中の名前・所属グループ・招待をJSONで確認する |
| `GROUP CREATE` | グループを作成する |
| `GROUP INVITE <name>` | 接続中のプレイヤーを自分のグループへ招待する |
| `GROUP JOIN <leader>` | リーダーのプレイヤー名を指定して招待を受ける |
| `GROUP LEAVE` | グループを抜ける |
| `QUIT` | 保存して切断する。参加前も使用可能 |

`<item>`・`<npc>` はID、または選択した言語の表示名を指定します。表示名やチャット本文に空白がある場合も引用符は付けません。招待には相手の `CONNECT` 名、`GROUP JOIN` にはグループIDではなくリーダーの名前を使います。

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

接続、名前の重複、不正コマンド、TCPデータの分割・連結、部屋の取得、チャット、保存・復元と保存失敗を確認します。GUIはFyneの `ci` ドライバを使い、X11やOpenGLなしで通信・JSON応答・世界データ・埋込画像・ステータスバーを確認します。実際のデスクトップでのウィンドウ操作は別途確認します。

### 実装メモ

- 接続ごとの処理をgoroutineで並行に動かし、複数人の共有データはMutexで同時更新を防ぎます。
- 各クライアントへの送信をキューに分け、遅い接続の影響を抑えます。
- GUIにはマップ、会話、戦闘などの表示とイラストを用意しています。
