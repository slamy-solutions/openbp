<template>
    <ActionWarningModal 
        :action-button="$t('modules.iot.fleet.delete.deleteButton')"
        :body-text="$t('modules.iot.fleet.delete.bodyText', { uuid: props.uuid, namespace: props.namespace })"
        :header="$t('modules.iot.fleet.delete.header')"
        :loading="loading"
        @action-clicked="deleteFleet"
    />
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar';
import api from 'src/boot/api';
import ActionWarningModal from 'src/components/ActionWarningModal.vue';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const emit = defineEmits<{
  (e: 'deleted', namespace: string, uuid: string): void
}>()

const props = defineProps<{
    namespace: string,
    uuid: string
}>()

const $q = useQuasar()
const $i18n = useI18n()

const loading = ref(false)

async function deleteFleet() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.iot.fleet.delete.deleteOperationNotify')
  })
  try {
      await api.iot.fleet.deleteFleet({
        namespace: props.namespace,
        uuid: props.uuid
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.iot.fleet.user.delete.deleteSuccessNotify')
      })
      emit('deleted', props.namespace, props.uuid)
  } catch (error) {
      notif({
          type: 'negative',
          message: $i18n.t('modules.iot.fleet.user.delete.deleteFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  } 
}

</script>