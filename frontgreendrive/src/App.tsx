import React, { useEffect, useRef } from 'react';
import axios from 'axios';
import * as THREE from 'three';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader';

function App() {
  const mountRef = useRef<HTMLDivElement>(null);

  const loadGLBModel = async (url: string) => {
    const loader = new GLTFLoader();

    const glb = await axios.get(url, { responseType: 'blob' });
    const glbUrl = URL.createObjectURL(glb.data);

    loader.load(
      glbUrl,
      (gltf) => {
        console.log(gltf); // 这里可以访问到GLTF模型的数据

        // 你可以在这里处理gltf.scene或gltf对象，提取需要的数据并以你想要的格式输出或使用
        // 例如，打印出scenes数组
        console.log('Scenes:', gltf.scenes);

        const scene = new THREE.Scene();
        scene.add(gltf.scene);

        const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        camera.position.z = 5;

        const renderer = new THREE.WebGLRenderer();
        renderer.setSize(window.innerWidth, window.innerHeight);
        mountRef.current?.appendChild(renderer.domElement);

        const animate = () => {
          requestAnimationFrame(animate);
          renderer.render(scene, camera);
        };

        animate();
      },
      undefined,
      (error) => {
        console.error(error);
      }
    );
  };

  useEffect(() => {
    return () => {
      if (mountRef.current?.firstChild) {
        mountRef.current.removeChild(mountRef.current.firstChild);
      }
    };
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <div ref={mountRef}></div>
        <button onClick={() => loadGLBModel('/car/car01/car2.gltf')}>加载3D模型</button>
      </header>
    </div>
  );
}

export default App;
