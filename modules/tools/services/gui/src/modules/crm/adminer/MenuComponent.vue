<template>
  <q-card class="full-width full-height bg-transparent" flat >
    <q-card-section>
        <q-list>
            <q-separator />
            <div v-for="view in mainViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName"/>
                <q-separator />
            </div>
        </q-list>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">

import MenuListItem from './MenuListItem.vue';

interface Viewdata {
    name: string
    title: string
    icon: string
    routerPageName: string
}

const props = defineProps<{
    selected: string
}>()

const mainViews = [
    {
        name: 'dashboard',
        icon: 'fa-solid fa-chart-line',
        title: 'Dashboard',
        routerPageName: 'crm_adminer_dashboard'
    },
    {
        name: 'settings',
        icon: 'settings',
        title: 'Settings',
        routerPageName: 'crm_adminer_settings'
    },
] as Viewdata[] 

</script>

<style>
.no-minimums {
    min-width: 0px;
    min-height: 0px;
}
</style>