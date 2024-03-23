# `GREENDRIVE_SHOWCASE`后端

拉取工具包

```bash
go mod tidy
```

## 怎样运行项目

> 先确保你的电脑上正确安装了`golang`环境

### 阿里云`oss`配置

本项目从`.bash_profile`里面获取以下两个变量

- `OSS_ACCESS_KEY_ID`

- `OSS_ACCESS_KEY_SECRET`

操作命令

```bash
vim ~/.bash_profile

export OSS_ACCESS_KEY_ID=XX
export OSS_ACCESS_KEY_SECRET=XXX

source ~/.bash_profile

# 验证是否配置成功
echo $OSS_ACCESS_KEY_ID
echo $OSS_ACCESS_KEY_SECRET
```

### 关于数据库

本项目使用`mongodb`的`gridfs`来存储`3d`汽车数据，在项目开始不需要配置任何数据库相关的东西，程序自动创建`tdCars`数据库

### 启动命令

```bash
go run main.go route.go
```

