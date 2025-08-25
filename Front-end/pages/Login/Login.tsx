import React from 'react';
import 'antd/dist/reset.css'
import imgBannerLogin from '../../assets/images/imageLogin.png'
import SchubLogo from '../../assets/images/Schub_logo.webp'
import LoginCard from '../../components/LoginCard/LoginCard.tsx'
import ChangePassword from '../../components/ChangePassword/ChangePassword.tsx'
import VerificatorAuth from '../../components/Verification/Verification.tsx'
import './style.css'


const Login: React.FC = () => {
  const [CardType, setCardType] = React.useState<string>('login');
  return (
    <div className='LoginContainer'>
      <div className='ImgBanner-Login'>
        <h1>TESTING <span>APPLICATION</span></h1>
        <img className='img-Banner' src={imgBannerLogin} alt="Banner image for the login page Schub tool" />
        <h3 className='aboutApp'>Application for<span> testing </span>both backend and frontend applications.</h3>
        <h3 className='powerBy-Schub'>Power by <span>Schub</span></h3>
      </div>
      <div className="login-right">
        <div className='CardLogin-container'>
          <img className='img-Banner-responsive' src={SchubLogo} alt="Banner image for the login page Schub tool" />
          {CardType === 'change' ? (
            <ChangePassword CardType={setCardType} />
          ) : CardType === 'verify' ? (
            <VerificatorAuth CardType={setCardType} />
          ) : (
            <LoginCard CardType={setCardType} />
          )}
        </div>
      </div>
    </div>
  );
};

export default Login;
