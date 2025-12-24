import './style.css'
import React from 'react';
import { Input, message } from 'antd';
import type { GetProps } from 'antd';
import axiosInstance from '../../axios';



interface CardGroupProps  {
  Id: number;
  Name: string;
} 

const CardGroup: React.FC<CardGroupProps> = ({Id, Name }) => {
    const sharedProps = {};
    return (
    <>
        <div id={Id.toString()}  className='ContainerGroupCard'>
          <div className='container-data'>
            <h2>{Name}</h2>
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
