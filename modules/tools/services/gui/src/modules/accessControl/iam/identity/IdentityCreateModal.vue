<template>
    <div class="full-height full-width row justify-center items-center">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.accessControl.iam.identity.create.header') }}</h3>
              </q-card-section>
          <q-card-section>
              <q-form class="q-gutter-md">
              <q-input disable square filled v-model="props.namespace" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.identity.create.namespaceInput')" />
              <q-input square filled clearable v-model="name" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.identity.create.nameInput')" />
              <q-item tag="label" v-ripple>
                <q-item-section>
                  <q-item-label>{{ $t('modules.accessControl.iam.identity.create.initiallyActiveInput') }}</q-item-label>
                  <q-item-label caption>{{ $t('modules.accessControl.iam.identity.create.initiallyActiveInputCaption') }}</q-item-label>
                </q-item-section>
                <q-item-section side >
                  <q-toggle color="secondary" v-model="initiallyActive" />
                </q-item-section>
              </q-item>
              </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.accessControl.iam.identity.create.createButton')" :loading="loading" :disabled="loading" @click="createIdentity" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
              <p class="text-grey-6">{{ $t('modules.accessControl.iam.identity.create.createHint') }}</p>
          </q-card-section>
          </q-card>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import api from "../../../../boot/api";
import { useQuasar } from "quasar";
import { useI18n } from "vue-i18n";
import { Identity } from "../../../../boot/api/accessControl/identity";

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string
}>()

const emit = defineEmits<{
  (e: 'created', identity: Identity): void
}>()

const name = ref('')
const initiallyActive = ref(true)

const loading = ref(false)

async function createIdentity() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.identity.create.createOperationNotify')
  })
  try {
      const createdIdentity = await api.accessControl.identity.create({
        name: name.value,
        initiallyActive: initiallyActive.value,
        namespace: props.namespace
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.accessControl.iam.identity.create.createSuccessNotify'),
          timeout: 5000
      })
      name.value = ""
      emit('created', createdIdentity)
  } catch (error) {
      console.error(error)
      notif({
          type: 'negative',
          message: $i18n.t('modules.accessControl.iam.identity.create.createFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style>

</style>