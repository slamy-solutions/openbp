<template>

</template>

<script setup lang="ts">
import { QPaginationProps, QTableProps, useQuasar } from 'quasar';
import { onMounted, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import api from '../../../boot/api';

import MenuComponent from '../MenuComponent.vue'
import FleetCreateModal from './FleetCreateModal.vue';
import FleetDeleteModal from './FleetDeleteModal.vue';
import IdentityViewComponent from '../../iam/identity/IdentityViewComponent.vue';
import { Fleet } from 'src/boot/api/iot/fleet';

const $i18n = useI18n()
const $q = useQuasar()

const displayableNamespace = ref("")

const tableColumns: Ref<QTableProps['columns']> = ref([
    {name: 'uuid', required: true, label: $i18n.t('modules.iot.fleet.list.uuidColumn'), align: 'left', sortable: false, field: 'uuid'},
    {name: 'name', required: true, label: $i18n.t('modules.iot.fleet.list.nameColumn'), align: 'left', sortable: false, field: 'name'},
    {name: 'description', required: true, label: $i18n.t('modules.iot.fleet.list.descriptionColumn'), align: 'left', sortable: false, field: 'description'},
    {name: 'actions', required: false, label: $i18n.t('modules.iot.fleet.list.actionsColumn'), align: 'right', sortable: false, field: 'actions'}
])
const tableData = ref([] as Array<Fleet>)
const dataLoading = ref(false)
const loadingError = ref("")
const tableRef = ref()

const tablePagination = ref({
page: 1,
rowsPerPage: 10,
rowsNumber: 0
})

const creationDialog = ref(false)
const fleetUUIDToDelete = ref('')
const deletionDialog = ref(false)


async function loadData(tableProps: QTableProps) {
    const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
    const page = (tableProps.pagination?.page || 1) - 1

    const notif = $q.notify({
    type: 'ongoing',
    message: $i18n.t('modules.iot.fleet.list.loadOperationNotify')
})
    dataLoading.value = true

    try {
    const response = await api.iot.fleet.listFleets({ namespace: displayableNamespace.value, skip: rowsPerPage*page, limit: rowsPerPage })
    tableData.value = response.fleets
    if (tablePagination.value != undefined) {
        tablePagination.value.page = page + 1
        tablePagination.value.rowsPerPage = rowsPerPage
        tablePagination.value.rowsNumber = response.totalCount
    }
    notif()
    loadingError.value = ""

    } catch (error) {
    notif({
        type: 'negative',
        message: $i18n.t('modules.iot.fleet.list.loadFailNotify', { error }),
        timeout: 5000
    })
    loadingError.value = String(error) 
    } finally {
    dataLoading.value = false
    }
}

async function onFleetCreated() {
creationDialog.value = false
tableRef.value.requestServerInteraction()
}

async function onFleetDeleted() {
deletionDialog.value = false
tableRef.value.requestServerInteraction()
}

onMounted(() => {
tableRef.value.requestServerInteraction()
})


</script>