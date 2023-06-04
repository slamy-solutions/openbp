<template>
    <q-item class="q-pa-none" :dense="props.dense">
    <q-item-section>
        <q-item-label>{{ props.manageData.service }}</q-item-label>
        <q-item-label caption lines="2">{{ props.manageData.reason }}</q-item-label>
    </q-item-section>

    <q-item-section side top>
        <q-item-label caption>{{ $t('modules.accessControl.iam.policy.list.managedByServiceColumn') }}</q-item-label>
        <q-item-label caption>{{ props.manageData.managementId }}</q-item-label>
    </q-item-section>
    </q-item>
</template>

<script setup lang="ts">
import { ServiceManaged } from './model';

const props = defineProps<{
    manageData: ServiceManaged,
    dense?: boolean
}>()

</script>