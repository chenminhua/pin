### pin 是啥

假设你在两个不同的局域网中各有一台机器，machineA 在 networkA中，machineB 在 networkB中。如何将machineA上的一个文件传递给machineB呢？

我希望能够有一种方法，就像复制黏贴一样自然地将文件从machineA传递到machineB，而pin就是一个为此而生的工具。pin 是一个运行在服务端的剪切板，当你需要在某个机器上传一个文件到另一个文件上的时候。

pin的原理非常简单，当在客户端执行pin copy操作时，pin会将被copy的文件或者字符串，或者任何可以变成二进制流的东西传递给pin server。
pin server会开辟一块内存专门存储这部分数据。
当客户端需要paste文件的时候，直接运行pin命令，就可以从pin server上将那块内存的数据拉取到本地。

### pin pipe

普通模式下，pin会将文件全部上传给服务端，并放在服务端内存中。
但是如果文件太大，这样做会打爆服务端内存，此时需要使用pipe模式。

原理也很简单，pipe模式下，server会和sender以及receiver分别维护一个连接。当从sender上读取到数据，就写给receiver。



### quick start



```sh
## in your server,run
pin --server

## in your client which you want to copy the file, run
pin --copy < file
## 推荐设置 alias pinc="pin --copy <"

## in your client which you want to paste the file, run
pin > file
## 推荐设置 alias pinp="pin >"
```

pin --server
pin file
pin --list
pin

http api

cli api

web interface
