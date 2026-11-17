<div align="center">

# astrocyte-s4-public-context-20261015

A small Unicode text counter in Go — words, code points, lines — plus the stdin command that prints them as JSON.

一个 Go 写的小文本统计器：词数、码点数、行数，外加一个把结果打成 JSON 的 stdin 命令。

![Go](https://img.shields.io/badge/go-1.27.2-00ADD8)
![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)
![Output](https://img.shields.io/badge/output-JSON-blue)
![License](https://img.shields.io/badge/license-not%20specified-lightgrey)

</div>

> **How to read this page / 怎么读这页**
> English first, then the same content in Chinese — not a summary, not a separate half.
> 每段英文在前、同一内容的中文紧跟在后；不是摘要，也没有单独的「中文版」。

---

## At a glance · 一眼看懂

| Field · 字段 | Value · 内容 |
|---|---|
| **Problem<br/>问题** | `wc -w` counts bytes-ish tokens and `wc -m` counts weirdly across locales; if you want code points, words by Unicode whitespace, and a line count that ignores a trailing newline, you end up hand-rolling it.<br/>`wc -w` 的分词口径粗糙，`wc -m` 在不同 locale 下结果也不一致；想要「按码点数字符、按 Unicode 空白分词、行数忽略结尾换行」就得自己手写。 |
| **Fix<br/>做法** | One `Count(text) Stats` function, one `CountAll([]string) Stats` for independent inputs, and a stdin command that JSON-encodes the result.<br/>一个 `Count(text) Stats`，一个给多段独立输入用的 `CountAll([]string) Stats`，再加一个把结果 JSON 化的 stdin 命令。 |
| **Cost<br/>代价** | Zero dependencies, zero network calls, Go 1.27.2, four source files.<br/>零依赖、零网络请求，Go 1.27.2，四个源码文件。 |
| **Not<br/>不是** | Not UTF-8 byte counting, not grapheme-cluster counting, no file arguments, no `--help` flags.<br/>不是按 UTF-8 字节数，也不是按字形簇数；不接受文件参数，没有 `--help` 之类的开关。 |

## What it does · 它做什么

Everything exported, from [textstats/stats.go](textstats/stats.go):

全部导出内容，来自 [textstats/stats.go](textstats/stats.go)：

| Symbol · 符号 | Signature · 签名 | Behaviour · 行为 |
|---|---|---|
| `Stats` | `struct { Words, Characters, Lines int }` | Plain result struct, no methods.<br/>纯结果结构体，没有方法。 |
| `Count` | `Count(text string) Stats` | Counts one string.<br/>统计单个字符串。 |
| `CountAll` | `CountAll(texts []string) Stats` | Sums `Count` over each element independently.<br/>对每个元素分别 `Count` 后求和。 |

The command in [cmd/textstats/main.go](cmd/textstats/main.go) reads all of stdin and writes one JSON object to stdout, or the error to stderr with exit code 1. It calls `Count` only — `CountAll` has no command-line surface.

[cmd/textstats/main.go](cmd/textstats/main.go) 里的命令把 stdin 全部读完，往 stdout 写一个 JSON 对象；出错则写 stderr 并以退出码 1 结束。它只调用 `Count`——`CountAll` 没有对应的命令行入口。

## Counting rules · 计数规则

Four rules, each one visible in the source and pinned by a test:

四条规则，每条都能在源码里看到，也都有测试钉住：

| Rule · 规则 | Implementation · 实现 | Consequence · 后果 |
|---|---|---|
| **Characters are code points**<br/>字符按码点 | `utf8.RuneCountInString` | `你好 world\n` is 9 characters, though it is 13 bytes of UTF-8; an emoji built from several runes counts several times.<br/>`你好 world\n` 算 9 个字符，但它是 13 个 UTF-8 字节；由多个 rune 组成的 emoji 会被算成多个字符。 |
| **Words are Unicode whitespace separated**<br/>词按 Unicode 空白分 | `len(strings.Fields(text))` | Any run of whitespace, CJK text included, splits words the same way Go's `Fields` does.<br/>任意一段空白都分词，中文文本跟 Go 的 `Fields` 口径一致。 |
| **A trailing newline is not a line**<br/>结尾换行不算一行 | count of `\n`, plus 1 if the text is non-empty and does not end with `\n` | `"a\n"` is 1 line, `"a"` is 1 line, `""` is 0, `"a\nb"` is 2.<br/>`"a\n"` 算 1 行，`"a"` 算 1 行，`""` 算 0 行，`"a\nb"` 算 2 行。 |
| **`CountAll` is per-input, not concatenation**<br/>`CountAll` 按段计，不是拼接 | sums each element's own `Stats` | `CountAll([]string{"one two", "three"})` gives 3 words and 2 lines, while `Count("one twothree")` would give 2 words and 1 line. The two are not interchangeable.<br/>`CountAll([]string{"one two", "three"})` 得 3 个词、2 行；而 `Count("one twothree")` 得 2 个词、1 行。两者不能互换。 |
| **JSON keys are Go field names**<br/>JSON 键名就是 Go 字段名 | `json.NewEncoder` over `Stats`, which has no struct tags | The object is shaped `{"Words":…,"Characters":…,"Lines":…}` with capitalised keys.<br/>对象形状是 `{"Words":…,"Characters":…,"Lines":…}`，键名首字母大写。 |

## Quickstart · 快速开始

### Requirements · 环境要求

Go 1.27.2 or newer, per [go.mod](go.mod). No third-party modules, so there is nothing to download.

按 [go.mod](go.mod) 需要 Go 1.27.2 或更新。没有第三方模块，因此没有东西要下载。

### 1. Count something · 数点东西

```bash
git clone https://github.com/songconmaisaix31-design/astrocyte-s4-public-context-20261015.git
cd astrocyte-s4-public-context-20261015
printf '你好 world\n' | go run ./cmd/textstats   # stdin in, one JSON object out · 读 stdin，输出一个 JSON 对象
```

The command takes no arguments; every byte it sees on stdin is counted, including a trailing newline. The example input is the same string the shipped test asserts on, so the expected numbers are 2 words, 9 characters, 1 line.

这个命令不接受任何参数；它看到的 stdin 每个字节都参与统计，包括结尾换行。示例输入就是自带测试里断言的那个字符串，所以期望值是 2 个词、9 个字符、1 行。

### 2. Validate · 验证

```bash
go test ./...     # the shipped test · 跑自带测试
go build ./...    # compiles package and command · 编译包与命令
```

These two commands are the ones the repository documents for validation. They were **not executed while writing this page** — the review was read-only, with no build. The commands above are quoted, not replayed.

这两条命令是仓库自己文档里给出的验证方式。撰写本页时**没有执行**它们——这次是只读审查，没有构建。上面的命令是引用，不是重放。

### 3. Use it as a package · 当包用

The module path is the repository URL, so the package is importable as `github.com/songconmaisaix31-design/astrocyte-s4-public-context-20261015/textstats`. Resolution through the Go module proxy was not verified.

模块路径就是仓库 URL，所以包可以按 `github.com/songconmaisaix31-design/astrocyte-s4-public-context-20261015/textstats` 导入。是否已能通过 Go module proxy 解析，未验证。

```go
import "github.com/songconmaisaix31-design/astrocyte-s4-public-context-20261015/textstats"

s := textstats.Count("你好 world\n")   // s.Words=2, s.Characters=9, s.Lines=1
```

## Tests as shipped · 自带的测试

One test function in [textstats/stats_test.go](textstats/stats_test.go), stdlib `testing` only, no CI workflow wired to it.

[textstats/stats_test.go](textstats/stats_test.go) 里有一个测试函数，只用标准库 `testing`，没有 CI 工作流在跑它。

| Assertion · 断言 | Expected · 期望 |
|---|---|
| `Count("你好 world\n")` | `Words=2, Characters=9, Lines=1` |
| `Count("")` | zero value `Stats{}` |
| `CountAll([]string{"one two", "three"})` | `Words=3, Lines=2` |

Three assertions in one function is the whole suite. There is no test for a large input, for CRLF line endings, or for the JSON command.

一个函数里三条断言，就是全部测试。大输入、CRLF 换行，以及那个 JSON 命令，都没有测试。

## Why this repository exists · 这个仓库为什么存在

Per the repository's own description, this is a public project owned by the Astrocyte development agents, used to exercise repository indexing, source version tracking, and Git worktree isolation. It contains no user source and no private data.

按仓库自己的说明，这是一个由 Astrocyte 开发 agent 所有的公开项目，用来验证仓库索引、源码版本追踪和 Git worktree 隔离。它不含用户源码，也不含私有数据。

Read that as the primary purpose: the text counter is deliberately small and complete so the thing under test is the pipeline around it, not the algorithm. If you are looking for a production text-stats library, this is a reference implementation of the counting contract, not a maintained one.

请把这句话当作主要用途：文本统计器是刻意做小而完整的，被测对象是它外围的流水线，不是算法。如果你在找一个生产用的文本统计库，这里是计数契约的参考实现，不是被长期维护的那个。

## Status and limits · 状态与边界

One commit (`4317d61`, 2026-10-15), five files, two packages, no third-party dependencies, no CI, no `LICENSE` file.

一次提交（`4317d61`，2026-10-15）、五个文件、两个包，无第三方依赖，无 CI，无 `LICENSE` 文件。

Known limits, stated rather than hidden:

边界写出来，而不是藏起来：

- **No CI.** There is no workflow badge on this page because there is no `.github/workflows` directory to point at.<br/>**没有 CI。** 本页没有工作流徽章，因为仓库里不存在 `.github/workflows` 目录可指。
- **Characters are runes, not grapheme clusters.** Combining marks and ZWJ emoji each count separately, so the number will not match what a human calls "characters" for such strings.<br/>**字符是 rune，不是字形簇。** 组合附加符号和 ZWJ emoji 会分别计数，所以对这类字符串，结果与人眼认定的「字符数」不一致。
- **Line counting has no CRLF handling.** `\r` is counted as a character and only `\n` starts a new line; a CRLF file gets one extra character per line.<br/>**行统计不处理 CRLF。** `\r` 会被当成一个字符，只有 `\n` 才换行；CRLF 文件每行都会多算一个字符。
- **No file arguments.** `cmd/textstats` only reads stdin, so counting files needs a shell redirect.<br/>**不接受文件参数。** `cmd/textstats` 只读 stdin，统计文件要自己用 shell 重定向。
- **No benchmarks or fuzzing.** No performance numbers are claimed on this page because none exist in the repository.<br/>**没有基准测试或模糊测试。** 本页不声明任何性能数字，因为仓库里没有这类数据。
- **No license.** No `LICENSE` file means no license is granted by default.<br/>**没有许可。** 没有 `LICENSE` 文件意味着默认不授予任何许可。

## License · 许可

No `LICENSE` file is present in this repository, so no license is granted by default and reuse needs permission from the author.

仓库里没有 `LICENSE` 文件，因此默认不授予任何许可，复用需要作者同意。
