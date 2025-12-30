import React, { useEffect, useState } from 'react';
import './GroupData.css'
import CollapseCard from '../../components/CollapseTests/Collapse';
import { Button, Form, Input, Space, type GetProps } from 'antd';
import FrontEnd_MetricsCard from '../../components/FrontEnd-Metrics/FrontEnd-Metrics';
import TableUsers_Group from '../../components/TableUsers-Group/TableUsers-Group';
import AddEndpoint from '../../components/AddEndpoint/AddEndpoint';
type SearchProps = GetProps<typeof Input.Search>;
const { Search } = Input;
import { useParams } from "react-router-dom";
import axiosInstance from '../../axios';

interface BackendTest {
  ID: number;
  Idgroup: number;
  Name: string;
  Httptype: string;
  Requesttype: string;
  Urlapi: string;
  ResponseHttpCode: number;
  Token: string;
  Response: JSON | null;

  Group: {
    ID: number;
    KeycloakID: string;
  };

  Header: Record<string, string>;

  Request: Record<string, any>;
}


const GroupData: React.FC = () => {
  const onSearch: SearchProps['onSearch'] = (value, _e, info) => setSearchEndpoint(value);
  const [search, setSearchEndpoint] = useState("")
  const { groupId } = useParams(); 
  const [frontUrl, setFrontUrl] = useState("")
  const [url, setUrl] = useState("")
  const [tests, setTest] = useState<BackendTest[]>([])
  const getTests = async () => {
    try{
      const response = await axiosInstance.get(`/api/tests/${groupId}/find-tests`, {
        params: { name: search }
      })
      setTest(response.data.Group)
      console.log(response.data.Group)
    }catch(error){
      console.log(error)
    }
  }

  useEffect(() =>{
    getTests()
  },[search])
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
          {tests?.map(test =>{
            return(
            <CollapseCard
              key={test.ID}
              Id={String(test.ID)}
              Name={test.Name}
              Type={test.Httptype}
              HTTPResult={test.ResponseHttpCode}
              urlEndpoint={test.Urlapi}
              RequestData={test.Request}
              ResponseData={test.Response}
            />)
          })}
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
        <TableUsers_Group />
      </div>
    </div>
    
  );
};

export default GroupData;
