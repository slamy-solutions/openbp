<template>
    <div class="full-height full-width row justify-center items-center">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.accessControl.iam.actor.user.create.header') }}</h3>
              </q-card-section>
          <q-card-section>
              <q-form class="q-gutter-md">
              <q-input disable square filled v-model="props.namespace" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.actor.user.create.namespaceInput')" />
              <q-input square filled clearable v-model="login" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.actor.user.create.loginInput')" />
              <q-input square filled clearable v-model="fullName" counter maxlength="64" type="text" :label="$t('modules.accessControl.iam.actor.user.create.fullNameInput')" />
              <q-input square filled clearable v-model="email" counter maxlength="64" type="text" :label="$t('modules.accessControl.iam.actor.user.create.emailInput')" />
              </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.accessControl.iam.actor.user.create.createButton')" :loading="loading" :disabled="loading" @click="createIdentity" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
              <p class="text-grey-6">{{ $t('modules.accessControl.iam.actor.user.create.createHint') }}</p>
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
import { User } from "src/boot/api/accessControl/actor/user";

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string
}>()

const emit = defineEmits<{
  (e: 'created', identity: User): void
}>()

const login = ref('')
const fullName = ref('')
const email = ref('')
const initiallyActive = ref(true)

const loading = ref(false)

async function createIdentity() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.actor.iam.user.create.createOperationNotify')
  })
  try {
      const createdUser = await api.accessControl.actor.user.create({
        login: login.value,
        email: email.value,
        fullName: fullName.value,
        namespace: props.namespace
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.accessControl.actor.iam.user.create.createSuccessNotify'),
          timeout: 5000
      })
      emit('created', createdUser)
  } catch (error) {
      console.error(error)
      notif({
          type: 'negative',
          message: $i18n.t('modules.accessControl.actor.iam.user.create.createFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style>

</style>