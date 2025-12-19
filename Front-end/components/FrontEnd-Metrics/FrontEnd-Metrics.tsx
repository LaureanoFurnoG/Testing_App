import './style.css'
import React, { useEffect, useState } from 'react';
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

  const fetchMetrics = async () => {
    const body = {
      Url: URL,
      Strategy:"Mobile"
    };
    const response = await axiosInstance.post('/api/tests/6/test-front', body)

    setSeoScore(response.data.result.lighthouseResult.categories.seo.score*100)
    setPerformanceScore(response.data.result.lighthouseResult.categories.performance.score*100)
    setAccessibilityScore(response.data.result.lighthouseResult.categories.accessibility.score*100)
    setPracticesScore(response.data.result.lighthouseResult.categories["best-practices"].score*100)
  };

  useEffect(() => {
    if (URL.trim() !== "") {
      fetchMetrics();  
    }
  }, [URL]); 

  const [SEOscore, setSeoScore] = useState<number>(0);
  const [Performancescore, setPerformanceScore] = useState<number>(0);
  const [Accessibilityscore, setAccessibilityScore] = useState<number>(0);
  const [Practicesscore, setPracticesScore] = useState<number>(0);

  const [ContentfulTime, setContentfulTime] = useState("0.0");
  const [BlockingTime, setBlockingTime] = useState("0.0");
  const [SpeedIndex, setSpeedIndex] = useState("0.0");
  const [LargestContentful, setLargestContentful] = useState("0.0");
  const [LayoutShift, setLayoutShift] = useState("0.0");
  return (
  <>
  <section className='CenterMetrics'>
    <div className='Metrics'>
      <div className='metric-Por'>
        <Progress size={80} percent={SEOscore} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
        <h4>SEO</h4>
      </div>
      <div className='metric-Por'>
        <Progress size={80} percent={Performancescore} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
        <h4>Performance</h4>
      </div>
      <div className='metric-Por'>
        <Progress size={80} percent={Accessibilityscore} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
        <h4>Accessibility</h4>
      </div>
      <div className='metric-Por'>
        <Progress size={80} percent={Practicesscore} strokeColor={conicColors} success={{ percent: 0 }} type="dashboard" />
        <h4>Best Practices</h4>
      </div>
    </div>
    <Divider />
    <div className='metrics-performance'>
      <div className='m-p-s'>
        <div className='mps-right'>
          <div>
            <p>First Contentful Paint</p>
            <h2>{ContentfulTime} s</h2>
          </div>
          <div>
              <Divider />
            <p>Total Blocking Time</p>
            <h2>{BlockingTime} s</h2>
          </div>
          <div>
              <Divider />
            <p>Speed Index</p>
            <h2>{SpeedIndex} s</h2>
          </div>
        </div>
        <div className='mps-left'>
          <div>
            <p>Largest Contentful Paint</p>
            <h2>{LargestContentful} s</h2>
          </div>
          <div>
            <Divider />
            <p>Cumulative Layout Shift</p>
            <h2>{LayoutShift} s</h2>
          </div>
        </div>
      </div>
    </div>
  </section>
  </>
  );
};

export default FrontEnd_MetricsCard;
