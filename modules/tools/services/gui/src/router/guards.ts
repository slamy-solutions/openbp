import api from 'src/boot/api';
import { useBootstrapStore } from 'src/stores/bootstrap-store';
import { useLoginStore } from 'src/stores/login-store';
import { NavigationGuardNext, RouteLocationNormalized } from 'vue-router';

// Checks if system is bootstrapped. If not - redirects to the bootstrap page.
export async function bootstrapGuard(to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) {
  const bootstrapStore = useBootstrapStore()
  if (!bootstrapStore.bootstrapped) {
    try {
      const status = await api.bootstrap.getStatus()
      bootstrapStore.updateBootstrapState(status)
    } catch {}
    
    if (!bootstrapStore.bootstrapped) {
      return next({ name: 'bootstrap' })
    }
  }

  return next()
}

// Checks if user is logged in. If not - redirects to the login page.
export async function loginGuard(to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) {
    const loginStore = useLoginStore()
    if (!loginStore.isLoggedIn) {
        loginStore.tryLoadLoginFromStorage()

        if (loginStore.isLoggedIn) {
          try {
              const tokenValid = await api.login.validateToken({
                  token: loginStore.refreshToken
              })
              if (!tokenValid) (
                  loginStore.logout()
              )
          } catch {
              loginStore.logout()
          }
        }

        if (!loginStore.isLoggedIn) {
            return next({ name: 'login' })
        }
    }
  
    return next()
  }