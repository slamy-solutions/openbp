<template>
  <q-card class="full-width full-height">
    <q-card-section>
      <div class="text-h6 gt-sm">IoT</div>
      <div class="text-subtitle2"></div>
    </q-card-section>

    <q-card-section>
        <q-list>
            <q-separator />
            <div v-for="view in mainViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName"/>
                <q-separator />
            </div>
        </q-list>
    </q-card-section>

    <q-card-section>
        <div>Integrations</div>
        <q-list>
            <q-separator />
            <div v-for="view in integrationViews" :key="view.name">
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
        name: 'fleets',
        icon: 'fa-solid fa-layer-group',
        title: 'Device fleets',
        routerPageName: 'iot_fleet_list'
    },
    {
        name: 'applications',
        icon: 'terminal',
        title: 'Applications',
        routerPageName: ''
    },
] as Viewdata[] 

const integrationViews = [
    {
        name: 'balena',
        icon: '',
        title: 'Balena',
        routerPageName: 'iot_integration_balena'
    }
] as Viewdata[]

</script>

<style>
.no-minimums {
    min-width: 0px;
    min-height: 0px;
}
</style>