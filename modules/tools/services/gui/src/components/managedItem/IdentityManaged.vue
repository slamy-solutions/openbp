<template>
    <q-item class="q-pa-none" :dense="props.dense">
    <q-item-section>
        <q-item-label>{{ props.manageData.identityUUID }} [Identity]</q-item-label>
        <q-item-label caption lines="2">Identity namespace: {{ props.manageData.identityNamespace }}</q-item-label>
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