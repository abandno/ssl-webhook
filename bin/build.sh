# windows 交叉编译 Linux 相关参数配置
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64


# 不加文件名, 生成和当前路径一致的名字
go build