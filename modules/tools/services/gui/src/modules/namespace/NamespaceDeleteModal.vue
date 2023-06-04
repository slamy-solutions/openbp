<template>
    <div class="window-height window-width row justify-center items-center">
      <div class="column">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1" style="border: red; border-style: solid; border-width: 10px;">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.namespace.delete.header') }}</h3>
              </q-card-section>
          <q-card-section>
            <span class="q-ma-sm text-bold">{{ $t('modules.namespace.delete.bodyText', {namespaceName: props.namespaceName}) }}</span>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="negative" size="lg" class="full-width" :label="$t('modules.namespace.delete.deleteButton')" :loading="loading" :disabled="loading" @click="deleteNamespace" />
          </q-card-actions>
          </q-card>
      </div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import api from "../../boot/api";
import { useQuasar } from "quasar";
import { useI18n } from "vue-i18n";

const $q = useQuasar()
const $i18n = useI18n()

const emit = defineEmits<{
  (e: 'deleted', name: string): void
}>()

const props = defineProps<{
    namespaceName: string
}>()

const loading = ref(false)

async function deleteNamespace() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.namespace.delete.deleteOperationNotify')
  })
  try {
      await api.namespace.list.delete({
        name: props.namespaceName
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.namespace.delete.deleteSuccessNotify'),
          timeout: 5000
      })
      emit('deleted', props.namespaceName)
  } catch (error) {
      notif({
          type: 'negative',
          message: $i18n.t('modules.namespace.delete.deleteFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style>

</style>