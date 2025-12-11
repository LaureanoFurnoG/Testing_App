import React from 'react';
import './Homepage.css'
import { useNavigate } from 'react-router-dom';
const Homepage: React.FC = () => {
  let navigate = useNavigate()

  const loginNavigate = () => {
    navigate('/login')
  };
  return (
    <div className='homeContPage'>
      <div className='text-frGreen'>
        <div>
          <h1>Welcome to the Testing App</h1>
          <h2>A testing of backends and frontends</h2>
          <p>To provide testing capabilities for both backends and frontends, all integrated into a single tool designed for simplicity and ease of use.</p>
          <button onClick={() => loginNavigate()}>LET'S GO</button>
        </div>
      </div>
    </div>
  );
};

export default Homepage;
