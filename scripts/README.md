# 安装、运行脚本

由于我的环境有限，大部分脚本没有考虑更多情况，所以这些脚本更多的是示例性质。

主要有用的是两个目录：
 * realtest：我自己用的脚本，可以参考 *（其中的windows.ps1和linux.sh、openwrt-remote.sh）*
 * services：比较有可移植性的脚本

## Windows.ps1
这个脚本会做下面这些事：
* 设置一个代理服务器（你需要先删除这个）
* 从github下载client.exe，和[winsw](https://github.com/winsw/winsw)，到`C:\Program Files\WireGuard`
* 根据环境变量（OneDriveConsumer、OneDrive）找到你的OneDrive路径，读取`{OneDrive}/Software/WireguardConfig/{ComputerName}.conf`这个文件
* 把client.exe和配置文件的**内容**注册到windows服务里（所以配置文件改变后需要重新注册）
* 把一个自动更新脚本放在`C:\Program Files\WireGuard`，并注册到Windows任务计划程序里
* 启动服务
