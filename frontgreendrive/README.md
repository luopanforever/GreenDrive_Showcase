# Getting Started with Create React App
You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).



## OPEN SSL 报错

```json
{
  opensslErrorStack: [ 'error:03000086:digital envelope routines::initialization error' ],
  library: 'digital envelope routines',
  reason: 'unsupported',
  code: 'ERR_OSSL_EVP_UNSUPPORTED'
}
```
`export NODE_OPTIONS=--openssl-legacy-provider`



## 接口一览

### 展示汽车

**`/car/show/:carName/*action`**

示例：

- `/car/show/car2/car2.gltf`

### 获取可用的有效汽车名

**`/car/names/available`**

### 获取汽车名字列表

**`/car/names/list`**

### 批量上传汽车文件

**`/car/upload/:carId`**

示例：

- `/car/upload/car1`

### 删除一辆汽车模型以及器所有相关资源

**`/car/upload/delete/:carName`**

示例：

- `/car/upload/delete/car1`

### 下载特定格式的汽车

**`/car/download/:format/:carName`**

示例：

- `/car/download/stf/car1`

## 接口详细介绍

下载特定格式的汽车

`/car/download/stf/car1`

- 根据倒数第二个参数决定转换哪一个格式

- 如果第二个为`gltf`

  - 直接开始下载文件

- 如果第二个为其他`stl、obj、stp、glb`

  - 则返回如下格式
    ```json
    {
        "code": 200,
        "data": {
            "fileUri": "https://m1.3dwhere.com/download.ashx?c=iy4M2S5902J4c70R&s=ae2e91fd42bb178025f6f07df88d6475"
        },
        "msg": "Download successful"
    }
    ```

    链接为转换后的地址



待处理`bug`

- 删除一辆车时不能自动跳转到已有的车辆
  - 在只有一辆车时删除汽车后不能正确的显示初始页面
