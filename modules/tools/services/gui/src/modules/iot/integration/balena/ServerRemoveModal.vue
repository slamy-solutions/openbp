<template>
    <ActionWarningModal 
        :action-button="$t('modules.iot.integration.balena.server.remove.deleteButton')"
        :body-text="$t('modules.iot.integration.balena.server.remove.bodyText', { uuid: props.uuid })"
        :header="$t('modules.iot.integration.balena.server.remove.header')"
        :loading="loading"
        @action-clicked="removeServer"
    />
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar';
import api from 'src/boot/api';
import ActionWarningModal from 'src/components/ActionWarningModal.vue';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const emit = defineEmits<{
  (e: 'removed', uuid: string): void
}>()

const props = defineProps<{
    namespace: string,
    uuid: string
}>()

const $q = useQuasar()
const $i18n = useI18n()

const loading = ref(false)

async function removeServer() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.iot.integration.balena.server.remove.deleteOperationNotify')
  })
  try {
      await api.iot.integration.balena.server.delete({
        uuid: props.uuid
      })
      notif({
          type: 'positive',
          message: $i18n.t('modulesiot.integration.balena.server.remove.deleteSuccessNotify')
      })
      emit('removed', props.uuid)
  } catch (error) {
      notif({
          type: 'negative',
          message: $i18n.t('modulesiot.integration.balena.server.remove.deleteFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  } 
}

</script>