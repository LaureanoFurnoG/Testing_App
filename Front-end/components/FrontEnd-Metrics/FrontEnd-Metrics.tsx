import './style.css'
import React from 'react';
import axiosInstance from '../../axios';
import { Divider, Progress } from 'antd';
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
    <section className='CenterMetrics'>
      <div className='Metrics'>
        <div className='metric-Por'>
          <Progress size={80} percent={100} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
          <h4>SEO</h4>
        </div>
        <div className='metric-Por'>
          <Progress size={80} percent={100} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
          <h4>Performance</h4>
        </div>
        <div className='metric-Por'>
          <Progress size={80} percent={100} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
          <h4>Accessibility</h4>
        </div>
        <div className='metric-Por'>
          <Progress size={80} percent={100} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
          <h4>Best Practices</h4>
        </div>
      </div>
      <Divider />
      <div className='metrics-performance'>
        <div className='m-p-s'>
          <div className='mps-right'>
            <div>
              <p>First Contentful Paint</p>
              <h2>2.5 s</h2>
            </div>
            <div>
               <Divider />
              <p>Total Blocking Time</p>
              <h2>2.5 s</h2>
            </div>
            <div>
               <Divider />
              <p>Speed Index</p>
              <h2>2.5 s</h2>
            </div>
          </div>
          <div className='mps-left'>
            <div>
              <p>Largest Contentful Paint</p>
              <h2>2.5 s</h2>
            </div>
            <div>
              <Divider />
              <p>Cumulative Layout Shift</p>
              <h2>2.5 s</h2>
            </div>
          </div>
        </div>
      </div>
    </section>
    </>
    );
};

export default FrontEnd_MetricsCard;
