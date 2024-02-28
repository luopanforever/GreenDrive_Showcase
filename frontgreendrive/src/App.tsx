import React, { useState } from 'react';
import axios from 'axios';

function App() {
  const [ state, setstate ] = useState()
  const fetchData = async () => {
    try {
      const response = await axios.get('/api/test'); // 确保这个URL与你的后端路由匹配
      console.log(response.data); // 这会在控制台中输出后端返回的数据
      setstate(response.data)
    } catch (error) {
      console.error('There was an error!', error);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        {/* 其他内容 */}
        {state}
        <button onClick={fetchData}>请求数据</button>
      </header>
    </div>
  );
}

export default App;
