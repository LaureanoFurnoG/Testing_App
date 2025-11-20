import './style.css'
import React from 'react';
import axiosInstance from '../../axios';
import { Progress } from 'antd';
import type { ProgressProps } from 'antd';



interface FrontEnd_Metrics  {
  URL: string;
} 

const FrontEnd_MetricsCard: React.FC<FrontEnd_Metrics> = ({ URL }) => {
    const conicColors: ProgressProps['strokeColor'] = {
      '0%': 'red',
      '50%': 'orange',
      '100%': 'green',
    };
    return (
    <>
      <Progress percent={30} strokeColor={conicColors} success={{ percent: 30 }} type="dashboard" />
    </>
    );
};

export default FrontEnd_MetricsCard;
