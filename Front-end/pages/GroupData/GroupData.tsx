import React from 'react';
import './GroupData.css'

import { Button, Input, Space } from 'antd';

const Users: React.FC = () => {
  return (
    <div className='container-group'>
      <h2 className='Group_Name'>G_Name</h2>
      <Space.Compact className='inputEndpoint'>
        <Input placeholder='API URL' />
        <Button type="primary">Add Endpoint</Button>
      </Space.Compact>
      <div className='backend-endpoints'>
        
      </div>
    </div>
  );
};

export default Users;
