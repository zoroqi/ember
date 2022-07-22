# ember

主要用来解析 kindle 的 "My Clippings.txt" 文件的工具.

"My Clippings.txt" 文件格式

```
书名 (作者)
标注索引和方式 | 标注时间

内容
======
```

解析结果会生成两部分
```
// 具体解析后的标注
type Clipping struct {
    Book          string        // 书名
    Text          string        // 高亮或笔记
    Time          time.Time     // 时间
    LocationStart int           // 标注的起始坐标, 当坐标只有一个值得时候 0
    LocationEnd   int           // 标注的结束坐标,
    Type          ClippingType  // 枚举: highlight,note,bookmark
}
// 未解析的块结构
type Block struct {
    num   int
    Book  string
    Index string
    Text  string
    Err   error
}
```

提供两个工具
1. cli, 直接以 json 输出
    1. 可以使用 jq 进行过滤或定制化输出
    2. 输出所有笔记内容 `./cli -f clipping.txt | jq '.[] |  select(.Type == "note") | {"Book":.Book,"Text":.Text}'`
2. specialcli, 一个简单的交互, 提供我自己的特殊格式输出

示例:
```go
clippingText, err := os.Open(*clippingFile)
if err != nil {
    fmt.Println(err)
    return
}
clips, errBlock := ember.ParseClippings(clippingText)
for _, clip := range clips {
    fmt.Println(clip)
}
```

具体的解析方式和系统语言, 书籍格式有关; 时间都需要进行特殊处理. 我只有中文和英文两种数据格式, 就只写了两种, 提供样例数据可以提供追加解析
