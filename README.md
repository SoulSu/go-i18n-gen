### i18n gen

生产通用的i18n错误码


#### 安装方法

`go get -v github.com/SoulSu/go-i18n-gen`

#### 使用方法

```bash
//go:generate go-i18n-gen SystemError "10000" "系统错误" "system error"

go generate
```
