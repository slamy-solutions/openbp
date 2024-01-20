<template>
    <q-card
        square
        style="min-width: 350px;"
    >
        <q-card-section>
            <div class="text-h6 text-center">{{ $t('modules.crm.adminer.department.create.title') }}</div>
        </q-card-section>

        <q-card-section>
            <q-input v-model="name" filled :label="$t('modules.crm.adminer.department.create.nameInput')" />
        </q-card-section>

        <q-card-actions>
            <q-btn
                :label="$t('modules.crm.adminer.department.create.createButton')"
                class="q-ma-none fit"
                unelevated
                outline
                size="md"
                :disable="loading"
                :loading="loading"
                @click="create"
            />
        </q-card-actions>

    </q-card>
</template>

<script setup lang="ts">
import { defineProps, ref } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from '../../../../boot/api';

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string
}>()

const emits = defineEmits<{
    (e: 'created'): void
}>()

const name = ref('')
const loading = ref(false)

async function create() {
    loading.value = true
    const notif = $q.notify({
        message: $i18n.t('modules.crm.adminer.department.create.createOperationNotify'),
        type: 'ongoing',
    })
    try {
        await api.crm.department.create({
            namespace: props.namespace,
            name: name.value
        })
        notif({
            message: $i18n.t('modules.crm.adminer.department.create.createSuccessNotify'),
            type: 'positive',
            timeout: 5000
        })
        emits('created')
    } catch (error) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.crm.adminer.department.create.createFailNotify', { error }),
            timeout: 5000
        })
        console.error(error)
    } finally {

        loading.value = false
    }
}


</script>