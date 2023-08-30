<template>
    <q-layout view="lHh Lpr lFf">
      <q-header bordered class="text-black bg-primary">
        <q-toolbar>
  
          <q-toolbar-title class="text-bold">

          </q-toolbar-title>
  
          <LayoutChangePopUpComponet class="q-ml-sm"/>
          <CurrentUserBadgeComponent dense class="q-ml-sm"/>
        </q-toolbar>
      </q-header>
  
      <q-drawer
        show-if-above
        bordered
        class="bg-primary"
      >
        <div class="q-pa-lg allign-center">
          <q-img src="/logo.svg"></q-img>
        </div>
        
        
  
        <q-list dense separator class="q-pt-md">
          <q-separator></q-separator>
          <div v-for="m in modules" :key="m.name" class="row">
          <div :class="'col ' + (isModuleRouteActive(m.activePrefix) ? 'bg-secondary' : '')" style="max-width: 5px;border-bottom: #ccc solid 1px;"></div>
          <q-item :class="'col ' + (isModuleRouteActive(m.activePrefix) ? 'bg-grey-2' : '')" clickable v-ripple @click="goToModule(m.routerPathName)" dense style="border-bottom: #ccc solid 1px;">
            <q-item-section avatar>
              <q-icon :name="m.icon" size="sm"></q-icon>
            </q-item-section>
  
            <q-item-section>
              <q-item-label>{{ $t(m.name) }}</q-item-label>
              <!-- <q-item-label caption lines="1">{{ $t(m.description) }}</q-item-label> -->
            </q-item-section>
          </q-item>
          </div>
        </q-list>
      </q-drawer>
  
      <q-page-container class="bg-primary">
        <router-view />
      </q-page-container>
    </q-layout>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import { useRouter, useRoute } from 'vue-router';

  import LayoutChangePopUpComponet from './LayoutChangePopUpComponet.vue';
  import CurrentUserBadgeComponent from '../modules/me/CurrentUserBadgeComponent.vue';
  
  interface ModuleInfo {
    name: string
    description: string
    routerPathName: string
    icon: string
    activePrefix: string
  }
  
  const $router = useRouter()
  const $route = useRoute()
  
  const modules = [
    {
      name: 'layout.me.modules.general.name',
      description: '',
      routerPathName: 'me_general',
      icon: 'person',
      activePrefix: 'me_general'
    },
    {
      name: 'layout.me.modules.auth.name',
      description: '',
      routerPathName: 'me_auth',
      icon: 'key',
      activePrefix: 'me_auth'
    },
    {
      name: 'layout.me.modules.security.name',
      description: '',
      routerPathName: 'me_security',
      icon: 'shield',
      activePrefix: 'me_security'
    }
  ] as ModuleInfo[]
  
  function isModuleRouteActive(prefix: string) {
    return $route.name?.toString().startsWith(prefix) ?? false
  }
  
  async function goToModule(pathName: string) {
    await $router.push({name: pathName, params: { currentNamespace: $route.params.currentNamespace }})
  }
  
  </script>
  