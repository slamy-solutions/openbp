<template>
      <div class="window-height window-width row justify-center items-center">
        <div class="column">
        <div class="row">
            <q-card square bordered class="q-pa-lg shadow-1">
                <q-card-section class="text-center">
                    <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.namespace.create.header') }}</h3>
                </q-card-section>
            <q-card-section>
                <q-form class="q-gutter-md">
                <q-input square filled clearable v-model="name" counter maxlength="32" type="text" :label="$t('modules.namespace.create.nameInput')" />
                <q-input square filled clearable v-model="fullName" counter maxlength="128" type="text" :label="$t('modules.namespace.create.fullNameInput')" />
                <q-input square filled clearable v-model="description" counter maxlength="512" type="textarea" :label="$t('modules.namespace.create.descriptionInput')" />
                </q-form>
            </q-card-section>
            <q-card-actions class="q-px-md">
                <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.namespace.create.createButton')" :loading="loading" :disabled="loading" @click="createNamespace" />
            </q-card-actions>
            <q-card-section class="text-center q-pa-none">
                <p class="text-grey-6">{{ $t('modules.namespace.create.createHint') }}</p>
            </q-card-section>
            </q-card>
        </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Namespace } from "../../boot/api/namespace/models";
import { ref } from "vue"
import api from "../../boot/api";
import { useQuasar } from "quasar";
import { useI18n } from "vue-i18n";

const $q = useQuasar()
const $i18n = useI18n()

const emit = defineEmits<{
    (e: 'created', namespace: Namespace): void
}>()

const name = ref('')
const fullName = ref('')
const description = ref('')

const loading = ref(false)

async function createNamespace() {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.namespace.create.createOperationNotify')
    })
    try {
        const createdNamespace = await api.namespace.list.create({
            name: name.value,
            fullName: fullName.value,
            description: description.value
        })
        notif({
            type: 'positive',
            message: $i18n.t('modules.namespace.create.createSuccessNotify'),
            timeout: 5000
        })
        name.value = ""
        fullName.value = ""
        description.value = ""
        emit('created', createdNamespace)
    } catch (error) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.namespace.create.createFailNotify', { error }),
            timeout: 5000
        })
    } finally {
        loading.value = false
    }   
}

</script>

<style>

</style>