<template>
<h6 class="row q-ma-md" v-if="policySelected">
    <h6 class="col-6 q-mt-none q-mb-none">{{ $t('modules.accessControl.iam.policy.view.header') }}</h6>
    <div class="col-6 text-right">
        <q-btn label="UPDATE"  outline color="positive" size="md" class="q-mr-md" v-if="updatesEnabled" @click="update"></q-btn>
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
        {{ $t('modules.accessControl.iam.policy.view.notSelected') }}
    </span>
</div>
<div class="row q-ma-xl full-width" v-if="policySelected && loadingError !== ''">
    <span class="col text-center text-h4 text-bold text-negative">
        {{ $t('modules.accessControl.iam.policy.view.error', { error: loadingError }) }}
    </span>
</div>
<div
    class="row full-width"
    v-if="policySelected && loadingError === ''"
>

    <div
        class="col-4 q-pl-md q-pr-md q-pb-md"
    >
    <q-list dense>
        <q-item>
            <q-input square filled dense disable class="full-width " hide-bottom-space v-model="props.namespace">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.policy.view.namespace') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="props.uuid">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.policy.view.uuid') }}:</div>
                </template>
            </q-input>
        </q-item>

        <q-separator class="q-mt-md q-mb-md"></q-separator>

        <q-item>
            <q-input square filled dense class="full-width" hide-bottom-space v-model="name" :disable="!updatesEnabled">
                <template v-slot:before>
                    <div class="text-h6 text-black">{{ $t('modules.accessControl.iam.policy.view.name') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense class="full-width" type="textarea" hide-bottom-space autogrow v-model="description" :disable="!updatesEnabled">
                <template v-slot:before>
                    <div class="text-h6 text-black">{{ $t('modules.accessControl.iam.policy.view.description') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item tag="label" v-ripple="updatesEnabled" :disable="!updatesEnabled" class="q-mr-md">
            <q-item-section>
            <q-item-label class="text-h6 text-black text-left">{{ $t('modules.accessControl.iam.policy.view.namespaceIndependent') }}</q-item-label>
            </q-item-section>
            <q-item-section side >
            <q-toggle color="blue" v-model="namespaceIndependent" val="battery" :disable="!updatesEnabled"/>
            </q-item-section>
        </q-item>
        <q-item class="text-center">
            <div class="text-h6 text-black">{{ $t('modules.accessControl.iam.policy.view.managedBy') }}: <ManagedByComponent :managed-by="managed" /></div>
        </q-item>

        <q-separator class="q-mt-md q-mb-md"></q-separator>

        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="created">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.policy.view.created') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="updated">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.policy.view.updated') }}:</div>
                </template>
            </q-input>
        </q-item>
        <q-item>
            <q-input square filled dense disable class="full-width" hide-bottom-space v-model="version">
                <template v-slot:before>
                    <div class="text-h6">{{ $t('modules.accessControl.iam.policy.view.version') }}:</div>
                </template>
            </q-input>
        </q-item>
    </q-list>
    
    
    </div>

    <div
        class="col-4 q-pl-md q-pr-md q-pb-md"
    >
    <h6 class="q-ma-sm text-uppercase q-pa-none">
                {{ $t('modules.accessControl.iam.policy.view.resources') }}
                <q-icon name="help" size="20px" class="q-mb-xs" >
                  <q-tooltip>{{ $t('modules.accessControl.iam.policy.view.resourcesCaption') }}</q-tooltip>
                </q-icon>
              </h6>
        <StringListInputComponent :value="resources" :disable="!updatesEnabled"></StringListInputComponent>
        
    </div>

    <div
        class="col-4 q-pl-md q-pr-md q-pb-md"
    >
    <h6 class="q-ma-sm text-uppercase q-pa-none">
                {{ $t('modules.accessControl.iam.policy.view.actions') }}
                <q-icon name="help" size="20px" class="q-mb-xs" >
                  <q-tooltip>{{ $t('modules.accessControl.iam.policy.view.actionsCaption') }}</q-tooltip>
                </q-icon>
              </h6>
        <StringListInputComponent :value="actions" :disable="!updatesEnabled"></StringListInputComponent>
        
    </div>

    <q-inner-loading :showing="loading">
        <q-spinner-gears size="50px" color="dark" class="q-mb-md"/>
        {{ $t('modules.accessControl.iam.policy.view.loading') }}
    </q-inner-loading>
</div>
</template>

<script setup lang="ts">
  import { useQuasar } from 'quasar';
import api from 'src/boot/api';
import { Policy } from 'src/boot/api/accessControl/policy';
import { onMounted,computed, Ref, ref, watch, PropType, toRefs } from 'vue';
import { useI18n } from 'vue-i18n';
  import StringListInputComponent from './StringListInputComponent.vue'
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

const resources = ref([] as Array<{
    id: string
    value: string   
}>)
const actions = ref([] as Array<{
    id: string
    value: string   
}>)
const namespaceIndependent = ref(false)

async function loadData() {
    loading.value = true

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.policy.view.loadOperationNotify')
    })

    try {
        const policy = await api.accessControl.policy.get({
            namespace: props.namespace,
            uuid: props.uuid
        })
        applyPolicyData(policy)
        notif()
        loadingError.value = ""
    } catch (error) {
        loadingError.value = String(error)
        notif({
            type: 'negative',
            message: $i18n.t('modules.accessControl.iam.policy.view.loadFailNotify', { error })
        })
    } finally {
        loading.value = false
    }
}

async function update() {
    loading.value = true
    try {
        const policy = await api.accessControl.policy.update({
            namespace: props.namespace,
            uuid: props.uuid,
            name: name.value,
            description: description.value,
            actions: actions.value.map((v) => v.value),
            resources: resources.value.map((v) => v.value),
            namespaceIndependent: namespaceIndependent.value
        })
        applyPolicyData(policy)
    } finally {
        loading.value = false
    }
}

function applyPolicyData(policy: Policy) {
    name.value = policy.name
    description.value = policy.description
    managed.value = policy.managed
    created.value = policy.created.toISOString()
    updated.value = policy.updated.toISOString()
    version.value = policy.version.toString()

    resources.value = policy.resources.map((r) => {
        return {
            id: Math.random().toString(),
            value: r
        }
    })

    actions.value = policy.actions.map((r) => {
        return {
            id: Math.random().toString(),
            value: r
        }
    })
    namespaceIndependent.value = policy.namespaceIndependent
}

watch(props, async (newProps)=>{
    if (newProps.uuid !== '') {
        console.log("Change")
        await loadData()
    }
})

</script>