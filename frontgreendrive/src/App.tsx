import React, { useEffect, useState } from "react"
import {
  Button,
  Upload,
  message,
  Select,
  Space,
  Menu,
  Popconfirm,
  Modal,
  Spin,
} from "antd"
import type { UploadProps } from "antd"
import { UploadOutlined, CloseCircleOutlined } from "@ant-design/icons"
import ShowModel from "./ShowModel"
import { RcFile } from "antd/es/upload"
import Request from "./api"
import { css } from "@emotion/css"
import axios from 'axios';


interface CarList {
  names: string[]
}
interface CarAvailable {
  availableName: string
}

const dowloadBtn = css`
  margin-bottom: 0.5em;
`

const App: React.FC = () => {
  // 所选择的上传文件
  const [selectedUploadFiles, setSelectedUploadFiles] = useState<RcFile[]>([])
  // 有效的汽车列表
  const [carList, setCarList] = useState<string[]>([]) //undefined初始化占个位
  // 选择框选择的汽车
  const [selectedCar, setSelectedCar] = useState<string>("")
  // 汽车有效名
  const [carAvailableName, setCarAvailableName] = useState("")
  // 触发Effect更新
  const [triggerEffect, setTriggerEffect] = useState(false)
  // 控制 Spin 显示
  const [loading, setLoading] = useState(false)

  // 添加新的状态
  const [isDownloadModalVisible, setIsDownloadModalVisible] = useState(false) // 新增状态 - 控制模态框是否可见
  const [downloadUrl, setDownloadUrl] = useState("") // 新增状态 - 存储下载链接

  // useEffect(() => {
  //   // 获取汽车列表
  //   Request.get<CarList>("/names/list").then((res) => {
  //     setCarList(res.data.names)
  //   })

  //   // 获取有效可用汽⻋名字
  //   Request.get<CarAvailable>("/names/available").then((res) => {
  //     setCarAvailableName(res.data.availableName)
  //   })
  // }, [selectedUploadFiles, triggerEffect])

  const fetchData = async () => {
    const carListRes = await Request.get<CarList>("/names/list")
    const carList = carListRes.data.names
    setCarList(carList)

    const carAvailableRes = await Request.get<CarAvailable>("/names/available")
    const availableName = carAvailableRes.data.availableName
    setCarAvailableName(availableName)

    // 如果只有car1表示没有车辆，则展示提示上传界面
    if (availableName === "car1") {
      // 显示提示上传的界面
      setSelectedCar("")
    } else {
      // 根据availableName展示上一个车辆
      const carNumber = parseInt(availableName.replace("car", ""))
      const displayCar = "car" + (carNumber - 1)
      setSelectedCar(displayCar)
    }
  }
  useEffect(() => {
    fetchData()
  }, [triggerEffect, selectedUploadFiles])

  const props: UploadProps = {
    fileList: selectedUploadFiles,
    beforeUpload(file) {
      const isZip = file.type === "application/zip"
      if (!isZip) {
        message.error("只能上传zip文件!")
        return false
      }
      const isLt50M = file.size / 1024 / 1024 < 50
      if (!isLt50M) {
        message.error("文件大小不能超过2MB!")
        return false
      }
      // 选择文件后，先不上传，只添加到 selectedFiles
      setSelectedUploadFiles((currentFiles) => [...currentFiles, file])
      // 返回 false 以阻止自动上传
      return false
    },
    onRemove: (file) => {
      // 移除文件时，从 selectedFiles 中删除
      setSelectedUploadFiles((currentFiles) =>
        currentFiles.filter((f) => f.uid !== file.uid)
      )
    },
  }

  const handleUpload = () => {
    const formData = new FormData()
    selectedUploadFiles.forEach((file) => {
      // formData.append("files[]", file)
      formData.append("file[]", file, file.name)
    })
    // 使用 axios 发送 formData
    Request.post<{ carNames: string[] }>(
      `/upload/${carAvailableName}`,
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      }
    )
      .then((res) => {
        // 处理成功响应
        message.success("上传成功！")
        // 上传成功后清空已选文件列表
        setSelectedUploadFiles([])
        setSelectedCar(res.data.carNames[0])
      })
      .catch((error) => {
        // 处理错误
        message.error("上传失败！")
      })
  }

  // 处理下载的方法
  const handleDownload = async (format: string) => {
    setIsDownloadModalVisible(false) // 关闭模态框
    setLoading(true) // 开始加载
    setLoading(true); // 开始加载

    // 设置最小等待时间为3分钟（180000毫秒）
    const minimumLoadingTime = 180000;
    const startTime = Date.now(); // 记录开始时间
    if (format === "gltf") {
      // 如果是gltf格式，发送请求后直接下载文件
      // 为什么下列提取不了文件？
      // downloadFile(`/download/${format}/${selectedCar}`,`${selectedCar}.zip`)

      // 这个却可以
      try {
        const response = await fetch(`/car/download/${format}/${selectedCar}`, {
          method: "GET",
          headers: {
            "Content-Type": "application/octet-stream", // 指定下载文件的类型
          },
        })

        const blob = await response.blob() // 将文件流转换为 Blob 对象
        const url = URL.createObjectURL(blob) // 创建 Blob 对象的 URL
        // const a = document.createElement('a'); // 创建一个隐藏的<a>元素
        // a.href = url; // 设置<a>元素的链接
        // a.download = `${selectedCar}.zip`; // 设置下载文件的名称
        // document.body.appendChild(a); // 将<a>元素添加到页面中
        // a.click(); // 模拟点击<a>元素进行下载
        // document.body.removeChild(a); // 下载完成后移除<a>元素
        // console.log(url)
        downloadFile(url, `${selectedCar}.zip`)
        URL.revokeObjectURL(url) // 释放 Blob 对象的 URL
      } catch (error) {
        message.error("下载失败，请稍后再试！")
        setLoading(false)
      } finally {
        setLoading(false) // 结束加载
      }
    } else {
      // 如果是其他格式，发送请求获取下载链接后下载文件
      // try {
      //   const response = await Request.get(`/download/${format}/${selectedCar}`, {
      //     timeout: 180000, // 3分钟超时
      //   });

      //   // 判断后端返回的状态码
      //   if (response.code === 200) {
      //     // 正常返回逻辑
      //     const downloadUrl = response.data.fileUri;
      //     downloadFile(downloadUrl, `${selectedCar}.zip`);
      //   } else if (response.code === 400) {
      //     // 超时或其他错误处理逻辑
      //     message.error(response.data.error || "下载失败，请稍后再试！");
      //   }
      //   setLoading(false)
      // } catch (error) {
      //   console.log(error)
      //   message.error("下载失败，请稍后再试！")
      //   // setLoading(false)
      // }
      axios.get(`car/download/${format}/${selectedCar}`).then(response => {
        // 处理响应
        const { code, data, msg } = response.data;
        if (code === 200) {
          // 正常返回逻辑
          const downloadUrl = data.fileUri;
          downloadFile(downloadUrl, `${selectedCar}.zip`);
          setLoading(false)
        } else if (code === 400) {
          // 超时或其他错误处理逻辑
          message.error(msg);
          setLoading(false)
        }
      }).catch(error => {
        // 检查error.response是否存在
        if (error.response) {
          console.log(error.response.status);
          setLoading(false)
          // 处理其他错误信息
        } else {
          // 处理错误对象不存在的情况
          console.log('Error', error.message);
          setLoading(false)
        }
      });
    }
  }





  // 下载文件的函数，用于代替window.open
  const downloadFile = (href: string, filename: string) => {
    // 创建隐藏的可下载链接
    const element = document.createElement("a")
    // element.setAttribute('href', href);
    // element.setAttribute('download', filename);
    element.href = href
    element.download = filename

    // 设置样式以保证它不会显示在页面上
    element.style.display = "none"

    // 将其加入到文档中
    document.body.appendChild(element)

    // 点击链接
    element.click()

    // 移除链接
    document.body.removeChild(element)
  }

  // 选择切换汽车
  const [selectOpen, setSelectOpen] = useState<boolean>(false)
  const handleSelectChange = async (value: string) => {
    setSelectedCar(value)
    setSelectOpen(false)
  }

  // 删除汽车
  const handleConfirmDelete = (carName: string) => {
    Request.delete<any>(`/upload/delete/${carName}`).then((res) => {
      message.success(res.msg)
      setTriggerEffect((prev) => !prev)
    })
  }
  return (
    <>
      {!!carList.length && (
        <Button
          type='primary'
          onClick={() => setIsDownloadModalVisible(true)}
          className={css`
            position: fixed;
            top: 50px;
            left: 0;
            z-index: 9999;
          `}
        >
          下载汽车模型
        </Button>
      )}
      <Modal
        title='下载汽车模型'
        visible={isDownloadModalVisible}
        onOk={() => setIsDownloadModalVisible(false)}
        onCancel={() => setIsDownloadModalVisible(false)}
        footer={null} // 不使用默认底部按钮
        zIndex={99999}
      >
        <div
          className={css`
            display: flex;
            flex-direction: column;
          `}
        >
          {/* 这里是模态框的内容，可以根据需要添加下载选项 */}
          <Button onClick={() => handleDownload("gltf")} className={dowloadBtn}>
            下载 glTF{" "}
            <span
              className={css`
                color: gray;
              `}
            >
              (原格式)
            </span>
          </Button>
          <Button className={dowloadBtn} onClick={() => handleDownload("fbx")}>
            下载 FBX
          </Button>
          <Button className={dowloadBtn} onClick={() => handleDownload("obj")}>
            下载 OBJ
          </Button>
          <Button className={dowloadBtn} onClick={() => handleDownload("glb")}>
            下载 GLB
          </Button>
        </div>
      </Modal>

      {/* <section
        style={{
          display: "flex",
          justifyContent: "space-between",
          marginBottom: "20px",
        }}
      > */}
      <div
        className={css`
          position: fixed;
          top: 90px;
          left: 0;
          z-index: 9999;
        `}
      >
        {!carList.length && (
          <p
            className={css`
              position: fixed;
              top: 0;
              left: 50%;
              transform: translateX(-50%);
            `}
          >
            数据库中还没有车辆请上传
          </p>
        )}
        <Upload {...props}>
          <Button icon={<UploadOutlined />}>上传</Button>
        </Upload>
        <Button
          type='primary'
          onClick={handleUpload}
          disabled={!selectedUploadFiles.length}
          className={css`
            margin-top: 16px;
            visibility: ${selectedUploadFiles.length ? "visible" : "hidden"};
          `}
        >
          开始上传
        </Button>
      </div>

      <Space
        wrap
        className={css`
          position: fixed;
          top: 0;
          right: 0;
          z-index: 9999;
        `}
      >
        {!!carList.length && (
          <Select
            style={{ width: 120 }}
            value={selectedCar}
            open={selectOpen}
            onDropdownVisibleChange={(visible) => setSelectOpen(visible)}
            dropdownRender={() => (
              <Menu>
                {carList.length &&
                  carList.map((item) => (
                    <Menu.Item
                      key={item}
                      onClick={() => handleSelectChange(item)}
                    >
                      <div
                        className={css`
                          display: flex;
                          justify-content: space-between;
                        `}
                      >
                        <span>{item}</span>

                        <Popconfirm
                          // onCancel={cancel}
                          title={`确定删除这个${item}模型?`}
                          okText='确定'
                          cancelText='取消'
                          onConfirm={() => handleConfirmDelete(item)}
                        >
                          <CloseCircleOutlined
                            className={css`
                              &:hover {
                                color: red;
                              }
                            `}
                          />
                        </Popconfirm>
                      </div>
                    </Menu.Item>
                  ))}
              </Menu>
            )}
          />
        )}
      </Space>
      {/* </section> */}

      <section>
        <Spin spinning={loading} tip='下载中...'>
          <ShowModel
            style={{ width: "100vw", height: "100vh" }}
            url={`/car/show/${selectedCar}/${selectedCar}.gltf`}
          />
        </Spin>
      </section>

      {/* <ShowModel style={{width:600,height:600}} url='https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Duck/glTF/Duck.gltf'/> */}
    </>
  )
}
export default App