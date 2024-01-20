<template>
    <div></div>
</template>

<script setup lang="ts">
import { ref, computed, defineEmits, defineExpose, watch } from 'vue';

const props = defineProps<{
    namespace: string,
    departmentUUID: string
}>()


const emits = defineEmits<{
    (e: 'onChanged'): void
}>()

const changed = ref(false)

watch(props, () => {
    changed.value = false
}, { deep: true })

async function applyChanges() {}

</script>