<template>
  <q-card class="full-width full-height bg-transparent" flat style="border-right: #ccc solid 1px;">
    <q-card-section>
        <div class="text-overline">Actors</div>
        <q-list>
            <q-separator />
            <div v-for="view in actorViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName" :icon-as-image="view.iconAsImage"/>
                <q-separator />
            </div>
        </q-list>
    </q-card-section>

    <q-card-section>
        <div class="text-overline">IAM</div>
        <q-list>
            <q-separator />
            <div v-for="view in iamViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName" :icon-as-image="view.iconAsImage"/>
                <q-separator />
            </div>
        </q-list>
    </q-card-section>

    <q-card-section>
        <div class="text-overline">Settings and security</div>
        <q-list>
            <q-separator />
            <div v-for="view in settingsViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName" :icon-as-image="view.iconAsImage"/>
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
    iconAsImage?: boolean
}

const props = defineProps<{
    selected: string
}>()

const actorViews = [
    {
        name: 'users',
        icon: 'person',
        title: 'Users',
        routerPageName: 'accessControl_actor_user_list'
    },
    {
        name: 'application',
        icon: 'dns',
        title: 'Applications',
        routerPageName: ''
    },
] as Viewdata[] 

const iamViews = [
    {
        name: 'identities',
        icon: 'badge',
        title: 'Identities',
        routerPageName: 'accessControl_iam_identity_list'
    },
    {
        name: 'roles',
        icon: 'fact_check',
        title: 'Roles',
        routerPageName: 'accessControl_iam_role_list'
    },
    {
        name: 'policies',
        icon: 'policy',
        title: 'Policies',
        routerPageName: 'accessControl_iam_policy_list'
    }
] as Viewdata[]

const settingsViews = [
    {
        name: 'oauth',
        icon: 'icons/oauth.logo.svg',
        iconAsImage: true,
        title: 'OAuth',
        routerPageName: 'accessControl_settings_oauth'
    },
    {
        name: '2fa',
        icon: 'mdi-two-factor-authentication',
        title: '2FA',
        routerPageName: ''
    },
    {
        name: 'access_log',
        icon: 'event_note',
        title: 'Access log',
        routerPageName: ''
    },
    {
        name: 'vault',
        icon: 'lock',
        title: 'Vault and HSM',
        routerPageName: ''
    },
    {
        name: 'sentry_mode',
        icon: 'warning',
        title: 'Sentry mode',
        routerPageName: ''
    }
] as Viewdata[]

</script>

<style>
.no-minimums {
    min-width: 0px;
    min-height: 0px;
}
</style>