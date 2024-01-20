<template>
  <q-field
    ref="fieldRef"
    outlined
    filled
    v-model="props.modelValue"
    :label="
      props.label !== undefined
        ? props.label
        : $t('modules.crm.adminer.department.selectInput.label')
    "
    @focus="startSelecting"
  >
    <template v-slot:control>
      <div
        v-if="props.modelValue !== null"
        class="self-center full-width no-outline"
      >
        {{ props.modelValue.name }}
      </div>
      <div v-else class="self-center full-width no-outline"></div>
    </template>

    <q-dialog v-model="selectionDialogVisible">
      <div class="bg-primary q-pa-md" style="width: 900px; max-width: 90%;">
        <DepartmentSelectorComponent
          :namespace="props.namespace"
          @selected="onSelected"
          @cancelled="onSelectionCanceled"
        />
      </div>
    </q-dialog>
  </q-field>
</template>

<script setup lang="ts">
import { Department } from "src/boot/api/crm/department";
import { ref } from "vue";

import DepartmentSelectorComponent from "./DepartmentSelectorComponent.vue";
import { QField } from "quasar";

const props = defineProps<{
  label?: string;
  namespace: string;
  modelValue: Department | null;
}>();
const emits = defineEmits<{
  (e: "update:modelValue", value: Department | null): void;
}>();

const selectionDialogVisible = ref(false);
const fieldRef = ref(null as null | QField);

function startSelecting() {
  selectionDialogVisible.value = true;
  fieldRef.value?.blur();
}

function onSelected(department: Department) {
  emits("update:modelValue", department);
  selectionDialogVisible.value = false;
}

function onSelectionCanceled() {
  selectionDialogVisible.value = false;
}
</script>
