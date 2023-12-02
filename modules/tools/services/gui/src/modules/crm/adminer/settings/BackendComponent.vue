<template>
    <q-card flat bordered class="bg-transparent">
        <q-card-section>
            <div class="text-h6">{{ $t('modules.crm.adminer.settings.backend.title') }}</div>
            <div class="text-caption">{{ $t('modules.crm.adminer.settings.backend.caption') }}</div> 
        </q-card-section>

        <q-card-section v-if="loading">
            <q-spinner size="xl" class="text-center allign-center full-width"></q-spinner>
        </q-card-section>

        <q-card-section v-if="!loading">
            <q-select filled  v-model="backend" class="fit" :label="$t('modules.crm.adminer.settings.backend.backendSelect')" :options="['NATIVE', 'ONE_C']" :disable="updating" @update:model-value="settingsModifiedByUser" />
        </q-card-section>

        <q-card-section v-if="!loading && backend == 'NATIVE'">
            <div>{{ $t('modules.crm.adminer.settings.backend.native.caption') }}</div>
        </q-card-section>

        <q-card-section v-if="!loading && backend == 'ONE_C'">
            <div>{{ $t('modules.crm.adminer.settings.backend.onec.caption') }}</div>
        </q-card-section>
        <q-card-section v-if="!loading && backend == 'ONE_C'" class="q-gutter-sm">
            <q-input class="fit q-mr-md" filled v-model="backendURL" :label="$t('modules.crm.adminer.settings.backend.onec.urlInput')" :disabled="updating" @update:model-value="settingsModifiedByUser" />
            <q-input class="fit" filled v-model="backendToken" type="password" :label="$t('modules.crm.adminer.settings.backend.onec.tokenInput')" :disabled="updating" @update:model-value="settingsModifiedByUser" />
            <q-btn class="fit" v-if="!oneCConnectionValid" outline @click="checkOneCConnection" :label="$t('modules.crm.adminer.settings.backend.onec.checkConnectionButton')" :disable="updating"></q-btn>
            <div v-if="oneCConnectionError != '' && !oneCConnectionValid" class="text-negative text-center text-bold">{{ `[${oneCConnectionStatusCode}]: ${oneCConnectionError}` }}</div>
            <div v-if="oneCConnectionError != '' && oneCConnectionValid" class="text-positive text-center">{{ $t('modules.crm.adminer.settings.backend.onec.connectionValid') }}</div>
        </q-card-section>

        <q-card-actions v-if="settingsModified && (backend !== 'ONE_C' || oneCConnectionValid)">
            <q-btn outline @click="updateSettings" class="fit q-ml-sm" :disabled="updating" :loading="updating" :label="$t('modules.crm.adminer.settings.backend.updateButton')" />
        </q-card-actions>
    </q-card>
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { onMounted, ref } from 'vue';
import { BackendType } from '../../../../boot/api/crm/settings'
import { api } from '../../../../boot/api'

const props = defineProps<{
    namespaceName: string
}>()
const $q = useQuasar()
const $i18n = useI18n()


const loading = ref(false)
const updating = ref(false)
const backend = ref<BackendType>('NATIVE')
const backendURL = ref('')
const backendToken = ref('')
const settingsModified = ref(false)

const oneCConnectionValid = ref(false)
const oneCConnectionError = ref('')
const oneCConnectionStatusCode = ref(0)

async function loadSettings() {
    loading.value = true
    try {
        const response = await api.crm.settings.get({namespace: props.namespaceName})
        backend.value = response.settings.backendType

        if (response.settings.backendType == 'ONE_C') {
            backendURL.value = response.settings.backendURL
            backendToken.value = response.settings.token
        }
    } catch (error) {
        $q.notify({
            type: 'negative',
            message: $i18n.t('modules.crm.adminer.settings.backend.loadSettingsFailNotify', { error }),
            timeout: 5000
        })
        console.error(error)
    } finally {
        loading.value = false
    }
}

async function updateSettings() {
    updating.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.crm.adminer.settings.backend.updateOperationNotify'),
    })

    try {
        await api.crm.settings.update({
            namespace: props.namespaceName,
            settings: {
                backendType: backend.value,
                backendURL: backendURL.value,
                token: backendToken.value
            }
        })
        notif({
            type: 'positive',
            message: $i18n.t('modules.crm.adminer.settings.backend.updateSuccessNotify'),
            timeout: 5000
        })
        settingsModified.value = false
    } catch (error) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.crm.adminer.settings.backend.updateFailNotify', { error }),
            timeout: 5000
        })
        console.error(error)
    } finally {
        updating.value = false
    }
}

function settingsModifiedByUser() {
    settingsModified.value = true
    oneCConnectionValid.value = false
}

async function checkOneCConnection() {
    try {
        const response = await api.crm.settings.checkOneCConnection({
            backendURL: backendURL.value,
            token: backendToken.value
        })
        oneCConnectionValid.value = response.success
        oneCConnectionError.value = response.message
        oneCConnectionStatusCode.value = response.statusCode
    } catch (e) {
        oneCConnectionValid.value = false
        oneCConnectionError.value = String(e)
    
    }
    
}

onMounted(async () => {
    await loadSettings()
})

</script>