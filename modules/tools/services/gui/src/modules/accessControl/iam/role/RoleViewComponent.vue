<template>
<h6 class="row q-ma-md" v-if="policySelected">
    <h6 class="col-6 q-mt-none q-mb-none">{{ $t('modules.accessControl.iam.role.view.header') }}</h6>
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
        class="col-4 q-pl-md q-pr-md q-pb-md"
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
        <q-item>
            <ManagedByComponent :managed-by="managed" />
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

    <div class="col-4 q-pl-md q-pr-md q-pb-md">
        
    </div>

    <div class="col-4 q-pl-md q-pr-md q-pb-md">
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
        message: $i18n.t('modules.accessControl.iam.role.view.loadOperationNotify')
    })

    try {
        const role = await api.accessControl.role.get({
            namespace: props.namespace,
            uuid: props.uuid
        })
        applyPolicyData(role)
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
    /*loading.value = true
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
    }*/
}

function applyPolicyData(role: Role) {
    name.value = role.name
    description.value = role.description
    managed.value = role.managed
    created.value = role.created.toISOString()
    updated.value = role.updated.toISOString()
    version.value = role.version.toString()

    /*resources.value = role.resources.map((r) => {
        return {
            id: Math.random().toString(),
            value: r
        }
    })

    actions.value = role.actions.map((r) => {
        return {
            id: Math.random().toString(),
            value: r
        }
    })*/
}

watch(props, async (newProps)=>{
    if (newProps.uuid !== '') {
        console.log("Change")
        await loadData()
    }
})

</script>