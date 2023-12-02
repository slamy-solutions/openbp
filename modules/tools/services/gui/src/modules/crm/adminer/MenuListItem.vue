<template>
    <q-item :clickable="!props.selected" :v-ripple="!props.selected" class="q-pa-xs no-minimums" :disable="props.selected" @click="goToPage">
        <q-item-section avatar class="q-mr-sm q-pa-none q-ma-none no-minimums">
            <q-icon :name="props.icon" :color="(props.selected) ? 'secondary' : ''"></q-icon>
        </q-item-section>
        <q-item-section class="gt-sm">
            <q-item-label>{{ props.title }}</q-item-label>
        </q-item-section>
    </q-item>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router';
const $route = useRoute()

const props = defineProps<{
    name: string,
    icon: string,
    title: string,
    selected: boolean,
    route: string
}>()

const $router = useRouter()

async function goToPage() {
    if (props.route != "" && !props.selected) {
        await $router.push({ name: props.route, params: { currentNamespace: $route.params.currentNamespace } })
    }
}
</script>

<style>

</style>