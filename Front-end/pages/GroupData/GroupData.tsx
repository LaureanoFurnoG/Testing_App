import React, { useState } from 'react';
import './GroupData.css'
import CollapseCard from '../../components/CollapseTests/Collapse';
import { Button, Input, Space, type GetProps } from 'antd';
import FrontEnd_MetricsCard from '../../components/FrontEnd-Metrics/FrontEnd-Metrics';
import TableUsers_Group from '../../components/TableUsers-Group/TableUsers-Group';
import AddEndpoint from '../../components/AddEndpoint/AddEndpoint';
type SearchProps = GetProps<typeof Input.Search>;
const { Search } = Input;


const GroupData: React.FC = () => {
  const onSearch: SearchProps['onSearch'] = (value, _e, info) => console.log(info?.source, value);

  const [frontUrl, setFrontUrl] = useState("")
  const [url, setUrl] = useState("")
  return (
    <div className='container-group'>
      <h2 className='Group_Name'>G_Name</h2>
      <Space.Compact className='inputEndpoint'>
        <Input placeholder='API URL' />
        <AddEndpoint />
      </Space.Compact>
      <div className='backend-endpoints'>
        <div className='Up-Title-search'>
          <div>
            <h2>Backend endpoints</h2>
            <Button className='btn-hover' type="primary">Run All Endpoints</Button>
          </div>
          <Search placeholder="input search endpoint" onSearch={onSearch} className='searchEndpoint' />
        </div>
        <div>
          <CollapseCard Id={"1"} Name='sd' Type="POST" HTTPResult={201} urlEndpoint="/example/apiendpoint" />
        </div>
      </div>

      <Space.Compact className='containerURLINPUT'>
        <Input className='urlfrontinput' placeholder='URL FRONT END'  onChange={(e) => setUrl(e.target.value)}/>
        <Button type="primary" onClick={() => setFrontUrl(url)}>Submit</Button>
      </Space.Compact>

      <div className='frontend-endpoints'>
        <div className='TopMetrics'>
          <FrontEnd_MetricsCard URL={frontUrl}/>
        </div>
      </div>
      <div className='Users-table'>
        <TableUsers_Group URL='al'/>
      </div>
    </div>
    
  );
};

export default GroupData;
