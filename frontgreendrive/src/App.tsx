import React from "react"
import ShowModel from "./ShowModel"
const App: React.FC = () => {
  return (
    <>
      <ShowModel
        style={{ width: 600, height: 600 }}
        url='http://localhost:8080/car/show/car1/car1.gltf'
      />
      {/* <ShowModel style={{width:600,height:600}} url='https://raw.githubusercontent.com/KhronosGroup/glTF-Sample-Models/master/2.0/Duck/glTF/Duck.gltf'/> */}
    </>
  )
}
export default App
