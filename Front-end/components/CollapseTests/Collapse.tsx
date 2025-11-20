import './style.css'
import React, { useState } from 'react';
import { Collapse, Flex, Input, Select } from 'antd';
import {
  DeleteOutlined,
  SaveOutlined,
  PlayCircleOutlined,
} from '@ant-design/icons';
const { Panel } = Collapse;

interface CollapseCardProps  {
  Id: number;
  Name: string;
  Endpoints: string;
  Members: string;
} 
const { TextArea } = Input;

const CollapseCard: React.FC<CollapseCardProps> = ({Id, Name, Endpoints, Members }) => {

  const CustomHeader = () => (
    <div className='ContainerHeader-Custom'>
      <div style={{display: 'flex', alignItems: 'center', gap:20}}>
        <h3>POST</h3>
        <span>GeeksforGeeks</span>
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
        <h3 style={{margin:0}}>201</h3>
      </div>
    </div>
  );

  const handleChange = (value: string) => {
    console.log(`selected ${value}`);
  };

  const [value, setValue] = useState('');
  
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
                  placeholder="Request"
                  autoSize={{ minRows: 6.2, maxRows: 6.2}}
                />
              </div>
              <div className="div3"> 
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
                <button>Manage Headers</button>
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
