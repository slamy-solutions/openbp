<template>
    <ActionWarningModal 
        :action-button="$t('modules.crm.adminer.performer.delete.deleteButton')"
        :body-text="$t('modules.crm.adminer.performer.delete.bodyText', { name: props.performer.name })"
        :header="$t('modules.crm.adminer.performer.delete.title')"
        :loading="loading"
        @action-clicked="deletePerformer"
    />
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar';
import api from 'src/boot/api';
import { Performer } from 'src/boot/api/crm/performer';
import ActionWarningModal from 'src/components/ActionWarningModal.vue';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const emit = defineEmits<{
  (e: 'deleted', performer: Performer): void
}>()

const props = defineProps<{
    namespace: string,
    performer: Performer
}>()

const $q = useQuasar()
const $i18n = useI18n()

const loading = ref(false)

async function deletePerformer() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.crm.adminer.performer.delete.deleteOperationNotify')
  })
  try {
    console.log(props)
      await api.crm.performer.delete({
        namespace: props.namespace,
        uuid: props.performer.uuid
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.crm.adminer.performer.delete.deleteSuccessNotify')
      })
      emit('deleted', props.performer as Performer)
  } catch (error) {
      notif({
          type: 'negative',
          message: $i18n.t('modules.crm.adminer.performer.delete.deleteFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  } 
}

</script>