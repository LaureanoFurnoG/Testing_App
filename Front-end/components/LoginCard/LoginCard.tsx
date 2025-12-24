import './style.css'
import React, {useRef, useState} from 'react';
import { Button, Form, Input, message } from 'antd';
import type { FormItemProps } from 'antd';
import ReCAPTCHA from "react-google-recaptcha";
import axiosInstance from '../../axios';

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
interface LoginCardProps {
  CardType: (val: string) => void;
} //typing of TypeScript.

const LoginCard: React.FC <LoginCardProps> = ({ CardType }) => {
    const captcha = useRef<ReCAPTCHA | null>(null);
    const [, setValidSession] = useState(false)
    const [captchaValid, setCaptchaValid] = useState<boolean | null>(null)
    const [, setUserValid] = useState(false)
    const [messageApi, ] = message.useMessage();
    const [loading, setLoading] = useState(false);

    function onChange() {
        if (captcha.current.getValue()) {
            setCaptchaValid(true)
        }
    }

    const onFinish = async (value: { email: string; password: string }) => {
        setLoading(true);

        try {
            const captchaValue = captcha.current.getValue();
            if (!captchaValue) {
                setUserValid(false);
                setCaptchaValid(false);
                setLoading(false);
                return; 
            }
            setUserValid(true);
            setCaptchaValid(true);
            try {
                const response = await axiosInstance.post("/api/user/login", value);
                setValidSession(true);
                localStorage.setItem('email', value.email);
                CardType("verify");
            } catch (error: any) {
                errorM(error); 
            }

        } catch (error: any) {
            setLoading(false);
            errorM(error);
        } finally {
            setLoading(false); 
        }
    };


    const errorM = (mText: String) => {
        messageApi.open({
        type: 'error',
        content: mText,
        });
    };
  return (
    <>
        <div className='Container-cardLogin'>
            <Form name="login-security-scan" className='form-login' layout="vertical" onFinish={onFinish}>
                <h1>WELCOME<span>!</span></h1>
                <MyFormItem name="email" label="Email">
                    <Input required={true} type='email' style={{height:54}} placeholder="Email"/>
                </MyFormItem>
                <MyFormItem  name="password" label="Password" >
                    <Input required={true} type='password' style={{height:54}} placeholder="Password" />
                </MyFormItem>

                <Button loading={loading}  style={{height:54, backgroundColor:'#236d55'}} type="primary" htmlType="submit">
                    LOGIN
                </Button>
            </Form>
            <button className='forgotPasswordBTN' onClick={() => CardType("change")}>Forgot Passoword?</button>
            <div className="captcha-container">
                <ReCAPTCHA className='captchaD'
                ref={captcha}
                sitekey={import.meta.env.VITE_CAPTCHA_KEY}
                onChange={onChange}
                />,
            </div>
            {captchaValid === false && <p className="error-captcha">Please, accept the captcha</p>}
        </div>
    </>
  );
};

export default LoginCard;
