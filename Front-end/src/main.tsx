import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom';
import App from './App.tsx'
import {GroupsProvider} from '../context/GroupsContext.tsx'
import { AuthProvider } from '../auth/AuthProvider.tsx'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <AuthProvider>   
        <GroupsProvider>
          <App />
        </GroupsProvider>
      </AuthProvider>
    </BrowserRouter>
  </StrictMode>,
)
