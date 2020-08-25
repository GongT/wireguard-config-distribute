Wireguard VPN配置文件分发工具

# 原理
Wireguard是一个P2P VPN，与传统VPN最大的区别就是任意两个节点都是直接通信的，不需要服务器。    
不过目前Wireguard配置还是比较麻烦的，如果节点多了，每个节点都要写其他所有节点的信息，特别难修改。    
所以就做了这么个玩意。

不过这个问题估计很快wireguard官方就会解决，所以这里还做了一些官方不可能提供的功能。

这是一个**服务器+客户端**模式的自动配置工具，附带一些其他功能。其中**必须**有一个**公网可访问**的“配置服务器”。    
公网仅能通过IPv6地址访问的情况也可以，而且IP不需要固定。    
这个服务器不需要加入wireguard网络，而是只提供配置功能（也可以选择加入）。

**主要功能有：**
* 支持linux、mac、openwrt、windows（win仅有客户端）
	* windows以服务形式运行
	* 示例systemd（linux）和procd（openwrt）脚本
* 配置Wireguard节点：
    * 服务器提供wireguard网络的基本信息，比如MTU、ip范围、掩码之类的
	* 每个节点提供自己的信息
* 通过*UPnP、*NAT-PMP映射入站端口
* 简单的流量混淆
* 识别相同内网的其他节点，直连ip并禁用混淆
* 同时配置多个无关的VPN网络
* 通过修改hosts文件实现节点之间的名称解析
* *TODO：支持移动设备（或ip偶尔会变的情况）*
* *TODO：自动配置路由转发（使某个节点成为传统VPN中的服务器，进而让两个内网节点可以通信）*
* *TODO：多用户*
* *TODO：高级混淆*
* *TODO：UDP打洞*

# 用法
## 安装
首先需要手动安装wireguard：https://www.wireguard.com/install/

目前不支持下载可执行文件，需要参考下面`开发`部分

## 运行
完整参数列表：
* Windows: `wireguard-config-client.exe /help`
* Others: `wireguard-config-client --help`

所有参数都可以通过环境变量传入，推荐这么做，详见help内容。    
程序运行需要管理员权限（或linux的CAP_NET_ADMIN），Windows上如果发现没有会自动提权。

### 服务端
服务器程序只支持linux（包括wrt）直接启动即可。   
普通linux推荐直接使用`./scripts/services/server.service`文件解决（但是必须手动修改里面的`Environment=`们传入正确的参数）

服务器有2种工作模式：
1. 明文模式：`--insecure`参数，仅用于有nginx之类的工具处理SSL的情况，当监听unix socket时只能用这个模式
1. 独立服务器模式：可以自签名（默认），或载入现有SSL证书

如果使用提供的service文件，则服务器会在`/var/lib/wireguard-config-server`文件夹中存放数据。其他方法调用时一般会存放在`~/.wireguard-config-server`。
* `vpns.json`: 可以定义多个VPN网络，设置各个网络的基本信息（比如MTU）
* `password.txt`：客户端连接时必须使用的口令

### 客户端
```powershell
wireguard-config-client --server=x.x.x.x:1234 --其他参数
```

linux自启动问题可以直接参考`./scripts/services`里的几个文件。    

### 客户端（Windows）
Windows既可以用`/参数 值`也可以用`--参数=值`    
除其他参数外，再加上`/install`参数，程序本身会复制到`C:\Program Files\WireGuard`，并把当前参数注册成开机启动的服务。用`/uninstall`停止并卸载服务。
```powershell
wireguard-config-client.exe /install /server x.x.x.x:1234 /其他参数
Start-Service wg_default ## 立即启动服务
```

# 开发
## 准备
1. 需要比较新的golang、powershell
1. 在项目外的某个目录运行：    
    ```bash
    go get -u github.com/GongT/go-generate-struct-interface/cmd/go-generate-struct-interface
	```
1. 运行： 
    ```bash
	./scripts/run.ps1 [client|server|tool] (...arguments)
	```
