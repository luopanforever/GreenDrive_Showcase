import React, { useEffect } from 'react';
import axios from 'axios';

function App() {
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('/api/test'); // 注意这里是相对路径
        console.log(response.data);
      } catch (error) {
        console.error('There was an error!', error);
      }
    };

    fetchData();
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        {/* 你的组件内容 */}
      </header>
    </div>
  );
}

export default App;
