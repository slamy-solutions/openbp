<template>
    <q-card flat>
        <q-card-section v-if="loadingError != ''">
            <div class="text-h5 text-negative">{{ $t('modules.accessControl.iam.auth.password.loadingError', { error: loadingError }) }}</div>
        </q-card-section>
        <q-card-section v-if="seted && loadingError == ''">
            <div class="text-h6 text-warning">{{ $t('modules.accessControl.iam.auth.password.header.enabled') }}</div>
            <div class="text-subtitle2">{{ $t('modules.accessControl.iam.auth.password.caption.enabled') }}</div>
        </q-card-section>
        <q-card-section v-if="!seted && loadingError == ''">
            <div class="text-h6 text-positive">{{ $t('modules.accessControl.iam.auth.password.header.disabled') }}</div>
            <div class="text-subtitle2">{{ $t('modules.accessControl.iam.auth.password.caption.disabled') }}</div>
        </q-card-section>
        <q-card-section v-if="loadingError == '' && editable">
            <q-input filled v-model="newPassword" label="New password" type="password" :disable="loadingError != ''"></q-input>
        </q-card-section>
        <q-card-actions class="row" v-if="loadingError == '' && editable">
            <q-btn
                :label="$t('modules.accessControl.iam.auth.password.disableButton')"
                class="col-xs-12 col-sm-4 col-md-2 q-mt-xs"
                color="dark"
                :disable="!seted || loading || loadingError != ''"
                @click="disablePassword"
                :loading="disabling"
            />
            <div class="col-xs-0 col-sm-1 col-md-4 q-mt-xs"></div>
            <q-btn
                :label="$t('modules.accessControl.iam.auth.password.setOrUpdateButton')"
                class="col-xs-12 col-sm-7 col-md-6 q-mt-xs"
                color="dark"
                :disable="loading || newPassword == '' || loadingError != ''"
                @click="updatePassword"
                :loading="updating"
            />
        </q-card-actions>
    
    </q-card>
    
    
</template>

<script lang="ts" setup>
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { onMounted, ref, watch } from 'vue';
import api from "../../../../boot/api";

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    identityNamespace: string,
    identityUUID: string,
    editable?: boolean
}>()

const newPassword = ref('')
const loading = ref(false)
const updating = ref(false)
const disabling = ref(false)

const loadingError = ref("")

const seted = ref(false)

async function loadPasswordInformation() {
    if (props.identityUUID == "") {
        return
    }
    
    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.auth.password.loadingOperationNotify')
    })
    loading.value = true
    loadingError.value = ""
    try {
        const status = await api.accessControl.auth.password.status({
            namespace: props.identityNamespace,
            identityUUID: props.identityUUID
        })
        seted.value = status.seted
        notif()
        loadingError.value = ""
    } catch (error) {
        loadingError.value = String(error)
        console.error(error)
        notif({
            type: 'negative',
            message: $i18n.t('modules.accessControl.iam.auth.password.loadingFailNotify', { error }),
            timeout: 5000
        })
    } finally {
        loading.value = false
    }
}

async function updatePassword() {
    loading.value = true
    updating.value = true
    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.auth.password.updateOperationNotify')
    })
    try {
        await api.accessControl.auth.password.setOrUpdate({
            namespace: props.identityNamespace,
            identityUUID: props.identityUUID,
            newPassword: newPassword.value
        })
        notif({
            type: 'positive',
            message: $i18n.t('modules.accessControl.iam.auth.password.updateSuccessNotify'),
            timeout: 3000
        })
        seted.value = true
    } catch (error) {
        console.error(error)
        notif({
            type: 'negative',
            message: $i18n.t('modules.accessControl.iam.auth.password.updateFailNotify', { error }),
            timeout: 5000
        })
    } finally {
        newPassword.value = ""
        loading.value = false
        updating.value = false
    }
}

async function disablePassword() {
    loading.value = true
    disabling.value = true
    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.auth.password.disableOperationNotify')
    })
    try {
        await api.accessControl.auth.password.disable({
            namespace: props.identityNamespace,
            identityUUID: props.identityUUID
        })
        notif({
            type: 'positive',
            message: $i18n.t('modules.accessControl.iam.auth.password.disableSuccessNotify'),
            timeout: 3000
        })
        seted.value = false
    } catch (error) {
        console.error(error)
        notif({
            type: 'negative',
            message: $i18n.t('modules.accessControl.iam.auth.password.disableFailNotify', { error }),
            timeout: 5000
        })
    } finally {
        loading.value = false
        disabling.value = false    }
}


onMounted(async () => await loadPasswordInformation())
watch(props, async (newProps)=>{
    if (newProps.identityUUID !== '') {
        await loadPasswordInformation()
    }
})

</script>