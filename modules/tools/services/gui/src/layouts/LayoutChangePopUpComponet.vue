<template>
    <q-btn square dense icon="apps">
        <q-menu class="q-pa-xs">
            <q-btn icon="settings" @click="goToLayout('home')" class="full-width">Adminer</q-btn>
            <q-btn icon="developer_board" @click="goToLayout('crm_home')" class="full-width">CRM</q-btn>
        </q-menu>
    </q-btn>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';


const $router = useRouter()

async function goToLayout(name: string) {
    await $router.push({ name })
}

</script>