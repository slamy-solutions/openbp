<template>
  <q-card class="full-width full-height">
    <q-card-section>
      <div class="text-h6 gt-sm">Access management</div>
      <div class="text-subtitle2"></div>
    </q-card-section>

    <q-card-section>
        <div>Actors</div>
        <q-list>
            <q-separator />
            <div v-for="view in actorViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName"/>
                <q-separator />
            </div>
        </q-list>
    </q-card-section>

    <q-card-section>
        <div>IAM</div>
        <q-list>
            <q-separator />
            <div v-for="view in iamViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName"/>
                <q-separator />
            </div>
        </q-list>
    </q-card-section>

    <q-card-section>
        <div>Settings and security</div>
        <q-list>
            <q-separator />
            <div v-for="view in settingsViews" :key="view.name">
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
        name: 'sso',
        icon: 'local_activity',
        title: 'SSO and 2FA',
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