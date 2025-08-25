import { Layout } from 'antd';
import SideMenu from '../components/SideMenu/SideMenu.tsx';
import { Outlet } from 'react-router-dom';

const { Content } = Layout;

export default function MainLayout() {
  return (
    <Layout style={{ height:'100%' }}>
      <SideMenu />
      <Layout className='ContainerContent-Page' style={{ padding: '24px', height:'100%' }}>
        <Content style={{ background: '#fff', borderRadius: 4 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
