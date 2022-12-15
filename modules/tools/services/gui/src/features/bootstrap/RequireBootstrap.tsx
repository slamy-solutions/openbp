import { useLocation, Navigate, Outlet } from 'react-router-dom'
import { store } from '../../store'

function RequireBootstrap() {
    const isBootstrapped = store.getState()['features/bootstrap'].fullyBootstrapped
    const location = useLocation()

    return (
        isBootstrapped
            ? <Outlet />
            : <Navigate to="/bootstrap" state={{ from: location }} replace />
    )
}

export default RequireBootstrap