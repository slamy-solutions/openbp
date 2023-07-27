<template>
<h6 class="row q-ma-md" v-if="policySelected">
    <h6 class="col-6 q-mt-none q-mb-none">{{ $t('modules.accessControl.iam.role.view.header') }}</h6>
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
<div class="row q-ma-xl full-width" v-if="!policySelected">
    <span class="col text-center text-h5 text-bold">
        {{ $t('modules.accessControl.iam.role.view.notSelected') }}
    </span>
</div>
<div class="row q-ma-xl full-width" v-if="policySelected && loadingError !== ''">
    <span class="col text-center text-h4 text-bold text-negative">
        {{ $t('modules.accessControl.iam.role.view.error', { error: loadingError }) }}
    </span>
</div>
<div
    class="row full-width"
    v-if="policySelected && loadingError === ''"
>

    <div
        class="col-sm-12 col-md-6 q-pl-md q-pr-md q-pb-md"
    >
    <q-list dense>
        <q-item>
            <q-input square filled dense disable class="full-width " hide-bottom-space v-model="props.namespace">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.role.view.namespace') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="props.uuid">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.role.view.uuid') }}:</div>
                </template>
            </q-input>
        </q-item>

        <q-separator class="q-mt-md q-mb-md"></q-separator>

        <q-item>
            <q-input square filled dense class="full-width" hide-bottom-space v-model="name" :disable="!updatesEnabled">
                <template v-slot:before>
                    <div class="text-h6 text-black">{{ $t('modules.accessControl.iam.role.view.name') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense class="full-width" type="textarea" hide-bottom-space autogrow v-model="description" :disable="!updatesEnabled">
                <template v-slot:before>
                    <div class="text-h6 text-black">{{ $t('modules.accessControl.iam.role.view.description') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-btn :label="$t('modules.accessControl.iam.role.view.updateButton')" color="dark" size="md" class="q-ma-sm full-width" v-if="updatesEnabled" @click="update"></q-btn>
        <q-item disabled>
            <div class="row full-width items-center">
                <div class="col-3 text-h6">{{ $t('modules.accessControl.iam.role.view.managedBy') }}:</div>
                <ManagedByComponent class="col-9" :managed-by="managed" />
            </div>
        </q-item>

        <q-separator class="q-mt-md q-mb-md"></q-separator>

        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="created">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.role.view.created') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="updated">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.role.view.updated') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="version">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.role.view.version') }}:</div>
                </template>
            </q-input>
        </q-item>
    </q-list>
    
    
    </div>

    <div class="col-sm-12 col-md-6 q-pl-md q-pr-md q-pb-md">
        <h6 class="q-ma-sm text-uppercase q-pa-none">
            {{ $t('modules.accessControl.iam.role.view.policies') }}
            <q-icon name="help" size="20px" class="q-mb-xs" >
                <q-tooltip>{{ $t('modules.accessControl.iam.role.view.policiesCaption') }}</q-tooltip>
            </q-icon>
        </h6>
        <PolicyListComponent
            :namespace="props.namespace"
            :policies="policies"
            :editable="updatesEnabled"
            @added="onPolicyAdded"
            @removed="onPolicyRemoved"
        />
    </div>

    <q-inner-loading :showing="loading">
        <q-spinner-gears size="50px" color="dark" class="q-mb-md"/>
        {{ $t('modules.accessControl.iam.role.view.loading') }}
    </q-inner-loading>
</div>
</template>

<script setup lang="ts">
  import { useQuasar } from 'quasar';
import api from 'src/boot/api';
import { Policy } from 'src/boot/api/accessControl/policy';
import { Role } from 'src/boot/api/accessControl/role';
import { onMounted,computed, Ref, ref, watch, PropType, toRefs } from 'vue';
import { useI18n } from 'vue-i18n';
  import StringListInputComponent from './StringListInputComponent.vue'
  import PolicyListComponent from '../policy/PolicyListComponent.vue';
import ManagedByComponent from '../../../../components/managedItem/ManagedByComponent.vue'
import { ManagedBy } from '../../../../components/managedItem/model'

const props = defineProps<{
    namespace: string,
    uuid: string,
    updatePossible?: boolean
}>()

const policySelected = computed(()=>{
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
const managed = ref({ type: 'none' } as ManagedBy)
const created = ref('')
const updated = ref('')
const version = ref('')

const policies = ref([] as Array<{namespace: string, uuid: string}>)

async function onPolicyAdded(namespace: string, uuid: string) {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.role.view.addPolicyOperationNotify')
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
            message: $i18n.t('modules.accessControl.iam.role.view.addPolicySuccessNotify')
        })
    } catch (error) {
        notif({
            type: 'negative',
            timeout: 5000,
            message: $i18n.t('modules.accessControl.iam.role.view.addPolicyFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}
async function onPolicyRemoved(namespace: string, uuid: string) {
    loading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.role.view.removePolicyOperationNotify')
    })
    try {
        const role = await api.accessControl.role.removePolicy({
            roleNamespace: props.namespace,
            roleUUID: props.uuid,
            policyNamespace: namespace,
            policyUUID: uuid
        })
        applyRoleData(role)
        notif({
            type: 'positive',
            message: $i18n.t('modules.accessControl.iam.role.view.removePolicySuccessNotify')
        })
    } catch (error) {
        notif({
            type: 'negative',
            timeout: 5000,
            message: $i18n.t('modules.accessControl.iam.role.view.removePolicyFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}

// const namespaceIndependent = ref(false)

async function loadData() {
    loading.value = true

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.role.view.loadOperationNotify')
    })

    try {
        const role = await api.accessControl.role.get({
            namespace: props.namespace,
            uuid: props.uuid
        })
        applyRoleData(role)
        notif()
        loadingError.value = ""
    } catch (error) {
        loadingError.value = String(error)
        notif({
            type: 'negative',
            message: $i18n.t('modules.accessControl.iam.role.view.loadFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}

async function update() {
    loading.value = true

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.role.view.updateOperationNotify')
    })

    try {
        const updateResponse = await api.accessControl.role.update({
            namespace: props.namespace,
            newDescription: description.value,
            newName: name.value,
            uuid: props.uuid
        })
        applyRoleData(updateResponse.role)

        notif({
            type: 'positive',
            message: $i18n.t('modules.accessControl.iam.role.view.updateSuccessfullNotify')
        })
    } catch(error) {
        notif({
            type: 'negative',
            timeout: 5000,
            message: $i18n.t('modules.accessControl.iam.role.view.updateFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}

function applyRoleData(role: Role) {
    name.value = role.name
    description.value = role.description
    policies.value = role.policies
    managed.value = role.managed
    created.value = role.created.toISOString()
    updated.value = role.updated.toISOString()
    version.value = role.version.toString()
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