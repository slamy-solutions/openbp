<template>
    <div class="full-height full-width row justify-center items-center">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.iot.fleet.create.header') }}</h3>
              </q-card-section>
          <q-card-section>
              <q-form class="q-gutter-md">
              <q-input disable square filled v-model="props.namespace" counter maxlength="32" type="text" :label="$t('modules.iot.fleet.create.namespaceInput')" />
              <q-input square filled clearable v-model="name" counter maxlength="32" type="text" :label="$t('modules.iot.fleet.create.nameInput')" />
              <q-input square filled clearable v-model="description" counter maxlength="128" type="text" :label="$t('modules.iot.fleet.create.descriptionInput')" />
              </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.iot.fleet.create.createButton')" :loading="loading" :disabled="loading" @click="createFleet" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
              <p class="text-grey-6">{{ $t('modules.iot.fleet.create.createHint') }}</p>
          </q-card-section>
          </q-card>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import api from "../../../boot/api";
import { useQuasar } from "quasar";
import { useI18n } from "vue-i18n";
import { Fleet } from "src/boot/api/iot/fleet";

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string
}>()

const emit = defineEmits<{
  (e: 'created', fleet: Fleet): void
}>()

const name = ref('')
const description = ref('')

const loading = ref(false)

async function createFleet() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.iot.fleet.create.createOperationNotify')
  })
  try {
      const createFleetResponse = await api.iot.fleet.createFleet({
        namespace: props.namespace,
        name: name.value,
        description: description.value
      })

      notif({
          type: 'positive',
          message: $i18n.t('modules.iot.fleet.create.createSuccessNotify'),
          timeout: 5000
      })
      emit('created', createFleetResponse.fleet)
  } catch (error) {
      console.error(error)
      notif({
          type: 'negative',
          message: $i18n.t('modules.iot.fleet.create.createFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style>

</style>