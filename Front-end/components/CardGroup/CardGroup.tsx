import './style.css'
import React from 'react';
import { Input, message } from 'antd';
import type { GetProps } from 'antd';
import axiosInstance from '../../axios';



interface CardGroupProps  {
  Id: number;
  Name: string;
  Endpoints: string;
  Members: string;
} 

const CardGroup: React.FC<CardGroupProps> = ({Id, Name, Endpoints, Members }) => {
    const sharedProps = {};
    return (
    <>
        <div className='ContainerGroupCard'>
          <div>
            <h2>{Name}</h2>
            <p><b>Endpoints</b>: {Endpoints}</p>
            <p><b>Members</b>: {Members}</p>
            <div className='Buttons-CardGroup'>
              <button className='open-button'>Open</button>
              <button className='delete-button'>Delete</button>
            </div>
          </div>
        </div>
    </>
    );
};

export default CardGroup;
