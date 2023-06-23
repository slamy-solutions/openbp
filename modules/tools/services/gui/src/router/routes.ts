
import { RouteRecordRaw } from 'vue-router';
import { bootstrapGuard, loginGuard } from './guards'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../modules/login/LoginPage.vue'),
    beforeEnter: [bootstrapGuard]
  },

  {
    path: '/logout',
    name: 'logout',
    component: () => import('../modules/login/LogoutPage.vue'),
    beforeEnter: [bootstrapGuard]
  },

  {
    path: '/bootstrap',
    name: 'bootstrap',
    component: () => import('../modules/bootstrap/BootstrapPage.vue')
  },

  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', name: 'home', component: () => import('../modules/IndexPage.vue') },

      { path: 'namespace/list', name: 'namespaceList', component: () => import('../modules/namespace/NamespaceListPage.vue')},

      { path: 'accessControl/iam/identity/list', name: 'accessControl_iam_identity_list', component: () => import('../modules/accessControl/IdentityListPage.vue')},
      { path: 'accessControl/iam/policy/list', name: 'accessControl_iam_policy_list', component: () => import('../modules/accessControl/iam/policy/PolicyListPage.vue')},
      { path: 'accessControl/iam/role/list', name: 'accessControl_iam_role_list', component: () => import('../modules/accessControl/iam/role/RoleListPage.vue')},
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
