import { Button } from "antd";
import React from "react";
import microsoftIcon from '../../assets/images/microsoftIcon.webp';
import PropTypes from 'prop-types';

interface MicrosoftButtonProps {
  onClick: () => void;
} //typing of TypeScript.

const MicrosoftButton: React.FC<MicrosoftButtonProps> = ({ onClick }) => {
  return (
    <Button 
      type="default" 
    
      icon={<img src={microsoftIcon} alt="Microsoft"  style={{ width: 25, height: 25 }} />} 
      onClick={onClick} 
      block
      style={{  borderRadius:"50px", height: 54}}
    >
      Login with Microsoft
    </Button>
  );
};



export default MicrosoftButton;
