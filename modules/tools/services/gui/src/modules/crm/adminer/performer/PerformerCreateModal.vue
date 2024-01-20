<template>
    <q-card
        flat
        class="q-pa-xs"
        square
        style="width: 90%; max-width: 350px;"
    >
        <q-card-section>
            <div class="text-h6 text-center">{{ $t('modules.crm.adminer.performer.create.title') }}</div>
        </q-card-section>

        <q-card-section>
            <UserSelectorInput v-model="user" :namespace="props.namespace" class="q-mb-sm"/>
            <DepartmentSelectorInput v-model="department" :namespace="props.namespace" />
        </q-card-section>
    
        <q-card-actions>
            <q-btn
                :label="$t('modules.crm.adminer.performer.create.createButton')"
                class="q-ma-none fit"
                unelevated
                outline
                size="md"
                :disable="loading || user === null || department === null"
                :loading="loading"
                @click="create"
            />
        </q-card-actions>
    </q-card>
</template>

<script setup lang="ts">
import { User } from 'src/boot/api/accessControl/actor/user';
import { Department } from 'src/boot/api/crm/department';
import { Performer } from 'src/boot/api/crm/performer';
import UserSelectorInput from 'src/modules/accessControl/actor/user/UserSelectorInput.vue';
import DepartmentSelectorInput from '../department/DepartmentSelectorInput.vue';
import { defineEmits, ref } from 'vue';
import api from 'src/boot/api';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    namespace: string,
}>()

const emits = defineEmits<{
    (e: 'created', performer: Performer): void
    (e: 'cancelled'): void
}>()

const loading = ref(false)

const user = ref(null as null | User)
const department = ref(null as null | Department)

async function create() {
    loading.value = true
    const notif = $q.notify({
        message: $i18n.t('modules.crm.adminer.performer.create.createOperationNotify'),
        type: 'ongoing',
    })
    try {
        const response = await api.crm.performer.create({
            namespace: props.namespace,
            userUUID: user.value?.uuid ?? '',
            departmentUUID: department.value?.uuid ?? ''
        })
        emits('created', response.performer)
        notif({
            message: $i18n.t('modules.crm.adminer.performer.create.createSuccessNotify'),
            type: 'positive',
            timeout: 5000
        })
    } catch (error) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.crm.adminer.performer.create.createFailNotify', { error }),
            timeout: 5000
        })
    } finally {
        loading.value = false
    }
}

</script>