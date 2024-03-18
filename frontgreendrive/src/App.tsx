import React, { useEffect, useState } from "react"
import axios from "axios"
import { Button, Upload, message, Select, Space, Menu } from "antd"
import type { UploadProps } from "antd"
import { UploadOutlined } from "@ant-design/icons"
import ShowModel from "./ShowModel"

type Response<T extends Record<string, any>> = {
  code: number
  data: T
  msg: string
}
interface CarList extends Response<{ names: string[] }> {}
interface CarAvaliable extends Response<{ availableName: string }> {}

const App: React.FC = () => {
  const [showCar, setShowCar] = useState("/car/show/car1/car1.gltf")
  // const [uploading, setUploading] = useState(false)
  const [carList, setCarList] = useState<string[]>(["undefined"]) //undefined初始化占个位
  const [selectedCar, setSelectedCar] = useState<string>("car1")

  useEffect(() => {
    // 获取汽车列表
    axios.get<CarList>("/car/names/list").then((res) => {
      setCarList(res.data.data.names)
    })
    // 获取有效可用汽⻋名字
    axios.get<CarAvaliable>("/car/names/available").then((res) => {
      console.log("有效汽车名", res.data.data.availableName)
    })
  }, [])

  useEffect(() => {
    axios.get("car/show/car1/car1.gltf").then((res) => {
      console.log(res.data)
    })
  }, [])

  // todo 上传
  const props: UploadProps = {
    // action: 'http://localhost:8080/car/upload/car2',
    beforeUpload(file) {
      const isZip = file.type === "application/zip"
      if (!isZip) {
        message.error("只能上传zip文件!")
      }
      const isLt50M = file.size / 1024 / 1024 < 50
      if (!isLt50M) {
        message.error("文件大小不能超过2MB!")
      }
      return isZip && isLt50M
    },
    onChange({ file, fileList }) {
      // 不在上传中
      if (file.status !== "uploading") {
        console.log(file, fileList)
      }
    },
    customRequest({ file, onSuccess, onError, onProgress }) {
      // onSuccess()
    },
  }

  // 选择切换汽车
  const [selectOpen, setSelectOpen] = useState<boolean>(false)
  const handleSelectChange = (value: string) => {
    // todo请求切换
    axios.get("http://localhost:8080/car/show/car1/car1.gltf")
    setSelectedCar(value)
    setSelectOpen(false)
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
        </div>

        <Space wrap>
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
        </Space>
      </section>

      <section>
        <ShowModel style={{ width: 600, height: 600 }} url={showCar} />
      </section>

      {/* <ShowModel style={{width:600,height:600}} url='https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Duck/glTF/Duck.gltf'/> */}
    </>
  )
}
export default App
