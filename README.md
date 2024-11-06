# SSHBannerChanger

SSH更新版本是一个比较折腾且风险很大的事，目前有些漏扫仅通过SSH banner获取版本就开始报一些非常抽象的很极端条件下才能用的漏洞，发现了就要求修，成批成批地告警非常的令人头疼。

因此有这个小工具，使用之后，无论是漏扫还是攻击者，都没法通过banner判断SSH版本，一定程度上也算是安全了。

使用方法：

```bash
# 要求root权限
chmod +x SSHBannerChanger
./SSHBannerChanger
```

使用效果：

```
root@test:/tmp# ./SSHBannerChanger
SSH port: 22, PID: 754, SSHD Locate: /usr/sbin/sshd
Current Banner: OpenSSH_8.9p1 Ubuntu-3ubuntu0.10
Successfully backup sshd executable to: /usr/sbin/sshd.1730859037
Successfully modified the SSHD banner
root@test:/tmp# nc 127.0.0.1 22
SSH-2.0-OpenSSH_fix                     
```

**如果有条件，该升级SSH还是要升级SSH版本的，本工具仅为折中妥协方案。**

请在[Release](https://github.com/Myts2/SSHBannerChanger/releases)下载二进制文件

**注意：如果不信任境外服务，请将代码下载到本地编译，repo仅保证公开代码部分的安全性**
