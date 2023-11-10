<template>
  <q-layout view="lHh Lpr lFf">
    <q-header bordered class="text-black bg-primary">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title class="text-bold">
          <span class="text-subtitle1">{{ $route.params.currentNamespace }}</span>
        </q-toolbar-title>

        <LayoutChangePopUpComponet class="q-ml-sm"/>
        <CurrentUserBadgeComponent dense class="q-ml-sm"/>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
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

const leftDrawerOpen = ref(false)
const modules = [
  {
    name: 'layout.main.modules.accessControl.name',
    description: 'layout.main.modules.accessControl.description',
    routerPathName: 'accessControl_iam_identity_list',
    icon: 'fingerprint',
    activePrefix: 'accessControl'
  },
  {
    name: 'layout.main.modules.iot.name',
    description: 'layout.main.modules.iot.description',
    routerPathName: 'iot_fleet_list',
    icon: 'fa-solid fa-cubes',
    activePrefix: 'iot'
  },
  {
    name: 'layout.main.modules.crm.name',
    description: 'layout.main.modules.crm.description',
    routerPathName: 'crm_settings',
    icon: 'fa-solid fa-cog',
    activePrefix: 'crm'
  }
] as ModuleInfo[]

if ($route.params.currentNamespace === '_global') {
  modules.unshift({
    name: 'layout.main.modules.namespace.name',
    description: 'layout.main.modules.namespace.description',
    routerPathName: 'namespaceList',
    icon: 'grid_view',
    activePrefix: 'namespace'
  })
}

function isModuleRouteActive(prefix: string) {
  return $route.name?.toString().startsWith(prefix) ?? false
}

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

async function goToModule(pathName: string) {
  await $router.push({name: pathName, params: { currentNamespace: $route.params.currentNamespace }})
}

</script>
