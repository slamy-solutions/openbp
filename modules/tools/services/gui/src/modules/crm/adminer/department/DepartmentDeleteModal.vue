<template>
    <div class="window-height window-width row justify-center items-center">
      <div class="column">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1" style="border: red; border-style: solid; border-width: 10px;">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.crm.adminer.department.delete.title') }}</h3>
              </q-card-section>
          <q-card-section>
            <span class="q-ma-sm text-bold">{{ $t('modules.crm.adminer.department.delete.bodyText', {namespace: props.namespace, uuid: props.departmentUUID, name: props.name}) }}</span>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="negative" size="lg" class="full-width" :label="$t('modules.crm.adminer.department.delete.deleteButton')" :loading="loading" :disabled="loading" @click="deleteDepartment" />
          </q-card-actions>
          </q-card>
      </div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import api from "../../../../boot/api";
import { useQuasar } from "quasar";
import { useI18n } from "vue-i18n";

const $q = useQuasar()
const $i18n = useI18n()

const emit = defineEmits<{
  (e: 'deleted', namespace: string, uuid: string): void
}>()

const props = defineProps<{
    name: string,
    namespace: string,
    departmentUUID: string
}>()

const loading = ref(false)

async function deleteDepartment() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.crm.adminer.department.delete.deleteOperationNotify')
  })
  try {
      await api.crm.department.delete({
        namespace: props.namespace,
        uuid: props.departmentUUID
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.crm.adminer.department.delete.deleteSuccessNotify'),
          timeout: 5000
      })
      emit('deleted', props.namespace, props.departmentUUID)
  } catch (error) {
      notif({
          type: 'negative',
          message: $i18n.t('modules.crm.adminer.department.delete.deleteFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style>

</style>