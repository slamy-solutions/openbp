<template>
    <q-btn outline dense icon="apps">
        <q-menu class="q-pa-xs">
            <q-btn icon="settings" @click="goToRoute('home')" class="full-width">Adminer</q-btn>
            <q-btn v-if="modulesStore.crm" icon="developer_board" @click="goToRoute('crm_home')" class="full-width">CRM</q-btn>
        </q-menu>
    </q-btn>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router';
import { useLoginStore } from 'src/stores/login-store'
import { useModulesStore } from 'src/stores/modules-store';

const $router = useRouter()
const $route = useRoute()
const loginStore = useLoginStore()
const modulesStore = useModulesStore()

let currentNamespace = $route.params.currentNamespace as string
if (!currentNamespace) {
    currentNamespace = loginStore.originalNamespace
}
if (!currentNamespace) {
    currentNamespace = "_global"
}

async function goToRoute(name: string) {
    await $router.push({ name: name, params: { currentNamespace } })
}
</script>