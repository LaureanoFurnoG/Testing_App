
/* pages */
import GroupsManagement from '../pages/Groups/Groups.tsx'
import Documentation from '../pages/Documentation/Docuementation.tsx'
import Users from '../pages/GroupData/GroupData.tsx'
import Settings from '../pages/Settings/Settings.tsx'
import Profile from '../pages/Profile/Profile.tsx'
import Login from '../pages/Login/Login.tsx'
import MainLayout from './MainLayout';
import Homepage from '../pages/Homepage/Homepage.tsx';
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
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<Homepage />} />
        <Route path="/" element={<MainLayout />}>
          <Route path="Groups" element={<GroupsManagement />} />
          <Route path="Users" element={<Users />} />
          <Route path="Documentation" element={<Documentation />} />
          <Route path="Settings" element={<Settings />} />
          <Route path="Profile" element={<Profile />} />
        </Route>
      </Routes>
    
    </>
  )
}

export default App
