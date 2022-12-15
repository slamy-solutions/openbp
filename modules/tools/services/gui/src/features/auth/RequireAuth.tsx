import { useLocation, Navigate, Outlet } from 'react-router-dom'
import { store } from '../../store'

function RequireAuth() {
    const isLoggedIn = store.getState()['features/auth'].loggedIn
    const location = useLocation()

    return (
        isLoggedIn
            ? <Outlet />
            : <Navigate to="/auth/login" state={{ from: location }} replace />
    )
}

export default RequireAuth