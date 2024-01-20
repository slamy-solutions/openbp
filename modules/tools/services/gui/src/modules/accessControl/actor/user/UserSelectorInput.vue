<template>
  <q-field
    ref="fieldRef"
    outlined
    filled
    v-model="props.modelValue"
    :label="
      props.label !== undefined
        ? props.label
        : $t('modules.accessControl.iam.actor.user.selectInput.label')
    "
    @focus="startSelecting"
  >
    <template v-slot:control>
      <div
        v-if="props.modelValue !== null"
        class="self-center full-width no-outline"
      >
        {{ props.modelValue.login }}
      </div>
      <div v-else class="self-center full-width no-outline"></div>
    </template>

    <q-dialog v-model="selectionDialogVisible">
      <div class="bg-primary q-pa-md" style="width: 900px; max-width: 90%">
        <UserSelectorComponent
          :namespace="props.namespace"
          @selected="onSelected"
          @cancelled="onSelectionCanceled"
        />
      </div>
    </q-dialog>
  </q-field>
</template>

<script setup lang="ts">
import { User } from "src/boot/api/accessControl/actor/user";
import { ref } from "vue";

import UserSelectorComponent from "./UserSelectorComponent.vue";
import { QField } from "quasar";

const props = defineProps<{
  label?: string;
  namespace: string;
  modelValue: User | null;
}>();
const emits = defineEmits<{
  (e: "update:modelValue", value: User | null): void;
}>();

const selectionDialogVisible = ref(false);
const fieldRef = ref(null as null | QField);

function startSelecting() {
  selectionDialogVisible.value = true;
  fieldRef.value?.blur();
}

function onSelected(user: User) {
  emits("update:modelValue", user);
  selectionDialogVisible.value = false;
}

function onSelectionCanceled() {
  selectionDialogVisible.value = false;
}
</script>
