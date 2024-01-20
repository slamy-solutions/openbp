<template>
    <q-page padding class="row">
        <div class="col-2">
            <MenuComponent selected="departments" ref="menu"/>
        </div>
        <div class="col-10">
            <q-table
                flat
                class="col-12 q-ma-none bg-transparent"
                :title="$t('modules.crm.adminer.department.list.header')"
                :columns="tableColumns"
                :rows="departments"
                row-key="uuid"
                :loading="loading"
                ref="tableRef"
                dense
                @row-click="selectDepartment"
            >
            <template v-slot:loading>
                <q-inner-loading showing color="secondary" />
            </template>
  
            <template v-slot:no-data="{}">
                <div class="full-width row flex-center q-gutter-sm">
                    <span v-if="loadingError === ''">
                    {{ $t('modules.crm.adminer.department.list.noData') }}
                    </span>
                    <span v-else class="text-negative">
                    {{ $t('modules.crm.adminer.department.list.failedToLoad', { error: loadingError }) }}
                    </span>
                </div>
            </template>
    
            <template v-slot:top-right>
                <q-btn
                    :label="$t('modules.crm.adminer.department.list.createButton')"
                    class="q-ma-none"
                    unelevated
                    outline
                    size="sm"
                    :disable="creationDialog"
                    @click.stop="creationDialog = true"
                />
            </template>
  
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn
                class="q-ma-none"
                unelevated
                outline
                size="sm"
                :disable="deletionDialog"
                icon="delete"
                @click.stop="deletionDialogDepartmentUUID = props.row.uuid; deletionDialogDepartmentName = props.row.name; deletionDialog = true;"
            />
          </q-td>
        </template>
  
        </q-table>
        </div>
    </q-page>

    <q-dialog
        v-model="creationDialog"
    >
        <DepartmentCreateModal
            :namespace="displayableNamespace"
            @created="async () => { creationDialog = false; await loadDepartments(); }"
        />
    </q-dialog>

    <q-dialog
        v-model="deletionDialog"
    >
        <DepartmentDeleteModal
            :namespace="displayableNamespace"
            :departmentUUID="deletionDialogDepartmentUUID"
            :name="deletionDialogDepartmentName"
            @deleted="async () => { deletionDialog = false; await loadDepartments(); }"
        />
    </q-dialog>
</template>

<script setup lang="ts">
import { QTableProps, useQuasar } from 'quasar';
import api from '../../../../boot/api';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import MenuComponent from '../MenuComponent.vue'
import { Department } from '../../../../boot/api/crm/department';

import DepartmentCreateModal from './DepartmentCreateModal.vue'
import DepartmentDeleteModal from './DepartmentDeleteModal.vue'

const $i18n = useI18n()
const $q = useQuasar()
const $router = useRouter()

const $route = useRoute()
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string

const loading = ref(false)
const departments = ref([] as Array<Department>)
const tableColumns = ref<QTableProps['columns']>([
    {name: 'name', required: true, label: $i18n.t('modules.crm.adminer.department.list.nameColumn'), align: 'left', sortable: false, field: 'name'},
    {name: 'actions', required: true, label: $i18n.t('modules.crm.adminer.department.list.actionsColumn'), align: 'right', sortable: false, field: 'actions'}
])
const loadingError = ref("")

const creationDialog = ref(false)

const deletionDialog = ref(false)
const deletionDialogDepartmentUUID = ref("")
const deletionDialogDepartmentName = ref("")

async function selectDepartment(_e: Event, row: Department) {
    await $router.push({ name: 'crm_adminer_department', params: { currentNamespace: $route.params.currentNamespace, departmentUUID: row.uuid } })
}

const menu = ref<{loadDepartments: () => Promise<void>} | null>(null)

async function loadDepartments() {
    loading.value = true
    try {
        const response = await api.crm.department.getAll({ namespace: displayableNamespace })
        departments.value = response.departments
    } catch (error) {
        $q.notify({
            type: 'negative',
            message: $i18n.t('modules.crm.adminer.department.list.loadFailNotify', { error }),
            timeout: 5000
        })
        console.error(error)
    } finally {
        loading.value = false
        menu.value?.loadDepartments()
    }
}


onMounted(async () => {
    await loadDepartments()
})
</script>
