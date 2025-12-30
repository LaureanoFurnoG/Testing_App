import './style.css'
import React, { useState } from 'react';
import { Collapse, Input, Select } from 'antd';
import {
  DeleteOutlined,
  SaveOutlined,
  PlayCircleOutlined,
} from '@ant-design/icons';
import AddHeaderDrawer from '../AddHeader-Drawer/AddHeader-Drawer';
const { Panel } = Collapse;

interface CollapseCardProps  {
  Id: string;
  Name: string;
  Type: string;
  HTTPResult: number;
  urlEndpoint: string;
  RequestData:Record<string, any>;
  ResponseData: JSON | null;
} 
const { TextArea } = Input;

const CollapseCard: React.FC<CollapseCardProps> = ({Id, Name, Type, HTTPResult, urlEndpoint, RequestData, ResponseData }) => {

  const CustomHeader = () => (
    <div className='ContainerHeader-Custom' id={Id}>
      <div style={{display: 'flex', alignItems: 'center', gap:20}}>
        <h3>{Type}</h3>
        <span>{Name}</span>
      </div>
      <div style={{display:'flex', alignItems: 'center', gap:20, marginRight:20}}>
        <DeleteOutlined onClick={(e) => {
          e.stopPropagation();
          console.log("Eliminar");}} style={{color:"red"}}/>
        <SaveOutlined onClick={(e) => {
          e.stopPropagation();
          console.log("Save");}} style={{color:"blue"}} />
        <PlayCircleOutlined onClick={(e) => {
          e.stopPropagation();
          console.log("Play");}} style={{color:"green"}} />
        <h3 style={{margin:0}}>{HTTPResult}</h3>
      </div>
    </div>
  );

  const handleChange = (value: string) => {
    console.log(`selected ${value}`);
  };

  const [value, setValue] = useState<string>(
    JSON.stringify(RequestData, null, 2)
  );

  const [Placeholder, _] = useState(`{
  "Name": "Example",
  "Password": "example"
}`)

  return (
  <>
    <div className='container-collapses'>
      <Collapse>
          <Panel header={<CustomHeader />} key="1">
            <div className="parent">
              <div className="div1"></div>
              <div className="div2"> 
                <TextArea
                  style={{padding:20, borderRadius:0, border:"none"}}
                  value={value}
                  onChange={(e) => setValue(e.target.value)}
                  placeholder={Placeholder}
                  autoSize={{ minRows: 6.2, maxRows: 6.2}}
                />
              </div>
              <div className="div3"> 
                <p>Endpoint URL: {urlEndpoint}</p>
                <div style={{display:"flex", justifyContent:"space-between"}}>
                  <div>
                    <Input placeholder="Token" />
                    <Select
                      defaultValue="Body"
                      style={{ width: 120 }}
                      onChange={handleChange}
                      options={[
                        { value: 'Body', label: 'Body' },
                        { value: 'Param', label: 'Param' },
                      ]}
                    />
                  </div>
                  <AddHeaderDrawer />
                </div>
              </div>
              <div className="div4"> </div>
            </div>
          </Panel>
      </Collapse>
    </div>
  </>
  );
};

export default CollapseCard;
