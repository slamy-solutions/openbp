<template>
<q-layout view="hHh Lpr lFf">
    <q-header bordered class="header-bar text-black">
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
        Kanban CRM
        </q-toolbar-title>

        <div>v0.0.3</div>
        <LayoutChangePopUpComponet class="q-ml-md"/>
        <CurrentUserBadgeComponent dense class="q-ml-sm"/>
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
        Tools
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

    <q-page-container class="layout-container">
    <router-view />
    </q-page-container>
</q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

import LayoutChangePopUpComponet from './LayoutChangePopUpComponet.vue';
import CurrentUserBadgeComponent from 'src/modules/me/CurrentUserBadgeComponent.vue';

interface ModuleInfo {
name: string
description: string
routerPathName: string
icon: string
}

const $router = useRouter()

const leftDrawerOpen = ref(false)
const modules = [] as ModuleInfo[]

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

<style lang="scss" scoped>
.layout-container {
    background-color: $crm-background;
}
.header-bar {
    background-color: $crm-background-panel;
}
</style>