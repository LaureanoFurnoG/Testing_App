import React, { useEffect, useState } from 'react';
import { Layout, Menu, Dropdown, Avatar, type MenuProps  } from 'antd';
import {
  DashboardOutlined,
  PieChartOutlined,
  UserAddOutlined,
  UsergroupAddOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import imageLogo from '../../assets/images/Schub_logo.webp'
import './SideMenu.css'
const { Sider } = Layout;

const SideMenu: React.FC = () => {
  // Local state to track collapsed/expanded status
  const [collapsed, setCollapsed] = useState(false);
  const [LocationApp, setLocationApp] = useState(['0'])
  const location = useLocation();
  let navigate = useNavigate()

  const dashboardNavigate = () => {
    navigate('/Dashboard')
  };
  const usersNavigate = () => {
    navigate('/Users')
  };

  const ClientsManagementNavigate = () => {
    navigate('/ClientManagement')
  };
   const SettingsNavigate = () => {
    navigate('/Settings')
  };
  const ProfileNavigate = () => {
    navigate('/Profile')
  }
  const LogOut = () => {
    localStorage.removeItem('token')
    navigate('/login')
  }
  const userMenu = (
    <Menu className='usersOptionMenu'
      items={[
        { key: 'profile', label: 'Profile', onClick: ProfileNavigate },
        { key: 'settings', label: 'Settings', onClick: SettingsNavigate },
        { type: 'divider' },
        { key: 'logout', label: 'Logout', onClick: LogOut },
      ]}
    />
  );
  type MenuItem = Required<MenuProps>['items'][number];

  const items: MenuItem[] = [
  {
    key: '1',
    label: 'Dashboard',
    icon: <DashboardOutlined />,
    onClick: dashboardNavigate
  },
  {
    key: '2',
    label: 'Users',
    icon: <UsergroupAddOutlined />,
    onClick: usersNavigate
  },
  {
    key: '3',
    label: 'Clients',
    icon: <PieChartOutlined />,
    children: [
      //i think that we can use a loop and travel the clients table, find the name and insert in the label
      {
        key: 'g1',
        label: 'Schub',
        onClick: usersNavigate
      },
    ],
  },
  {
    key: '4',
    label: 'Clients Management',
    icon: <UserAddOutlined />,
    onClick: ClientsManagementNavigate
  },
  ]

  useEffect(() =>{
   const pathname = location.pathname
   switch(true){
    case pathname.includes("Dashboard"):
      setLocationApp(['1'])
      break;
    case pathname.includes("Users"):
      setLocationApp(['2'])
      break;
    case pathname.includes("ClientManagement"):
      setLocationApp(['4'])
      break;
    default:
        setLocationApp(['0']);
   }
  },[location])
  return (
     <>
    <div className={collapsed ? 'sideMenu collapsed' : 'sideMenu'}>
    <Sider
      collapsible // enable collapsible behavior
      collapsed={collapsed}
      onCollapse={(value) => setCollapsed(value)}
      width={240}
      style={{
        background: '#fff',
        borderRight: '1px solid #f0f0f0',
      }}
      trigger={null}
    >
      <div onClick={() => setCollapsed(!collapsed)} className='SideMenu-activeButton'>
        {collapsed ?
          <div className='disableMenu-label'>
            <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#1f1f1f"><path d="M120-240v-80h720v80H120Zm0-200v-80h720v80H120Zm0-200v-80h720v80H120Z" /></svg>
          </div>
          : <img className='img_LogoSchub' src={imageLogo} alt="" />}
      </div>
      <Menu
        mode="inline"
        theme="light"
        selectedKeys={LocationApp}
        style={{ height: '90%', borderRight: 0 }}
        items={items}
      >
      </Menu>

    </Sider>
    
    <div className='dataUser-sideMenu' >
        <Dropdown overlay={userMenu} trigger={['click']}>
          <div style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
            <Avatar
              src="https://joeschmoe.io/api/v1/random"
              size="large"
              icon={<UserOutlined />} />
            {collapsed ? '' :
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                <p style={{ marginLeft: 8, marginBottom: 1, fontWeight: 'bold', fontSize: 14 }}>Serati Ma</p>
                <p style={{ marginLeft: 8, marginTop: 1, fontSize: 12, color: 'gray' }}>Serati@ma.com</p>
              </div>}
          </div>
        </Dropdown>
      </div>
    </div>
  </>
    
  );
};

export default SideMenu;
