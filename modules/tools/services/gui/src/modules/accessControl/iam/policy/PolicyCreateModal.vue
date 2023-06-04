<template>
    <div class="row justify-center items-center full-width full-height">
      <div class="column">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.accessControl.iam.policy.create.header') }}</h3>
              </q-card-section>
          <q-card-section>
              <q-form class="q-gutter-md">
              <q-input disable square filled v-model="props.namespace" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.policy.create.namespaceInput')" />
              <q-input square filled clearable v-model="name" counter maxlength="64" type="text" :label="$t('modules.accessControl.iam.policy.create.nameInput')" />
              <q-input square filled clearable v-model="description" counter maxlength="256" type="textarea" :label="$t('modules.accessControl.iam.policy.create.descriptionInput')" />
              <q-item tag="label" v-ripple>
                <q-item-section>
                  <q-item-label>{{ $t('modules.accessControl.iam.policy.create.namespaceIndependentInput') }}</q-item-label>
                  <q-item-label caption>{{ $t('modules.accessControl.iam.policy.create.namespaceIndependentInputCaption') }}</q-item-label>
                </q-item-section>
                <q-item-section side >
                  <q-toggle color="secondary" v-model="namespaceIndependent" />
                </q-item-section>
              </q-item>
              <h6 class="q-ma-sm text-uppercase q-pa-none">
                {{ $t('modules.accessControl.iam.policy.create.resourcesList.header') }}
                <q-icon name="help" size="20px" class="q-mb-xs" >
                  <q-tooltip>{{ $t('modules.accessControl.iam.policy.create.resourcesList.caption') }}</q-tooltip>
                </q-icon>
              </h6>
              <StringListInputComponent v-model:value="resources" class="q-mt-none"/>
              <h6 class="q-ma-sm text-uppercase q-pa-none">
                {{ $t('modules.accessControl.iam.policy.create.actionsList.header') }}
                <q-icon name="help" size="20px" class="q-mb-xs" >
                  <q-tooltip>{{ $t('modules.accessControl.iam.policy.create.actionsList.caption') }}</q-tooltip>
                </q-icon>
              </h6>
              <StringListInputComponent v-model:value="actions" class="q-mt-none"/>
              </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.accessControl.iam.policy.create.createButton')" :loading="loading" :disabled="loading" @click="createPolicy" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
              <p class="text-grey-6">{{ $t('modules.accessControl.iam.policy.create.createHint') }}</p>
          </q-card-section>
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

import StringListInputComponent from "./StringListInputComponent.vue"
import { Policy } from "../../../../boot/api/accessControl/policy";

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string
}>()

const emit = defineEmits<{
  (e: 'created', policy: Policy): void
}>()

const name = ref('')
const description = ref('')
const namespaceIndependent = ref(false)
const actions = ref([] as Array<{id: string, value: string}>)
const resources = ref([] as Array<{id: string, value: string}>)

const loading = ref(false)

async function createPolicy() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.policy.create.createOperationNotify')
  })
  try {
      const createdPolicy = await api.accessControl.policy.create({
        name: name.value,
        namespace: props.namespace,
        description: description.value,
        namespaceIndependent: namespaceIndependent.value,
        actions: actions.value.map((v) => v.value),
        resources: resources.value.map((v) => v.value)
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.accessControl.iam.policy.create.createSuccessNotify'),
          timeout: 5000
      })
      name.value = ""
      emit('created', createdPolicy)
  } catch (error) {
      console.log(error)
      notif({
          type: 'negative',
          message: $i18n.t('modules.accessControl.iam.policy.create.createFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style>

</style>