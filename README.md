# wordcloud

```bash
$ podman compose up -d
$ podman compose exec app go run cmd/main.go -w [input word]
...
$ ./output.png
```

## Usage

```bash
# デフォルト（output.png）
$ podman compose exec app go run cmd/main.go -w "テキスト"

# カスタムファイル名を指定
$ podman compose exec app go run cmd/main.go -w "テキスト" -o custom_name.png

# --output でも可能（エイリアス）
$ podman compose exec app go run cmd/main.go -w "テキスト" --output result.png

# ファイル入力と組み合わせ
$ podman compose exec app go run cmd/main.go -f input.txt -o result.png
```

### example

```
$ podman compose exec app go run cmd/main.go -w メロスは激怒した。必ず、かの邪智暴虐の王を除かなければならぬと決意した。メロスには政治がわからぬ。メロスは、 村の牧人である。笛を吹き、羊と遊んで暮して来た。けれども邪悪に対しては、人一倍に敏感であった。きょう未明メロスは村を出発し、野を越え山越え、十里はなれた此このシラクスの市にやって来た。メロスには父も、母も無い。女房も無い。 十六の、内気な妹と二人暮しだ。この妹は、村の或る律気な一牧人を、近々、花婿はなむことして迎える事になっていた。結婚式も間近かなのである。メロスは、それゆえ、花嫁の衣裳やら祝宴の御馳走やらを買いに、はるばる市にやって来たのだ 。先ず、その品々を買い集め、それから都の大路をぶらぶら歩いた。メロスには竹馬の友があった。セリヌンティウスである。今は此のシラクスの市で、石工をしている。その友を、これから訪ねてみるつもりなのだ。久しく逢わなかったのだから 、訪ねて行くのが楽しみである。歩いているうちにメロスは、まちの様子を怪しく思った。ひっそりしている。もう既に日も落ちて、まちの暗いのは当りまえだが、けれども、なんだか、夜のせいばかりでは無く、市全体が、やけに寂しい。のんき なメロスも、だんだん不安になって来た。路で逢った若い衆をつかまえて、何かあったのか、二年まえに此の市に来たときは、夜でも皆が歌をうたって、まちは賑やかであった筈はずだが、と質問した。若い衆は、首を振って答えなかった。しばら く歩いて老爺ろうやに逢い、こんどはもっと、語勢を強くして質問した。老爺は答えなかった。メロスは両手で老爺のからだをゆすぶって質問を重ねた。老爺は、あたりをはばかる低声で、わずか答えた。
```

![sample_1](./sample_output_1.png)
![sample_2](./sample_output_2.png)

## Credits
- x0y0pxFreeFont: http://www17.plala.or.jp/xxxxxxx/00ff/
