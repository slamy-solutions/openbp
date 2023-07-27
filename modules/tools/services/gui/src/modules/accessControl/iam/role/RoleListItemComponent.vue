<template>
    <q-item :clickable="!props.editable" :v-ripple="!props.editable" @click="detailsDialog = true">
          <q-item-section>
              <q-item-label>{{ displayName }}</q-item-label>
              <q-item-label caption>{{ displayCaption }}</q-item-label>
          </q-item-section>
          <q-item-section side top v-if="props.editable">
              <q-btn rounded icon="delete" size="md" color="negative" @click="emit('onDeleteClicked')"></q-btn>
          </q-item-section>

          <q-inner-loading :showing="loading">
            <q-spinner-gears size="sm" color="dark"/>
          </q-inner-loading>
      </q-item>
      <q-dialog v-model="detailsDialog">
        <div class="bg-primary" style="max-width: 95%; width: 95%;">
            <RoleViewComponent
                :namespace="props.namespace"
                :uuid="props.uuid"
                :update-possible="true"
            />
        </div>
      </q-dialog>
</template>

<script setup lang="ts">
import api from 'src/boot/api';
import { Role } from 'src/boot/api/accessControl/role';
import RoleViewComponent from './RoleViewComponent.vue';
import { Ref, computed, onMounted, reactive, ref, watch } from 'vue';
    const props = defineProps<{
        namespace: string,
        uuid: string,
        editable?: boolean
    }>()

    const emit = defineEmits<{
        (e: 'onDeleteClicked'): void,
      }>()

    const loading = ref(true)
    const error = ref("")
    const role = ref(null as Role | null)
    const detailsDialog = ref(false)

    const displayName = computed(()=>{
        if (role.value) {
            return `${role.value.name} [${props.uuid}]`
        }

        return props.uuid
    })

    const displayCaption = computed(()=>{
        if (!loading.value) {
            if (error.value != "") {
                return error.value
            }
            return role.value?.description
        }
    })

    async function loadRoleInformation() {
        try {
            role.value = await api.accessControl.role.get({ namespace: props.namespace, uuid: props.uuid })
        } catch(e) {
            error.value = String(e)
        } finally {
            loading.value = false
        }
    }

    onMounted(loadRoleInformation)
</script>

<style>

</style>