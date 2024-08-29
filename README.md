# Albatross.swift


# これは何？

2024-08-22 から 2024-08-24 にかけて開催された [iOSDC Japan 2024](https://iosdc.jp/2024/) の中の企画、Swift コードバトルのシステムです。

[サイトはこちら (現在は新規にプレイすることはできません)](https://t.nil.ninja/iosdc-japan/2024/code-battle/)


# サンドボックス化の仕組み

ユーザから任意のコードを受け付ける関係上、何も対策をしないと深刻な脆弱性を抱えてしまいます。

このシステムでは、送信されたコードを [SwiftWasm](https://swiftwasm.org/) によって WebAssembly に変換することでサンドボックス化をおこなっています。


# License

The contents of the repository are licensed under The MIT License, except for

* [backend/admin/assets/css/normalize.css](backend/admin/assets/normalize.css),
* [backend/admin/assets/css/sakura.css](backend/admin/assets/sakura.css) and
* [frontend/public/favicon.svg](frontend/public/favicon.svg).

See [LICENSE](./LICENSE) for copylight notice.
