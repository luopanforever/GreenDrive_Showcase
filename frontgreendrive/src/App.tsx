import React, { useEffect, useState } from "react"
import axios from "axios"
import { Button, Upload, message, Select, Space, Menu } from "antd"
import type { UploadProps } from "antd"
import { UploadOutlined } from "@ant-design/icons"
import ShowModel from "./ShowModel"
import { RcFile } from "antd/es/upload"

type Response<T extends Record<string, any>> = {
  code: number
  data: T
  msg: string
}
interface CarList extends Response<{ names: string[] }> {}
interface CarAvailable extends Response<{ availableName: string }> {}

/* interface CarList {
  names: string[]
}
interface CarAvailable {
  availableName: string
} */

const App: React.FC = () => {
  // 所选择的上传文件
  const [selectedUploadFiles, setSelectedUploadFiles] = useState<RcFile[]>([])
  // 有效的汽车列表
  const [carList, setCarList] = useState<string[]>(["undefined"]) //undefined初始化占个位
  // 选择框选择的汽车
  const [selectedCar, setSelectedCar] = useState<string>("car1")

  const [carAvailableName, setCarAvailableName] = useState("")
  useEffect(() => {
    // 获取汽车列表
    axios.get<CarList>("/car/names/list").then((res) => {
      setCarList(res.data.data.names)
    })
    // 获取有效可用汽⻋名字
    axios.get<CarAvailable>("/car/names/available").then((res) => {
      setCarAvailableName(res.data.data.availableName)
    })
  }, [selectedUploadFiles])

  const props: UploadProps = {
    // action: `/car/upload/${carAvailableName}}`,
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
    axios
      .post(`/car/upload/${carAvailableName}`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      })
      .then((response) => {
        // 处理成功响应
        message.success("上传成功！")
        // 上传成功后清空已选文件列表
        setSelectedUploadFiles([])
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
  /*  useEffect(() => {
    axios.get(`/car/show/${selectedCar}/${selectedCar}.gltf`).then((res) => {
      // setShowCar(res.data)
    })
  }, [selectedCar]) */

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
            style={{ marginTop: 16 }}
          >
            开始上传
          </Button>
        </div>

        <Space wrap>
          {carList.length && (
            <Select
              style={{ width: 120 }}
              value={selectedCar}
              open={selectOpen}
              onDropdownVisibleChange={(visible) => setSelectOpen(visible)}
              dropdownRender={() => (
                <Menu>
                  {carList.map((item) => (
                    <Menu.Item
                      key={item}
                      onClick={() => handleSelectChange(item)}
                    >
                      {item}
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
      </section>

      {/* <ShowModel style={{width:600,height:600}} url='https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Duck/glTF/Duck.gltf'/> */}
    </>
  )
}
export default App
