<template>
    <q-field outlined filled ref="fieldRef" v-model="props.modelValue" :label="props.label !== undefined ? props.label : $t('modules.iot.device.selectInput.label')" @focus.stop="startSelecting">
        <template v-slot:control>
            <div v-if="props.modelValue !== null" class="self-center full-width no-outline">{{ props.modelValue.name }}</div>
            <div v-else class="self-center full-width no-outline"></div>
        </template>
    </q-field>

    <q-dialog v-model="selectionDialogVisible">
        <div class="bg-primary" style="width: 800px; max-width: 90%;">
            <DeviceSelectorComponent
            :namespace="props.namespace"
            @selected="onSelected"
            @canceled="onSelectionCanceled"
        />
        </div>
    </q-dialog>
</template>

<script setup lang="ts">
import { Device } from 'src/boot/api/iot/device';
import { ref } from 'vue';

import DeviceSelectorComponent from './DeviceSelectorComponent.vue';
import { QField } from 'quasar';

const props = defineProps<{
    label?: string,
    namespace: string,
    modelValue: Device | null
}>()
const emits = defineEmits<{
    (e: 'update:modelValue', value: Device | null): void
}>()

const fieldRef = ref(null as null | QField)
const selectionDialogVisible = ref(false)

function startSelecting() {
    selectionDialogVisible.value = true
    fieldRef.value?.blur()
}

function onSelected(device: Device) {
    emits('update:modelValue', device)
    selectionDialogVisible.value = false
}

function onSelectionCanceled() {
    selectionDialogVisible.value = false
}

</script>