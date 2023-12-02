
import { RouteRecordRaw } from 'vue-router';
import { bootstrapGuard, loginGuard } from './guards'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: {name: 'bootstrap' }
  },
  {
    path: '/_bootstrap',
    name: 'bootstrap',
    component: () => import('../modules/bootstrap/BootstrapPage.vue')
  },
  {
    path: '/:currentNamespace/login',
    name: 'login',
    component: () => import('../modules/login/LoginPage.vue'),
    beforeEnter: [bootstrapGuard]
  },
  {
    path: '/:currentNamespace/login/finalizeoauth/:provider',
    name: 'login_oauth_finalize',
    component: () => import('../modules/login/FinalizeOAuthLoginPage.vue'),
    beforeEnter: [bootstrapGuard]
  },
  {
    path: '/:currentNamespace/logout',
    name: 'logout',
    component: () => import('../modules/login/LogoutPage.vue'),
    beforeEnter: [bootstrapGuard]
  },
  {
    path: '/:currentNamespace',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', name: 'home', component: () => import('../modules/IndexPage.vue') },

      { path: 'namespace/list', name: 'namespaceList', component: () => import('../modules/namespace/NamespaceListPage.vue')},

      { path: 'accessControl/actor/user/list', name: 'accessControl_actor_user_list', component: () => import('../modules/accessControl/actor/user/UserListPage.vue')},
      { path: 'accessControl/iam/identity/list', name: 'accessControl_iam_identity_list', component: () => import('../modules/accessControl/iam/identity/IdentityListPage.vue')},
      { path: 'accessControl/iam/policy/list', name: 'accessControl_iam_policy_list', component: () => import('../modules/accessControl/iam/policy/PolicyListPage.vue')},
      { path: 'accessControl/iam/role/list', name: 'accessControl_iam_role_list', component: () => import('../modules/accessControl/iam/role/RoleListPage.vue')},
      { path: 'accessControl/settings/oauth', name: 'accessControl_settings_oauth', component: () => import('../modules/accessControl/settings/oauth/OAuthSettingsPage.vue')},


      { path: 'iot/fleet/list', name: 'iot_fleet_list', component: () => import('../modules/iot/fleet/FleetListPage.vue') },
      { path: 'iot/integration/balena', name: 'iot_integration_balena', component: () => import('../modules/iot/integration/balena/BalenaIntegrationPage.vue') },

      { path: 'crm/adminer/dashboard', name: 'crm_adminer_dashboard', component: () => import('../modules/crm/adminer/dashboard/DashboardPage.vue') },
      { path: 'crm/adminer/settings', name: 'crm_adminer_settings', component: () => import('../modules/crm/adminer/settings/SettingsPage.vue') },
    ],
    beforeEnter: [bootstrapGuard, loginGuard]
  },
  {
    path: '/_me',
    component: () => import('layouts/MeLayout.vue'),
    children: [
      { path: '/general', name: 'me_general', component: () => import('../modules/me/generalInfo/GeneralInfoPage.vue') },
      { path: '/auth', name: 'me_auth', component: () => import('../modules/me/auth/AuthPage.vue') },
      { path: '/auth/oauth/finalize/:provider', name: 'me_auth_oauth_finalize', component: () => import('../modules/me/auth/FinalizeOAuthRegistrationPage.vue')},
    ],
    beforeEnter: [bootstrapGuard, loginGuard]
  },

  {
    path: '/crm',
    component: () => import('layouts/CRMLayout.vue'),
    children: [
      { path: '', name: 'crm_home', component: () => import('../modules/crm/IndexPage.vue') }
    ],
    beforeEnter: [bootstrapGuard, loginGuard]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('../modules/ErrorNotFound.vue'),
  },
];

export default routes;
