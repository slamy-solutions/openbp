<template>
    <q-list bordered separator>
      <RoleListItemComponent
          v-for="item in props.roles" :key="item.namespace + ':' + item.uuid"
          :namespace="item.namespace"
          :uuid="item.uuid"
          :editable="editable"
          class="q-pa-xs"
          @on-delete-clicked="emit('removed', item.namespace, item.uuid)"
      />

      <q-item class="justify-center items-center" v-if="props.editable">
        <q-btn dense round color="positive" icon="add" @click="addNewDialog = true"/>
      </q-item>
    </q-list>
    <q-dialog v-model="addNewDialog">
      <RoleSelectorComponent :namespace="props.namespace" @canceled="onAddCancel" @selected="onAddSelected" class="bg-primary"/>
    </q-dialog>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
import RoleListItemComponent from './RoleListItemComponent.vue';
import RoleSelectorComponent from './RoleSelectorComponent.vue';
      const props = defineProps<{
          namespace: string,
          roles: Array<{ namespace: string, uuid: string }>,
          editable?: boolean
      }>()

      const emit = defineEmits<{
        (e: 'added', namespace: string, uuid: string): void,
        (e: 'removed', namespace: string, uuid: string): void
      }>()

      const addNewDialog = ref(false)

      function onAddCancel() {
        addNewDialog.value = false
      }
      function onAddSelected(namespace: string, uuid: string) {
        emit("added", namespace, uuid)
        addNewDialog.value = false
      }
  </script>
  
  <style>
  
  </style>