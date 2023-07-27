<template>
    <q-item class="q-pa-none" :dense="props.dense">
    <q-item-section>
        <q-item-label>BUILT IN</q-item-label>
        <q-item-label caption lines="2">The resource is built in</q-item-label>
    </q-item-section>
    </q-item>
</template>

<script setup lang="ts">
import { BuiltInManaged } from './model';

const props = defineProps<{
    manageData: BuiltInManaged,
    dense?: boolean
}>()

</script>