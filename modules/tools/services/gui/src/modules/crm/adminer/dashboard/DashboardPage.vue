<template>
    <q-page padding class="row">
        <div class="col-2">
            <MenuComponent selected="dashboard"/>
        </div>
        <div class="col-10">
            <div class="row q-gutter-sm">
                <q-card flat bordered class="bg-transparent col-12">
                    <q-card-section>
                        <div class="text-h6">{{ $t("modules.crm.adminer.dashboard.backend") }}</div>
                        <div class="text-caption" v-if="!loading">{{ backendType }}</div>
                        <q-spinner v-if="loading" size="xs" />
                    </q-card-section>
                </q-card>

                <q-card flat bordered class="bg-transparent col" v-if="backendType == 'ONE_C'">
                    <SyncStatusTableComponent :namespace="displayableNamespace" />
                </q-card>
            </div>
        </div>
    </q-page>
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar';
import api from '../../../../boot/api';
import { BackendType } from '../../../../boot/api/crm/settings';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import SyncStatusTableComponent from '../onec/SyncStatusTableComponent.vue'
import MenuComponent from '../MenuComponent.vue'

const $i18n = useI18n()
const $q = useQuasar()

const $route = useRoute()
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string

const loading = ref(false)
const backendType = ref<BackendType>('NATIVE')

async function loadSettings() {
    loading.value = true
    try {
        const response = await api.crm.settings.get({namespace: displayableNamespace})
        backendType.value = response.settings.backendType
    } catch (error) {
        $q.notify({
            type: 'negative',
            message: $i18n.t('modules.crm.adminer.dashboard.loadSettingsFailNotify', { error }),
            timeout: 5000
        })
        console.error(error)
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    await loadSettings()
})
</script>
