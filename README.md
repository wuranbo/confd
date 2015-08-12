# Usage

> git clone https://github.com/wuranbo/confd.git
> cd confd
> ./build
> cp bin/confd ~/bin/

# wuranbo-v0.8.2

- 添加了新的backend类型`json`，用于通过一定约定的json文件来定义key
- 约定的json文件格式样例请看[integration/json/test.sh](integration/json/test.sh)

[下载](docs/installation.md)

# wuranbo-v0.8.1

- 添加了递归查找功能: 注意toml文件的src是相对config根目录的相对路径，有子目录的也要对应填上。
- 添加template函数concat，bytetoM, 四则运算
- 增加了只在内存内处理模板


# confd

[![Build Status](https://travis-ci.org/wuranbo/confd.png?branch=master)](https://travis-ci.org/wuranbo/confd)

`confd` is a lightweight configuration management tool focused on:

* keeping local configuration files up-to-date using data stored in [etcd](https://github.com/coreos/etcd),
  [consul](http://consul.io), or env vars and processing [template resources](docs/template-resources.md).
* reloading applications to pick up new config file changes

## Community

* IRC: `#confd` on Freenode
* Mailing list: [Google Groups](https://groups.google.com/forum/#!forum/confd-users)
* Website: [www.confd.io](http://www.confd.io)

## Getting Started

Before we begin be sure to [download and install confd](docs/installation.md).

* [quick start guide](docs/quick-start-guide.md)

## Next steps

Check out the [docs directory](docs) for more docs.
