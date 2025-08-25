
/* pages */
import Dashboard from '../pages/Dashboard/Dashboard.tsx'
import ClientManagement from '../pages/ClientManagement/ClientManagement.tsx'
import Users from '../pages/Users/Users.tsx'
import Settings from '../pages/Settings/Settings.tsx'
import Profile from '../pages/Profile/Profile.tsx'
import Login from '../pages/Login/Login.tsx'
import MainLayout from './MainLayout';

import { Layout } from 'antd';
//import HeaderBar from '../components/HeaderBar/HeaderBar.tsx';
import { Routes, Route, useLocation } from "react-router-dom";

import './App.css'
import { useEffect, useState } from 'react'

function App() {
  const location = useLocation();
  const [loginPage, setLoginPage] = useState(false)

  useEffect(() => {
    if(location.pathname.includes("/login")){
      setLoginPage(true)
    }
  }, [loginPage])

  return (
    <>
      <Routes>
        <Route path="/" element={<Login />} />

        <Route path="/page" element={<MainLayout />}>
          <Route path="Dashboard" element={<Dashboard />} />
          <Route path="Users" element={<Users />} />
          <Route path="ClientManagement" element={<ClientManagement />} />
          <Route path="Settings" element={<Settings />} />
          <Route path="Profile" element={<Profile />} />
        </Route>
      </Routes>
    
    </>
  )
}

export default App
