<div style="position: relative; height: 70vh;">
<h1 style="position:absolute; top: 50%; left: 50%; transform : translate(-50%,-50%);">TinyGoで組み込み開発を始めよう!!</h1>
</div>

====

# おまえ誰よ？

- メカトロソフトエンジニア
- Pythonista & Gopher
- なんでも Go で書いちゃうひと
- Go歴は７年目
- サイト: <http://golang.rdy.jp/>
- 会社: 144Lab
- HN: nobonobo

====

# 本資料のURL
#### *のぼのぼ.GitHub.IO/ぷれぜんてーしょんず/たいにーごー２*
### https://nobonobo.github.io/presantations/tinygo2/

====

# 前回

https://nobonobo.github.io/presantations/tinygo/
# 「TinyGoでIoTを始めよう」

<h2 class="fragment">Let's try embedded development<br/>
with TinyGo!!</h2>

====

# Why TinyGo?

[Go-pain-points](https://github.com/tinygo-org/tinygo/wiki/Go-pain-points)

主な理由を抜粋

- Go標準ライブラリはランタイムと密結合
- Non-OS(ベアメタル)をサポートしない

====

# リンク

- ドキュメント: https://tinygo.org/
- リポジトリ: https://github.com/tinygo-org/tinygo

====

# Important TinyGo members

- Ayke van Laethem(A.K.A. @aykevl)
- Ron Evans(A.K.A. @deadprogram)

====

# TinyGo 0.9.0 released.

<h2 class="fragment">🎉Congratulations!!🎉</h2>

====

# 準備

- docker setup
- docker pull tinygo/tinygo

<b class="fragment">これだけ!</b>

====

# 確認方法
```sh
> docker run -it --rm tinygo/tinygo tinygo version
tinygo version 0.9.0 linux/amd64 (using go version go1.13.1)
```

====

### 最小のサンプル

```go
package main
import (
    "machine"
    "time"
)
func main() {
    LED := machine.Pin(6) // <- ピン番号（ターゲットに合わせる）
    LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
    for {
        LED.High()
        time.Sleep(500 * time.Millisecond)
        LED.Low()
        time.Sleep(500 * time.Millisecond)
    }
}
```

```sh
> docker run -it --rm -v $PWD:/go/src/app -w /go/src/app -e GOPATH=/go \
  tinygo/tinygo tinygo build -target pca10059 -o app.hex .
```

====

# for micro:bit

```go
package main
import (
	"machine"
	"time"
)
func main() {
	ledrow := machine.LED_ROW_1
	ledrow.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledcol := machine.LED_COL_1
	ledcol.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledcol.Low()
	for {
		ledrow.Low()
		time.Sleep(500 * time.Millisecond)
		ledrow.High()
		time.Sleep(500 * time.Millisecond)
	}
}
```

```sh
> docker run -it --rm -v $PWD:/go/src/app -w /go/src/app -e GOPATH=/go \
  tinygo/tinygo tinygo build -target microbit -o app.hex .
```

====

# target指定について

ボードの識別名を指定します

一覧： https://github.com/tinygo-org/tinygo

- circuitplay-express
- microbit
- etc...

同時にこの識別名はビルド時にタグ指定として渡されます

====

# ボード毎の実装

tinygo/src/machine配下にあるファイル群

- 「board_識別名.go」
- 「machine_識別名.go」

これらのうちtargetと合ったファイルだけがビルドに含まれる<br/>
（Go言語の標準仕様）

====

# machineパッケージの役割

ターゲットボード毎に以下の宣言が入ってる

- 主要なピンレイアウトの宣言
- デバッグ出力用UART
- デフォルトI2Cバス
- あればデフォルトSPIバス

====

# for Circuit Playground Express

```go
package main
import (
	"image/color"
	"machine"
	"math"
	"time"
	"tinygo.org/x/drivers/ws2812"
)
func calc(p float64) uint8 {
	i := int(64 * (math.Sin(2*math.Pi*p) + 1.0))
	if i > 255 {
		i = 255
	}
	if i < 0 {
		i = 0
	}
	return uint8(i)
}
func increase(p float64) float64 {
	p += 0.1
	if p > 1 {
		p -= float64(int(p))
	}
	return p
}
func main() {
	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(neo)
	leds := make([]color.RGBA, 10)
	rp, gp, bp := 0.0, 1.0/3, 2.0/3
	for {
		for i := range leds {
			leds[i] = color.RGBA{
				R: calc(rp + 0.1*float64(i)),
				G: calc(gp + 0.1*float64(i)),
				B: calc(bp + 0.1*float64(i))}
		}
		ws.WriteColors(leds)
		time.Sleep(100 * time.Millisecond)
		rp = increase(rp)
		gp = increase(gp)
		bp = increase(bp)
	}
}
```

====

# ライブラリを含むビルド方法

```sh
> go get golang.org/dl/go1.12.12
> go1.12.12 download
...
> GO111MODULE=on go1.12.12 get tinygo.org/x/drivers/ws2812
> GO111MODULE=on go1.12.12 mod vendor
> docker run -it --rm -v $PWD:/go/src/app -w /go/src/app -e GOPATH=/go \
  tinygo/tinygo tinygo build -target circuitplay-express -o app.uf2 .
```

<h2 class="fragment">go1.13のモジュールパスは<br/>厳格にホスト名が必要になった</h2>

====

# TinyGoの出力ファイル名

- 省略不可
- 指定拡張子に合わせてフォーマットが自動選択
- .bin/.hex/.uf2などをサポート

====

# サポート外のボード

- custom.jsonとcustom.ldを用意して
- `-target custom.json`指定にてカスタムボードを扱える
- machineのボード特有定義が存在しないので
- その定義の肩代わりをアプリ作者が代行する

====

# `custom.json`

```json
{
  "inherits": ["cortex-m"],
  "llvm-target": "armv7em-none-eabi",
  "build-tags": ["nrf52840", "nrf", "pca10056"],
  "cflags": [
    "--target=armv7em-none-eabi",
    "-mfloat-abi=soft",
    "-Qunused-arguments",
  ],
  "ldflags": ["-T", "custom.ld"],
  "extra-files": ["lib/nrfx/mdk/system_nrf52840.c", "src/device/nrf/nrf52840.s"]
}
```

====

# `custom.ld`

```ldscript
MEMORY
{
    FLASH_TEXT (rw) : ORIGIN = 0x00000000 + 0x00026000 , LENGTH = 1M - 0x00026000
    RAM (xrw)       : ORIGIN = 0x20000000 + 0x000039c0,  LENGTH = 256K  - 0x000039c0
}
_stack_size = 4K;
INCLUDE "targets/arm.ld"
```

nRF52840の場合スペックシートによると

- FLASHメモリは0x000000000から1MiB
- RAMは0x20000000から256KiB
- 0x00000000 -> 0x26000まではsoftdevice領域
- 0x2000000 -> 0x200039c0まではsoftdeviceが利用

====

# 未サポートボードの対応追加

- custom.jsonである程度動作を確認
- machineパッケージの実装を追加
- [Adding-a-new-board](https://github.com/tinygo-org/tinygo/wiki/Adding-a-new-board)を参照
- Let's コントリビュート！

====

# 組込で必要になる機能

====

# 基本機能の利用

- GPIO
- PWM出力
- ADC/DAC
- UART
- I2C
- SPI
- I2S

<h2 class="fragment">
TinyGoが標準サポート<br/>
ターゲットに関係なく透過的に利用可能</h2>

====

# インラインアセンブラ

```go
  arm.Asm("wfi")
  arm.AsmFull(`
    str {value}, {result}
    `,
    map[string]interface{}{
      "value":  1
      "result": &dest,
    }
  )
```

====

# volatile

```go
import "runtime/volatile"
func foo() {
  var i volatile.Register32
  for{
    i++
  }
}
```

====

# レジスタアクセス(nRF52例)

```go
import　"device/nrf"

func foo() {
  nrf.UART0.PSELTXD.Set(8)
}
```

レジスタの名称はスペックシートで確認

====

# タイマー割り込み(nRF52例)

割り込みハンドラは一式がWeak宣言されていて、<br/>
同名の関数をエクスポートすることで上書きできます。

```go
//go:export TIMER1_IRQHandler
func timerHandler(ptr uint32) {
  // do something
  if nrf.TIMER1.EVENTS_COMPARE[0].Get() != 0 {
    nrf.TIMER1.EVENTS_COMPARE[0].Set(0)
  }
}
```

タイマーの設定はレジスタアクセスにて。

====

# CGOでC資産を利用

Go本家とやり方は同じ
```go
/*
#include "sdk_config.h"
#include "SEGGER_RTT.h"
*/
import "C"
```
clang-cでCコード部分はコンパイルされます。<br/>
そこへ渡したいFLAGSはcustom.jsonに追記します。<br/>
（例えばBLEのSDKヘッダファイルへのインクルードパス追加など）

<p class="fragment">
ホストにあるSDKなどのファイル群は<br/>
dockerでtinygo側にマウントしましょう<br/>
実装例： https://github.com/144lab/tinygo-sample1
</p>

====

# 高機能ハードウェア

- カメラモジュール
- 高機能センサなど
- LCD/OLED/E-Ink
- BLE/Bluetooth
- Ether/Wi-Fi
- USB機器/ホスト
- LoRa/3G/LTE

<h2 class="fragment">
ドライバサポートが必要<br/>
https://github.com/tinygo-org/drivers/ を参照
</h2>

====

# TinyGoがArduinoに統合予定

### [TinyGo on Arduino](https://blog.arduino.cc/2019/08/23/tinygo-on-arduino/)

====

# BLEサポート準備中

プロポーザル
```go
type UUID [4]uint32
type Address [6]uint8
type Bluetooth struct {
    // ...
}
func (b *Bluetooth) Enable(config BluetoothConfig) error {}
func (b *Bluetooth) Disable() error {}
func (b *Bluetooth) Advertise(interval int8, advertisement, scanResponse []byte) {}
type ScanResult struct {
    Address
    // ...
}
func (b *Bluetooth) Scan(callback func(*ScanResult)) error {}
func (b *Bluetooth) StopScan() error {}
```

====

# 余談（時間が許せば）

====

# GoとTinyGoの実験

以下のコードをGoとTinyGoとでビルドして・・・
```go
package main
func recurse(n int) int {
	if n <= 0 {
		return 0
	}
	return n + recurse(n-1)
}
func main() {
	println(recurse(2000000))
}
```

====

# 性能比較

```sh
$ /usr/bin/time -l ./exp-go
2000001000000
        0.42 real         0.17 user         0.04 sys
 132681728  maximum resident set size
```

```sh
$ /usr/bin/time -l ./exp-tinygo
2000001000000
        0.25 real         0.00 user         0.00 sys
    700416  maximum resident set size
```

<h3 class="fragment">TinyGoの方が早い！！＆メモリ使用量少ない！！</h3>

====

# ファイルサイズ

```sh
$ ls -lh
...
-rwxr-xr-x@ 1 nobo  staff   1.1M 10 24 16:59 exp-go
-rwxr-xr-x@ 1 nobo  staff    13K 10 24 16:59 exp-tinygo
...
```

<h3 class="fragment">TinyGoの出力サイズが1.1%！！</h3>

====

# こういうの鵜呑みにしない

- Goは実用にフォーカスしてる
- Goでの再帰呼び出しはメモリを浪費し遅くなる
- LLVMは再帰処理を含む広範囲の最適化を行う
- Goの出力はlibc相当を内包してる
- libcは2MiBくらいある

====

# TinyGoのPros

- Goの最適化とLLVMの最適化の両方が利く
- Cの資産をCGO経由で利用可能
- Goの資産を取り込めるようになる予定
- Goの良さの多くを継承している
- AVR系を除きGCを持っていてメモリ管理が楽
- 組込開発に必要な基本フィーチャーは出そろってきた
- ATSAMD向けUSBCDCサポートが追加

====

# TinyGoのCons

- LLVMバックエンドが重い
- 環境づくりもビルドタイムも時間が必要
- goroutineが本物ではない(LLVMのcoroutine)
- 構造体フィールドタグにアクセスできない<br/>（鋭意対応に向けて活動中ではある）
- 標準jsonエンコーダなどが動かない

====

# まとめ

- 公開からたった一年半で急成長中
- WASMもだいぶ使えるようになってきた
- WebGLサンプルが圧縮で9KBサイズになった事例あり
- 本家の代わりに使う、WASM勢、LLVM勢などが参入する可能性
- RISC-V、ゲームボーイアドバンスの開発も可能になった
- TinyGoがWindowsでも動くようになりつつある

====

# 質問？

====

<div style="position: relative; height: 70vh;">
<h1 style="position:absolute; top: 50%; left: 50%; transform : translate(-50%,-50%);">おわり</h1>
</div>
