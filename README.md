# SSL 证书自动部署 webhook web服务

> 暂仅仅支持 OHTTPS 平台的自动部署。https://ohttps.com/

1. 构建
```shell
cd {项目根目录}
sh bin\build.sh 
```

2. 上传

生成 ssl-webhook 可执行文件。
上传至服务器。
添加执行权限
```shell
chmod 777 ssl-webhook
```

3. 启动

启动服务（后台）
```shell
./ssl-webhook -d=true
```

4. 部署证书

证书管理》点选要部署的证书》部署证书

结果：
会在 nginx 证书目录下（"/etc/nginx/cert"），生成对应域名的目录（域名即使目录名）。
内含两个文件：
```
cert.key // 私钥
fullchain.cer // 证书
```

改域名下原有证书目录，重命名，添加当时时间后缀备份。

![image.png](https://s2.loli.net/2022/11/19/5kI1PARGSaHs32j.png)