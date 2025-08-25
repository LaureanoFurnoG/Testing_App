import './style.css'
import React, {} from 'react';
import { Input, message } from 'antd';
import type { GetProps } from 'antd';
import axiosInstance from '../../axios.js';
import { useNavigate } from 'react-router-dom';


interface CardType {
  CardType: (val: string) => void;
} 

type OTPProps = GetProps<typeof Input.OTP>;

const VerificatorAuth: React.FC <CardType> = ({ CardType }) => {
  const [messageApi, contextHolder] = message.useMessage();
  let navigate = useNavigate()

  const dashboardNavigate = () => {
    navigate('/Dashboard')
  };
  const onChange: OTPProps['onChange'] = async (text) => {
    try {
      const value = {
        email: localStorage.getItem('email'),
        encriptedCode: localStorage.getItem('codeEncripted'),
        code: text
      }
      const response = await axiosInstance.post("/login/authenticationCode", value);
      localStorage.setItem('token', response.data.token)
      localStorage.removeItem('email')
      localStorage.removeItem('codeEncripted')
      dashboardNavigate()
      
    } catch (error: any) {
      console.log(error)
      errorM(error)
    }
  };

  const errorM = (mText: String) => {
    messageApi.open({
    type: 'error',
    content: mText,
    });
  };
  const sharedProps: OTPProps = {
    onChange,
  };
   //ADD LOGIC FOR CHANGE THE PASSWORD
  return (
    <>
        <div className='Container-cardVerification'>
          <h1>Authentication<span> Code</span></h1>
          <p style={{textAlign:'center'}}>If you have an account, you will receive an email with a code</p>
          <hr style={{ width:'100%', marginTop: 0, marginBottom:20 }}/>
          <div className='codeBox'>
            <Input.OTP  formatter={(str) => str.toUpperCase()} {...sharedProps} />
          </div>
          <button className='loginRedirect' onClick={() => CardType("login")}>Go to the <span style={{color:"#916BF3"}}>Login</span></button>
        </div>
    </>
  );
};

export default VerificatorAuth;
