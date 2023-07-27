<template>
    <div class="row justify-center items-center full-width full-height">
      <div class="column">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.accessControl.iam.role.create.header') }}</h3>
              </q-card-section>
          <q-card-section>
              <q-form class="q-gutter-md">
              <q-input disable square filled v-model="props.namespace" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.role.create.namespaceInput')" />
              <q-input square filled clearable v-model="name" counter maxlength="64" type="text" :label="$t('modules.accessControl.iam.role.create.nameInput')" />
              <q-input square filled clearable v-model="description" counter maxlength="256" type="textarea" :label="$t('modules.accessControl.iam.role.create.descriptionInput')" />
              </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.accessControl.iam.role.create.createButton')" :loading="loading" :disabled="loading" @click="createRole" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
              <p class="text-grey-6">{{ $t('modules.accessControl.iam.role.create.createHint') }}</p>
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
import { Role } from "src/boot/api/accessControl/role";

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string
}>()

const emit = defineEmits<{
  (e: 'created', role: Role): void
}>()

const name = ref('')
const description = ref('')
const namespaceIndependent = ref(false)

const loading = ref(false)

async function createRole() {
  loading.value = true
  const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.role.create.createOperationNotify')
  })
  try {
      const createdRole = await api.accessControl.role.create({
        name: name.value,
        namespace: props.namespace,
        description: description.value,
      })
      notif({
          type: 'positive',
          message: $i18n.t('modules.accessControl.iam.role.create.createSuccessNotify'),
          timeout: 5000
      })
      name.value = ""
      emit('created', createdRole)
  } catch (error) {
      console.log(error)
      notif({
          type: 'negative',
          message: $i18n.t('modules.accessControl.iam.role.create.createFailNotify', { error }),
          timeout: 5000
      })
  } finally {
      loading.value = false
  }   
}

</script>

<style>

</style>