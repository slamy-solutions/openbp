<template>
    <q-item class="q-pa-none" :dense="props.dense">
    <q-item-section>
        <q-item-label>{{ props.manageData.identityNamespace }}</q-item-label>
        <q-item-label caption lines="2">{{ props.manageData.identityUUID }}</q-item-label>
    </q-item-section>

    <q-item-section side top>
        <q-item-label caption>{{ $t('modules.accessControl.iam.policy.list.managedByIdentityColumn') }}</q-item-label>
    </q-item-section>
    </q-item>
</template>

<script setup lang="ts">
import { IdentityManaged } from './model';

const props = defineProps<{
    manageData: IdentityManaged,
    dense?: boolean
}>()

</script>