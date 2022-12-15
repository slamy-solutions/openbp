
import {Routes, Route} from 'react-router-dom'

import RequireBootstrap from './features/bootstrap/RequireBootstrap'
import BootstrapPage from './features/bootstrap/BootstrapPage'

import LoginPage from './features/auth/LoginPage'

function App() {
  return (
      <Routes>
        <Route path="/bootstrap" element={ <BootstrapPage /> }></Route>
        
        <Route path="/" element={<RequireBootstrap />}>
          <Route path="auth/login" element={ <LoginPage /> } />
        </Route>
      </Routes>
  );
}

export default App;
