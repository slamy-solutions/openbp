<template>
    <q-field ref="fieldRef" outlined filled v-model="props.modelValue" :label="props.label !== undefined ? props.label : $t('modules.iot.fleet.selectInput.label')" clearable @focus="startSelecting">
        <template v-slot:control>
            <div v-if="props.modelValue !== null" class="self-center full-width no-outline">{{ props.modelValue.name }}</div>
            <div v-else class="self-center full-width no-outline"></div>
        </template>
    </q-field>

    <q-dialog v-model="selectionDialogVisible">
        <div class="bg-primary" style="width: 900px; max-width: 90%;">
        <FleetSelectorComponent
            :namespace="props.namespace"
            @selected="onSelected"
            @canceled="onSelectionCanceled"
        />
    </div>
    </q-dialog>
</template>

<script setup lang="ts">
import { Fleet } from 'src/boot/api/iot/fleet';
import { ref } from 'vue';

import FleetSelectorComponent from './FleetSelectorComponent.vue';
import { QField } from 'quasar';

const props = defineProps<{
    label?: string,
    namespace: string,
    modelValue: Fleet | null
}>()
const emits = defineEmits<{
    (e: 'update:modelValue', value: Fleet | null): void
}>()

const selectionDialogVisible = ref(false)
const fieldRef = ref(null as null | QField)

function startSelecting() {
    selectionDialogVisible.value = true
    fieldRef.value?.blur()
}

function onSelected(fleet: Fleet) {
    emits('update:modelValue', fleet)
    selectionDialogVisible.value = false
}

function onSelectionCanceled() {
    selectionDialogVisible.value = false
}

</script>