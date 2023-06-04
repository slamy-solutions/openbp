<template>
  <q-layout view="hHh Lpr lFf">
    <q-header elevated style="background: #333333;">
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
          OpenBP
        </q-toolbar-title>

        <div>v0.0.3</div>
        <q-btn icon="logout" dense class="q-ml-md" @click="goToLogout()"></q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
    >
      <q-list>
        <q-item-label
          header
        >
          Modules
        </q-item-label>

        <q-item v-for="m in modules" :key="m.name" clickable v-ripple @click="goToModule(m.routerPathName)">
          <q-item-section avatar>
            <q-icon :name="m.icon" size="md"></q-icon>
          </q-item-section>

          <q-item-section>
            <q-item-label>{{ $t(m.name) }}</q-item-label>
            <q-item-label caption lines="1">{{ $t(m.description) }}</q-item-label>
          </q-item-section>
        </q-item>
        

        <!-- <EssentialLink
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
        /> -->
      </q-list>
    </q-drawer>

    <q-page-container class="bg-primary">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

interface ModuleInfo {
  name: string
  description: string
  routerPathName: string
  icon: string
}

const $router = useRouter()

const leftDrawerOpen = ref(false)
const modules = [
  {
    name: 'layout.main.modules.namespace.name',
    description: 'layout.main.modules.namespace.description',
    routerPathName: 'namespaceList',
    icon: 'grid_view'
  },
  {
    name: 'layout.main.modules.accessControl.name',
    description: 'layout.main.modules.accessControl.description',
    routerPathName: 'accessControl_iam_identity_list',
    icon: 'fingerprint'
  }
] as ModuleInfo[]

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

async function goToModule(pathName: string) {
  await $router.push({name: pathName})
}

async function goToLogout() {
  await $router.push({name: "logout"})
}

</script>
