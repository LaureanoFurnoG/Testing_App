import React from 'react';
import { Layout, Avatar, Dropdown, Menu } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Header } = Layout;

const HeaderBar: React.FC = () => {
  // Optional dropdown menu
  let navigate = useNavigate()
   const SettingsNavigate = () => {
    navigate('/Settings')
  };
  const ProfileNavigate = () => {
    navigate('/Profile')
  }
  const userMenu = (
    <Menu className='usersOptionMenu'
      items={[
        { key: 'profile', label: 'Profile', onClick: ProfileNavigate },
        { key: 'settings', label: 'Settings', onClick: SettingsNavigate },
        { type: 'divider' },
        { key: 'logout', label: 'Logout' },
      ]}
    />
  );

  return (
    <Header
      style={{
        background: '#fff',
        borderBottom: '1px solid #f0f0f0',
        padding: '0 24px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
      }}
    >
      {/* Right side: Avatar + User Name + Optional Dropdown */}
      <Dropdown overlay={userMenu} trigger={['click']}>
        <div style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
          <Avatar
            src="https://joeschmoe.io/api/v1/random"
            size="large"
            icon={<UserOutlined />}
          />
          <span style={{ marginLeft: 8 }}>Serati Ma</span>
        </div>
      </Dropdown>
    </Header>
  );
};

export default HeaderBar;
