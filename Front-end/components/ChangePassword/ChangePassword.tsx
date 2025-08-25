import './style.css'
import React, {useRef, useState} from 'react';
import { Button, Form, Input, message } from 'antd';
import type { FormItemProps } from 'antd';
import MicrosoftButton from '../MicrosoftLogin/MicrosoftLogin'

const MyFormItemContext = React.createContext<(string | number)[]>([]);

interface MyFormItemGroupProps {
  prefix: string | number | (string | number)[];
}

function toArr(str: string | number | (string | number)[]): (string | number)[] {
  return Array.isArray(str) ? str : [str];
}

const MyFormItemGroup: React.FC<React.PropsWithChildren<MyFormItemGroupProps>> = ({
  prefix,
  children,
}) => {
  const prefixPath = React.useContext(MyFormItemContext);
  const concatPath = React.useMemo(() => [...prefixPath, ...toArr(prefix)], [prefixPath, prefix]);

  return <MyFormItemContext.Provider value={concatPath}>{children}</MyFormItemContext.Provider>;
};

const MyFormItem = ({ name, ...props }: FormItemProps) => {
  const prefixPath = React.useContext(MyFormItemContext);
  const concatName = name !== undefined ? [...prefixPath, ...toArr(name)] : undefined;

  return <Form.Item name={concatName} {...props} />;
};

interface ChangePassword {
  CardType: (val: string) => void;
} 

const ChangePassword: React.FC <ChangePassword> = ({ CardType }) => {
    

   //ADD LOGIC FOR CHANGE THE PASSWORD
  return (
    <>
        <div className='Container-cardLogin'>
            <Form name="login-security-scan" className='form-login' layout="vertical" /*onFinish={empty for the moment}*/>
                <h1>Change<span> Password</span></h1>
                <p style={{textAlign:'center'}}>If you have an account, you will receive an email with a code</p>
                <hr style={{ width:'100%', marginTop: 0, marginBottom:20 }}/>
                <MyFormItem name="code" label="Code">
                    <Input required={true} style={{height:54}} placeholder="Code"/>
                </MyFormItem>
                <MyFormItem  name="New Password" label="New Password" >
                    <Input required={true} type='password' style={{height:54}} placeholder="New Password" />
                </MyFormItem>

                <Button style={{height:54, backgroundColor:'#916BF3'}} type="primary" htmlType="submit">
                    CHANGE PASSWORD
                </Button>

            </Form>
            <button className='goToLogin' onClick={() => CardType("login")}>Go to the <span style={{color:"#916BF3"}}>Login</span></button>
            <MicrosoftButton onClick={() => console.log("Login con Microsoft")}/>
        </div>
    </>
  );
};

export default ChangePassword;
