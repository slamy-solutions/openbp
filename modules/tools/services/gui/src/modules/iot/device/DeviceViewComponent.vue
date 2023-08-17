<template>

</template>

<script setup lang="ts">
import api from 'src/boot/api';

const props = defineProps<{
    namespace: string,
    uuid: string
}>()

async function loadData() {
    const response  = await api.iot.device.getDevice({ namespace: props.namespace, uuid: props.uuid })
}

</script>