import React, { useEffect, useState } from 'react';
import { Layout, Menu, Dropdown, Avatar, type MenuProps  } from 'antd';
import {
  FileDoneOutlined,
  UsergroupAddOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import './SideMenu.css'
const { Sider } = Layout;

const SideMenu: React.FC = () => {
  // Local state to track collapsed/expanded status
  const [collapsed, setCollapsed] = useState(false);
  const [LocationApp, setLocationApp] = useState(['0'])
  const location = useLocation();
  let navigate = useNavigate()

  const manageGroupsNavigate = () => {
    navigate('/Groups')
  };

  const groupNavigate = () => {
    navigate('/GroupData')
  };

  const DocumentationNavigate = () => {
    navigate('/Documentation')
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
    label: 'Groups',
    icon: <UsergroupAddOutlined />,
    children: [
      //i think that we can use a loop and travel the clients table, find the name and insert in the label
      {
        key: 'g1',
        label: 'Manage Groups',
        onClick: manageGroupsNavigate
      },
      {
        key: 'g2',
        label: 'G_Name',
        onClick: groupNavigate
      },
    ],
  },
  {
    key: '2',
    label: 'Documentation',
    icon: <FileDoneOutlined />,
    onClick: DocumentationNavigate
  },
  ]

  useEffect(() =>{
   const pathname = location.pathname
   switch(true){
    case pathname.includes("Groups"):
      setLocationApp(['0'])
      break;
    case pathname.includes("Users"):
      setLocationApp(['1'])
      break;
    case pathname.includes("Documentation"):
      setLocationApp(['2'])
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
            <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#ffffffff"><path d="M120-240v-80h720v80H120Zm0-200v-80h720v80H120Zm0-200v-80h720v80H120Z" /></svg>
          </div>
          : <h1 style={{margin: 10, marginLeft: 0}}>TESTING APP</h1>}
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
                <p style={{ marginLeft: 8, marginBottom: 1, fontWeight: 'bold', fontSize: 14, color:'white'}}>Serati Ma</p>
                <p style={{ marginLeft: 8, marginTop: 1, fontSize: 12, color: '#ffffffda' }}>Serati@ma.com</p>
              </div>}
          </div>
        </Dropdown>
      </div>
    </div>
  </>
    
  );
};

export default SideMenu;
