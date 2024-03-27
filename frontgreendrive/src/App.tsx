import React, { useEffect, useState } from "react"
import {
  Button,
  Upload,
  message,
  Select,
  Space,
  Menu,
  Popconfirm,
  FloatButton,
} from "antd"
import type { UploadProps } from "antd"
import { UploadOutlined, CloseCircleOutlined } from "@ant-design/icons"
import ShowModel from "./ShowModel"
import { RcFile } from "antd/es/upload"
import Request from "./api"
import { css } from "@emotion/css"
import { useLocalStorageState } from "ahooks"
import { ModelType } from "./constants/enum"
/* type Response<T extends Record<string, any>> = {
  code: number
  data: T
  msg: string
}
interface CarList extends Response<{ names: string[] }> {}
interface CarAvailable extends Response<{ availableName: string }> {} */

interface CarList {
  names: string[]
}
interface CarAvailable {
  availableName: string
}

const App: React.FC = () => {
  // 所选择的上传文件
  const [selectedUploadFiles, setSelectedUploadFiles] = useState<RcFile[]>([])
  // 有效的汽车列表
  const [carList, setCarList] = useState<string[]>([""]) //undefined初始化占个位
  // 选择框选择的汽车
  const [selectedCar, setSelectedCar] = useLocalStorageState<string>(
    "currentCar",
    {
      defaultValue: carList[0],
    }
  )
  // 汽车有效名
  const [carAvailableName, setCarAvailableName] = useState("")
  // 触发Effect更新
  const [triggerEffect, setTriggerEffect] = useState(false)
  useEffect(() => {
    // 获取汽车列表
    Request.get<CarList>("/names/list").then((res) => {
      setCarList(res.data.names)
    })

    // 获取有效可用汽⻋名字
    Request.get<CarAvailable>("/names/available").then((res) => {
      setCarAvailableName(res.data.availableName)
    })
  }, [selectedUploadFiles, triggerEffect])

  useEffect(() => {
    // 根据 carList 设置 selectedCar 的默认值
    if (carList.length > 0) {
      setSelectedCar(carList[0])
    }
  }, [carList])

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
      setSelectedCar(carList[0])
      setTriggerEffect((prev) => !prev)
    })
  }

  const currentModelType = "gltf"
  const handleDownload = () => {
    console.log(Object.values(ModelType))
    /* if (currentModelType === ModelType.GLTF) {
      console.log('直接下载')
    } else {
      Request.get<any>(`/download/${currentModelType}/${selectedCar}`).then(
        (res) => {
          console.log("isffhasidoIJPO", res.data)
        }
      )
    } */
  }
  return (
    <>
      <section
        style={{
          display: "flex",
          justifyContent: "space-between",
          marginBottom: "20px",
        }}
      >
        <div>
          {!carList.length && <p>数据库中还没有车辆请上传</p>}
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

        <Space wrap>
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
      </section>

      <section>
        <ShowModel
          style={{ width: 600, height: 600 }}
          url={`/car/show/${selectedCar}/${selectedCar}.gltf`}
        />

        <div>
          {/* <FloatButton onClick={() => console.log("onClick")} /> */}
          <FloatButton.Group
            trigger='click'
            type='primary'
            style={{ right: 24 }}
            tooltip={<div>下载不同类型的模型</div>}
          >
            <FloatButton description="GLB" onClick={handleDownload} />
            <FloatButton description="GLB" onClick={handleDownload} />
            <FloatButton description="GLB" onClick={handleDownload} />
          </FloatButton.Group>
        </div>
      </section>

      {/* <ShowModel style={{width:600,height:600}} url='https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Duck/glTF/Duck.gltf'/> */}
    </>
  )
}
export default App
