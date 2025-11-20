import React from 'react';
import './GroupData.css'
import CollapseCard from '../../components/CollapseTests/Collapse';
import { Button, Input, Space, type GetProps } from 'antd';
type SearchProps = GetProps<typeof Input.Search>;
const { Search } = Input;

const Users: React.FC = () => {
  const onSearch: SearchProps['onSearch'] = (value, _e, info) => console.log(info?.source, value);

  return (
    <div className='container-group'>
      <h2 className='Group_Name'>G_Name</h2>
      <Space.Compact className='inputEndpoint'>
        <Input placeholder='API URL' />
        <Button type="primary">Add Endpoint</Button>
      </Space.Compact>
      <div className='backend-endpoints'>
        <div className='Up-Title-search'>
          <div>
            <h2>Backend endpoints</h2>
            <button>Create Group</button>
          </div>
          <Search placeholder="input search endpoint" onSearch={onSearch} className='searchEndpoint' />
        </div>
        <div>
          <CollapseCard Id={1} Name='sd' Endpoints='sdk' Members='293' />
        </div>
      </div>

      <Input placeholder='URL FRONT END' className='urlfrontinput'/>

    </div>
  );
};

export default Users;
