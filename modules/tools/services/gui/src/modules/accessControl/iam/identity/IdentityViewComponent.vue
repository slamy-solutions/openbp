<template>
    <h6 class="row q-ma-md" v-if="policySelected">
        <h6 class="col-6 q-mt-none q-mb-none">{{ $t('modules.accessControl.iam.identity.view.header') }}</h6>
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
            {{ $t('modules.accessControl.iam.identity.view.notSelected') }}
        </span>
    </div>
    <div class="row q-ma-xl full-width" v-if="policySelected && loadingError !== ''">
        <span class="col text-center text-h4 text-bold text-negative">
            {{ $t('modules.accessControl.iam.identity.view.error', { error: loadingError }) }}
        </span>
    </div>
    <div
        class="row full-width"
        v-if="policySelected && loadingError === ''"
    >
    
        <div
            class="col-sm-12 col-md-4 q-pl-md q-pr-md q-pb-md"
        >
        <q-list dense>
            <q-item>
                <q-input square filled dense disable class="full-width " hide-bottom-space v-model="props.namespace">
                    <template v-slot:before>
                        <div class="text-h6">{{ $t('modules.accessControl.iam.identity.view.namespace') }}:</div>
                    </template>
                </q-input>
            </q-item>
            <q-item>
                <q-input square filled dense disable class="full-width" hide-bottom-space v-model="props.uuid">
                    <template v-slot:before>
                        <div class="text-h6">{{ $t('modules.accessControl.iam.identity.view.uuid') }}:</div>
                    </template>
                </q-input>
            </q-item>
    
            <q-separator class="q-mt-md q-mb-md"></q-separator>
    
            <q-item :disable="!updatesEnabled">
                <q-input square filled dense class="full-width" hide-bottom-space v-model="name" :disable="!updatesEnabled">
                    <template v-slot:before>
                        <div class="text-h6 text-black">{{ $t('modules.accessControl.iam.identity.view.name') }}:</div>
                    </template>
                </q-input>
            </q-item>

            <q-btn label="UPDATE"  color="dark" size="md" class="q-ma-md full-width" v-if="updatesEnabled" @click="update"></q-btn>

            <q-item tag="label" :v-ripple="updatesEnabled" :disable="!updatesEnabled" class="q-mr-md">
                <q-item-section>
                <q-item-label class="text-h6">{{ $t('modules.accessControl.iam.identity.view.disabled') }}:</q-item-label>
                </q-item-section>
                <q-item-section side top>
                <q-toggle color="red" v-model="disabled" val="picture" :disable="!updatesEnabled" @update:model-value="setDisabled"/>
                </q-item-section>
            </q-item>

            <q-item disabled>
                <div class="row full-width items-center">
                    <div class="col-3 text-h6">{{ $t('modules.accessControl.iam.identity.view.managedBy') }}:</div>
                    <ManagedByComponent class="col-9" :managed-by="managed" />
                </div>
            </q-item>
    
            <q-separator class="q-mt-md q-mb-md"></q-separator>
    
            <q-item>
                <q-input square filled dense disable class="full-width" hide-bottom-space v-model="created">
                    <template v-slot:before>
                        <div class="text-h6">{{ $t('modules.accessControl.iam.identity.view.created') }}:</div>
                    </template>
                </q-input>
            </q-item>
            <q-item>
                <q-input square filled dense disable class="full-width" hide-bottom-space v-model="updated">
                    <template v-slot:before>
                        <div class="text-h6">{{ $t('modules.accessControl.iam.identity.view.updated') }}:</div>
                    </template>
                </q-input>
            </q-item>
            <q-item>
                <q-input square filled dense disable class="full-width" hide-bottom-space v-model="version">
                    <template v-slot:before>
                        <div class="text-h6">{{ $t('modules.accessControl.iam.identity.view.version') }}:</div>
                    </template>
                </q-input>
            </q-item>

            <q-separator class="q-mt-md q-mb-md"></q-separator>

            
        </q-list>
        
        
        </div>
        
        <div class="col-sm-12 col-md-8 q-pr-md q-pl-md q-pb-md">

        
        <q-tabs v-model="additionalInfoTab" class="bg-primary" indicator-color="secondary" >
            <q-tab :label="$t('modules.accessControl.iam.identity.view.tabs.privileges')" name="privileges"/>
            <q-tab :label="$t('modules.accessControl.iam.identity.view.tabs.password')" name="password"/>
            <q-tab :label="$t('modules.accessControl.iam.identity.view.tabs.oauth')" name="oauth"/>
            <q-tab :label="$t('modules.accessControl.iam.identity.view.tabs.tokens')" name="tokens"/>
            <q-tab :label="$t('modules.accessControl.iam.identity.view.tabs.2fa')" name="2fa"/>
            <q-tab :label="$t('modules.accessControl.iam.identity.view.tabs.certificates')" name="certificates"/>
        </q-tabs>

        <q-tab-panels v-model="additionalInfoTab" animated class="transparent">
            <q-tab-panel name="privileges" class="row bg-transparent">
                <div class="col-sm-12 col-md-6 q-pl-md q-pr-md q-pb-md">
                    <h6 class="q-ma-sm text-uppercase q-pa-none">
                        {{ $t('modules.accessControl.iam.identity.view.roles') }}
                        <q-icon name="help" size="20px" class="q-mb-xs" >
                            <q-tooltip>{{ $t('modules.accessControl.iam.identity.view.rolesCaption') }}</q-tooltip>
                        </q-icon>
                    </h6>
                    <RoleListComponent
                        :namespace="props.namespace"
                        :roles="roles"
                        :editable="updatesEnabled"
                        @added="onRoleAdded"
                        @removed="onRoleRemoved"
                    />
                </div>

                <div class="col-sm-12 col-md-6 q-pl-md q-pr-md q-pb-md">
                    <h6 class="q-ma-sm text-uppercase q-pa-none">
                        {{ $t('modules.accessControl.iam.identity.view.policies') }}
                        <q-icon name="help" size="20px" class="q-mb-xs" >
                            <q-tooltip>{{ $t('modules.accessControl.iam.identity.view.policiesCaption') }}</q-tooltip>
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
            </q-tab-panel>

            <q-tab-panel name="password" class="q-pa-xs q-pt-md">
                <PasswordConfigComponent
                    class="full-width transparent"
                    :identityNamespace="props.namespace"
                    :identityUUID="props.uuid"
                    :editable="updatesEnabled"
                />
            </q-tab-panel>

            <q-tab-panel name="oauth" class="q-pa-xs q-pt-md">
                Not Implemented
            </q-tab-panel>

            <q-tab-panel name="tokens" class="q-pa-xs q-pt-md">
                Not Implemented
            </q-tab-panel>

            <q-tab-panel name="2fa" class="q-pa-xs q-pt-md">
                Not Implemented
            </q-tab-panel>

            <q-tab-panel name="certificates" class="q-pa-xs q-pt-md">
                <IdentityCertificatesListComponent
                    class="full-width"
                    :identityNamespace="props.namespace"
                    :identityUUID="props.uuid"
                    :editable="updatesEnabled"
                />
            </q-tab-panel>
        </q-tab-panels>
    </div>
    
        <q-inner-loading :showing="loading">
            <q-spinner-gears size="50px" color="dark" class="q-mb-md"/>
            {{ $t('modules.accessControl.iam.identity.view.loading') }}
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
      import RoleListComponent from '../role/RoleListComponent.vue';
      import PasswordConfigComponent from '../auth/PasswordConfigComponent.vue'
    import ManagedByComponent from '../../../../components/managedItem/ManagedByComponent.vue'
    import IdentityCertificatesListComponent from './IdentityCertificatesListComponent.vue';
    import { ManagedBy } from '../../../../components/managedItem/model'
import { Identity } from 'src/boot/api/accessControl/identity';
    
    const props = defineProps<{
        namespace: string,
        uuid: string,
        updatePossible?: boolean
    }>()
    
    const policySelected = computed(()=>{
        return props.uuid != ''
    })
    
    const additionalInfoTab = ref("privileges")

    const $q = useQuasar()
    const $i18n = useI18n()
    
    const updatesEnabled = ref(false)
    
    const loading = ref(false)
    const loadingError = ref('')
    
    const name = ref('')
    const disabled = ref(false)
    const managed = ref({ type: 'none' } as ManagedBy)
    const created = ref('')
    const updated = ref('')
    const version = ref('')
    
    const policies = ref([] as Array<{
        namespace: string
        uuid: string
    }>)

    const roles = ref([] as Array<{
        namespace: string
        uuid: string
    }>)
    
    async function onPolicyAdded(namespace: string, uuid: string) {
        loading.value = true
    
        const notif = $q.notify({
            type: 'ongoing',
            message: $i18n.t('modules.accessControl.iam.identity.view.addPolicyOperationNotify')
        })

        try {
            const identity = await api.accessControl.identity.addPolicy({
                identityNamespace: props.namespace,
                identityUUID: props.uuid,
                policyNamespace: namespace,
                policyUUID: uuid
            })
            applyIdentityData(identity)

            notif({
                type: 'positive',
                message: $i18n.t('modules.accessControl.iam.identity.view.addPolicySuccessNotify')
            })
        } catch(error) {
            notif({
                type: 'negative',
                timeout: 5000,
                message: $i18n.t('modules.accessControl.iam.identity.view.addPolicyFailNotify', { error })
            })
        } finally {
            loading.value = false
        }
    }
    async function onPolicyRemoved(namespace: string, uuid: string) {
        loading.value = true
    
        const notif = $q.notify({
            type: 'ongoing',
            message: $i18n.t('modules.accessControl.iam.identity.view.removePolicyOperationNotify')
        })

        try {
            const identity = await api.accessControl.identity.removePolicy({
                identityNamespace: props.namespace,
                identityUUID: props.uuid,
                policyNamespace: namespace,
                policyUUID: uuid
            })
            applyIdentityData(identity)

            notif({
                type: 'positive',
                message: $i18n.t('modules.accessControl.iam.identity.view.removePolicySuccessNotify')
            })
        } catch(error) {
            notif({
                type: 'negative',
                timeout: 5000,
                message: $i18n.t('modules.accessControl.iam.identity.view.removePolicyFailNotify', { error })
            })
        } finally {
            loading.value = false
        }
    }

    async function onRoleAdded(namespace: string, uuid: string) {
        loading.value = true
    
        const notif = $q.notify({
            type: 'ongoing',
            message: $i18n.t('modules.accessControl.iam.identity.view.addRoleOperationNotify')
        })

        try {
            const identity = await api.accessControl.identity.addRole({
                identityNamespace: props.namespace,
                identityUUID: props.uuid,
                roleNamespace: namespace,
                roleUUID: uuid
            })
            applyIdentityData(identity)

            notif({
                type: 'positive',
                message: $i18n.t('modules.accessControl.iam.identity.view.addRoleSuccessNotify')
            })
        } catch(error) {
            notif({
                type: 'negative',
                timeout: 5000,
                message: $i18n.t('modules.accessControl.iam.identity.view.addRoleFailNotify', { error })
            })
        } finally {
            loading.value = false
        }
    }
    async function onRoleRemoved(namespace: string, uuid: string) {
        loading.value = true
    
        const notif = $q.notify({
            type: 'ongoing',
            message: $i18n.t('modules.accessControl.iam.identity.view.removeRoleOperationNotify')
        })

        try {
            const identity = await api.accessControl.identity.removeRole({
                identityNamespace: props.namespace,
                identityUUID: props.uuid,
                roleNamespace: namespace,
                roleUUID: uuid
            })
            applyIdentityData(identity)

            notif({
                type: 'positive',
                message: $i18n.t('modules.accessControl.iam.identity.view.removeRoleSuccessNotify')
            })
        } catch(error) {
            notif({
                type: 'negative',
                timeout: 5000,
                message: $i18n.t('modules.accessControl.iam.identity.view.removeRoleFailNotify', { error })
            })
        } finally {
            loading.value = false
        }
    }
    
    
    async function loadData() {
        loading.value = true
    
        const notif = $q.notify({
            type: 'ongoing',
            message: $i18n.t('modules.accessControl.iam.identity.view.loadOperationNotify')
        })
    
        try {
            const identity = await api.accessControl.identity.get({
                namespace: props.namespace,
                uuid: props.uuid
            })
            applyIdentityData(identity)
            notif()
            loadingError.value = ""
        } catch (error) {
            loadingError.value = String(error)
            notif({
                type: 'negative',
                message: $i18n.t('modules.accessControl.iam.identity.view.loadFailNotify', { error })
            })
        } finally {
            loading.value = false
        }
    }
    
    async function update() {
        loading.value = true
    
        const notif = $q.notify({
            type: 'ongoing',
            message: $i18n.t('modules.accessControl.iam.identity.view.updateOperationNotify')
        })
    
        try {
            const updateResponse = await api.accessControl.identity.update({
                namespace: props.namespace,
                newName: name.value,
                uuid: props.uuid
            })
            applyIdentityData(updateResponse)
    
            notif({
                type: 'positive',
                message: $i18n.t('modules.accessControl.iam.identity.view.updateSuccessfullNotify')
            })
        } catch(error) {
            notif({
                type: 'negative',
                timeout: 5000,
                message: $i18n.t('modules.accessControl.iam.identity.view.updateFailNotify', { error })
            })
        } finally {
            loading.value = false
        }
    }

    async function setDisabled(disabled: boolean) {
        loading.value = true
        const notif = $q.notify({
            type: 'ongoing',
            message: $i18n.t('modules.accessControl.iam.identity.view.activeOperationNotify')
        })
        try {
            const response = await api.accessControl.identity.setActive({
                namespace: props.namespace,
                uuid: props.uuid,
                active: !disabled
            })
            applyIdentityData(response)

            notif({
                type: 'positive',
                message: $i18n.t('modules.accessControl.iam.identity.view.activeSuccessfullNotify')
            })
        } catch (error) {
            notif({
                type: 'negative',
                timeout: 5000,
                message: $i18n.t('modules.accessControl.iam.identity.view.activeFailNotify', { error })
            })
        } finally {
            loading.value = false
        }
    }
    
    function applyIdentityData(identity: Identity) {
        name.value = identity.name
        policies.value = identity.policies
        roles.value = identity.roles
        managed.value = identity.managed
        disabled.value = !identity.active
        created.value = identity.created.toISOString()
        updated.value = identity.updated.toISOString()
        version.value = identity.version.toString()
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