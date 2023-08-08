<template>
<h6 class="row q-ma-md" v-if="fleetSelected">
    <h6 class="col-6 q-mt-none q-mb-none">{{ $t('modules.iot.fleet.view.header') }}</h6>
    <div class="col-6 text-right">
        <q-btn color="dark" outline label="" icon="menu">
                <q-menu>
                <q-list>
                    <q-item clickable v-close-popup class="q-pa-xs" v-if="!updatesEnabled && updatePossible">
                    <q-btn v-close-popup label="Enable updates"  outline color="dark" size="sm" class="full-width" @click="updatesEnabled = true"></q-btn>
                    </q-item>
                    <q-item clickable v-close-popup class="q-pa-xs">
                    <q-btn v-close-popup label="Delete" outline color="negative" size="sm" class="full-width"></q-btn>
                    </q-item>
                </q-list>
                </q-menu>
            </q-btn>
    </div>
</h6>
<div class="row q-ma-xl full-width" v-if="!fleetSelected">
    <span class="col text-center text-h5 text-bold">
        {{ $t('modules.iot.fleet.view.notSelected') }}
    </span>
</div>
<div class="row q-ma-xl full-width" v-if="fleetSelected && loadingError !== ''">
    <span class="col text-center text-h4 text-bold text-negative">
        {{ $t('modules.iot.fleet.view.error', { error: loadingError }) }}
    </span>
</div>
<div
    class="row full-width"
    v-if="fleetSelected && loadingError === ''"
>

    <div
        class="col-sm-12 col-md-4 q-pl-md q-pr-md q-pb-md"
    >
    <q-list dense>
        <q-item>
            <q-input square filled dense disable class="full-width " hide-bottom-space v-model="props.namespace">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.iot.fleet.view.namespace') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="props.uuid">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.iot.fleet.view.uuid') }}:</div>
                </template>
            </q-input>
        </q-item>

        <q-separator class="q-mt-md q-mb-md"></q-separator>

        <q-item>
            <q-input square filled dense class="full-width" hide-bottom-space v-model="name" :disable="!updatesEnabled">
                <template v-slot:before>
                    <div class="text-h6 text-black">{{ $t('modules.iot.fleet.view.name') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense class="full-width" type="textarea" hide-bottom-space autogrow v-model="description" :disable="!updatesEnabled">
                <template v-slot:before>
                    <div class="text-h6 text-black">{{ $t('modules.iot.fleet.view.description') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-btn :label="$t('modules.iot.fleet.view.updateButton')" color="dark" size="md" class="q-ma-sm full-width" v-if="updatesEnabled" @click="update"></q-btn>

        <q-separator class="q-mt-md q-mb-md"></q-separator>

        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="created">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.iot.fleet.view.created') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="updated">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.iot.fleet.view.updated') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="version">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.iot.fleet.view.version') }}:</div>
                </template>
            </q-input>
        </q-item>
    </q-list>
    
    
    </div>

    <div class="col-sm-12 col-md-8 q-pl-md q-pr-md q-pb-md">
        <FleetDevicesListComponent
            :namespace="props.namespace"
            :uuid="props.uuid"
            :editable="updatesEnabled"
        />
    </div>

    <q-inner-loading :showing="loading">
        <q-spinner-gears size="50px" color="dark" class="q-mb-md"/>
        {{ $t('modules.iot.fleet.view.loading') }}
    </q-inner-loading>
</div>
</template>

<script setup lang="ts">
    import { useQuasar } from 'quasar';
import api from 'src/boot/api';
import { onMounted,computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
    import FleetDevicesListComponent from './FleetDevicesListComponent.vue';
import { Fleet } from 'src/boot/api/iot/fleet';

const props = defineProps<{
    namespace: string,
    uuid: string,
    updatePossible?: boolean
}>()

const fleetSelected = computed(()=>{
    return props.uuid != ''
})

const $q = useQuasar()
const $i18n = useI18n()

const updatesEnabled = ref(false)

const loaded = ref(false)
const loading = ref(false)
const loadingError = ref('')

const name = ref('')
const description = ref('')
const created = ref('')
const updated = ref('')
const version = ref('')

/*

const policies = ref([] as Array<{namespace: string, uuid: string}>)

async function onPolicyAdded(namespace: string, uuid: string) {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.iot.fleet.view.addPolicyOperationNotify')
    })
    try {
        const role = await api.accessControl.role.addPolicy({
            roleNamespace: props.namespace,
            roleUUID: props.uuid,
            policyNamespace: namespace,
            policyUUID: uuid
        })
        applyRoleData(role)
        notif({
            type: 'positive',
            message: $i18n.t('modules.iot.fleet.view.addPolicySuccessNotify')
        })
    } catch (error) {
        notif({
            type: 'negative',
            timeout: 5000,
            message: $i18n.t('modules.iot.fleet.view.addPolicyFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}
async function onDeviceRemoved(namespace: string, uuid: string) {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.iot.fleet.view.removePolicyOperationNotify')
    })
    try {
        const role = await api.accessControl.role.removePolicy({
            roleNamespace: props.namespace,
            roleUUID: props.uuid,
            policyNamespace: namespace,
            policyUUID: uuid
        })
        applyFleetData(role)
        notif({
            type: 'positive',
            message: $i18n.t('modules.iot.fleet.view.removePolicySuccessNotify')
        })
    } catch (error) {
        notif({
            type: 'negative',
            timeout: 5000,
            message: $i18n.t('modules.iot.fleet.view.removePolicyFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}
*/

// const namespaceIndependent = ref(false)

async function loadData() {
    loading.value = true

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.iot.fleet.view.loadOperationNotify')
    })

    try {
        const response = await api.iot.fleet.getFleet({
            namespace: props.namespace,
            uuid: props.uuid
        })
        applyFleetData(response.fleet)
        notif()
        loadingError.value = ""
    } catch (error) {
        loadingError.value = String(error)
        notif({
            type: 'negative',
            message: $i18n.t('modules.iot.fleet.view.loadFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}

async function update() {
    loading.value = true

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.iot.fleet.view.updateOperationNotify')
    })

    try {
        const updateResponse = await api.iot.fleet.updateFleet({
            namespace: props.namespace,
            newDescription: description.value,
            uuid: props.uuid
        })
        applyFleetData(updateResponse.fleet)

        notif({
            type: 'positive',
            message: $i18n.t('modules.iot.fleet.view.updateSuccessfullNotify')
        })
    } catch(error) {
        notif({
            type: 'negative',
            timeout: 5000,
            message: $i18n.t('modules.iot.fleet.view.updateFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}

function applyFleetData(fleet: Fleet) {
    name.value = fleet.name
    description.value = fleet.description
    created.value = fleet.created.toISOString()
    updated.value = fleet.updated.toISOString()
    version.value = fleet.version.toString()
}

watch(props, async (newProps)=>{
    if (newProps.uuid !== '') {
        await loadData()
    }
})

onMounted(async () => {
    if (props.uuid !== '') {
        await loadData()
    }
})

</script>