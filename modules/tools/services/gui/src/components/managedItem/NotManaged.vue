<template>
    <q-item class="q-pa-none" :dense="props.dense">
    <q-item-section>
        <q-item-label>Not managed</q-item-label>
        <q-item-label caption lines="2">The resource management is unknown</q-item-label>
    </q-item-section>

    <q-item-section side top>
        <q-item-label caption>{{ $t('modules.accessControl.iam.policy.list.managedByIdentityColumn') }}</q-item-label>
    </q-item-section>
    </q-item>
</template>

<script setup lang="ts">
const props = defineProps<{
    dense?: boolean
}>()
</script>