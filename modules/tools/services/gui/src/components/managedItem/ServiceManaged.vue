<template>
    <q-item class="q-pa-none" :dense="props.dense">
    <q-item-section>
        <q-item-label>{{ props.manageData.service }} [Service]</q-item-label>
        <q-item-label caption lines="2">Reason: {{ props.manageData.reason }}</q-item-label>
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