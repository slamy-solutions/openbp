<template>
    <div class="full-height row justify-center items-center card-style">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.iot.integration.balena.server.add.header') }}</h3>
              </q-card-section>
          <q-card-section>
              <q-form class="q-gutter-md" ref="inputForm">
              <q-input disable square filled v-model="props.namespace" counter maxlength="32" type="text" :label="$t('modules.iot.integration.balena.server.add.namespaceInput')" />
              <q-input square filled clearable v-model="name" counter maxlength="64" type="text" :label="$t('modules.iot.integration.balena.server.add.nameInput')" />
              <q-input square filled clearable v-model="description" counter maxlength="256" type="text" :label="$t('modules.iot.integration.balena.server.add.descriptionInput')" />
              <q-input square filled clearable v-model="url" counter maxlength="128" type="text" :label="$t('modules.iot.integration.balena.server.add.urlInput')" shadow-text="https://api.balena-cloud.com"/>
              <q-input square filled clearable v-model="apiToken" counter maxlength="1024" type="password" :label="$t('modules.iot.integration.balena.server.add.apiTokenInput')" />
              </q-form>
          </q-card-section>
          <q-card-section>
            <q-btn square outlined :loading="loading" :label="$t('modules.iot.integration.balena.server.add.validateConnectionButton')" @click="validateConnection" class="fit" :disable="apiToken == '' || url == ''"></q-btn>
            <p class="text-subtitle text-info text-center q-mt-sm" v-if="!creadentialsValidated">{{ $t('modules.iot.integration.balena.server.add.connectionDataNotValidated') }}</p> 
            <p class="text-subtitle text-negative text-center q-mt-sm" v-if="credentialsStatus != 'OK'">{{ credentialsStatus + ": " + credentialsMesssage }}</p>
            <p class="text-subtitle text-positive text-center q-mt-sm" v-if="creadentialsValidated && credentialsStatus == 'OK'">{{ $t('modules.iot.integration.balena.server.add.connectionDataIsValid') }}</p> 
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.iot.integration.balena.server.add.createButton')" :loading="loading" :disable="loading || !creadentialsValidated || credentialsStatus != 'OK'" @click="addServer" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
              <p class="text-grey-6">{{ $t('modules.iot.integration.balena.server.add.createHint') }}</p>
          </q-card-section>
          </q-card>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import api from "../../../../boot/api";
import { QForm, useQuasar } from "quasar";
import { useI18n } from "vue-i18n";
import { Device } from "src/boot/api/iot/device";
import { ConnectionStatus } from "src/boot/api/iot/integration/balena/tools";
import { Server } from "src/boot/api/iot/integration/balena/models";

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string
}>()

const emit = defineEmits<{
  (e: 'added', server: Server): void
}>()

const inputForm = ref<QForm | null>(null)
const name = ref('')
const description = ref('')
const url = ref('https://api.balena-cloud.com')
const apiToken = ref('')

const creadentialsValidated = ref(false)
const credentialsStatus = ref('OK' as ConnectionStatus)
const credentialsMesssage = ref('')

const loading = ref(false)

async function validateConnection() {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.iot.integration.balena.server.add.validateOperationNotify')
    })
    try {
        const validateCredentialsResponse = await api.iot.integration.balena.tools.verifyConnectionData({
            apiToken: apiToken.value,
            url: url.value,
        })
        creadentialsValidated.value = true
        credentialsStatus.value = validateCredentialsResponse.status
        credentialsMesssage.value = validateCredentialsResponse.message
        notif()
    } catch (error) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.iot.integration.balena.server.add.validateFailNotify', { error }),
            timeout: 5000
        })
    } finally {
        loading.value = false
    }
}

async function addServer() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.iot.integration.balena.server.add.addOperationNotify')
  })
  try {
      const addServerResponse = await api.iot.integration.balena.server.create({
        authToken: apiToken.value,
        baseUrl: url.value,
        namespace: props.namespace,
        name: name.value,
        description: description.value,
      })

      notif({
          type: 'positive',
          message: $i18n.t('modules.iot.integration.balena.server.add.addSuccessNotify'),
          timeout: 5000
      })
      emit('added', addServerResponse.server)
  } catch (error) {
      console.error(error)
      notif({
          type: 'negative',
          message: $i18n.t('modules.iot.integration.balena.server.add.addFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style scoped>
.card-style {
  width: 800px;
  max-width: 90%;
}
</style>