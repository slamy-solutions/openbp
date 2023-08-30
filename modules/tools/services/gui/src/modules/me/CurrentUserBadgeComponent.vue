<template>
    <q-btn outline icon="person">
        <q-menu class="q-pa-xs bg-primary" size="xs" style="border: #ccc solid 1px" unelevated flat >
            <q-btn outline size="sm" icon="person" @click="goToRoute('me_general')" class="full-width q-mt-xs">Preferences</q-btn>
            <q-btn outline size="sm" icon="logout" @click="goToRoute('logout')" class="full-width q-mt-xs">LOGout</q-btn>
        </q-menu>
    </q-btn>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router';
import { useLoginStore } from 'src/stores/login-store'

const $router = useRouter()
const $route = useRoute()
const loginStore = useLoginStore()


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