<template>
    <q-table
        class="row q-ma-none bg-transparent"
        flat
        :title="$t('modules.crm.adminer.settings.onec.sync.table.header')"
        :columns="tableColumns"
        :rows="tableData"
        row-key="uuid"
        :loading="dataLoading"
        ref="tableRef"
        @request="loadData"
        @row-click="selectSyncEvent"
        :selected="selectedEvent ? [selectedEvent] : []"
        dense
        v-model:pagination="tablePagination"
    >
        <template v-slot:loading>
            <q-inner-loading showing color="secondary" />
        </template>

    <template v-slot:no-data="{}">
        <div class="full-width row flex-center q-gutter-sm">
        <span v-if="loadingError === ''">
            {{ $t('modules.crm.adminer.settings.onec.sync.table.noData') }}
        </span>
        <span v-else class="text-negative">
            {{ $t('modules.crm.adminer.settings.onec.sync.table.failedToLoad', { error: loadingError }) }}
        </span>
        </div>
    </template>

    <template v-slot:top-right>
        <q-btn
            :label="$t('modules.crm.adminer.settings.onec.sync.table.syncNowButton')"
            class="q-ma-none"
            unelevated
            outline
            size="md"
            :disable="syncing"
            :loading="syncing"
            @click="syncNow"
        />
    </template>

    <template v-slot:body-cell-success="props">
      <q-td :props="props">
        <q-icon name="check_circle" color="positive" v-if="props.row.success" size="sm"/>
        <q-icon name="cancel" color="negative" v-if="!props.row.success" size="sm"/>
      </q-td>
    </template>
    </q-table>

    <q-dialog v-model="eventDetailsDialog">
      <q-card>
        <SyncEventViewComponent v-if="selectedEvent" :event="selectedEvent" />
      </q-card>
    </q-dialog>
</template>
  
<script setup lang="ts">
  import { QPaginationProps, QTableProps, useQuasar } from 'quasar';
  import { onMounted, Ref, ref } from 'vue';
  import { useRoute } from 'vue-router';
  import { useI18n } from 'vue-i18n';
  import api from '../../../../boot/api';
import { SyncEvent } from '../../../../boot/api/crm/onec';
import SyncEventViewComponent from './SyncEventViewComponent.vue';
  
  const $i18n = useI18n()
  const $q = useQuasar()
  
  const $route = useRoute()
  const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string
  
  const tableColumns: Ref<QTableProps['columns']> = ref([
    {name: 'timestamp', required: true, label: $i18n.t('modules.crm.adminer.settings.onec.sync.table.timestampColumn'), align: 'left', sortable: false, field: (row) => row.timestamp.toLocaleString()},
    {name: 'errorMessage', required: false, label: $i18n.t('modules.crm.adminer.settings.onec.sync.table.errorMessageColumn'), align: 'left', sortable: false, field: 'errorMessage'},
    {name: 'success', required: true, label: $i18n.t('modules.crm.adminer.settings.onec.sync.table.successColumn'), align: 'right', sortable: false, field: 'success'}
  ])
  const tableData = ref([] as Array<SyncEvent>)
  const dataLoading = ref(false)
  const loadingError = ref("")
  const tableRef = ref()
  const selectedEvent = ref(undefined as SyncEvent | undefined)
  const eventDetailsDialog = ref(false)


  const tablePagination = ref({
    page: 1,
    rowsPerPage: 10,
    rowsNumber: 0
  })
  
  function selectSyncEvent(_evt: Event, event: SyncEvent) {
    selectedEvent.value = event
    eventDetailsDialog.value = true
  }
  
  const syncing = ref(false)
  async function syncNow() {
    syncing.value = true
    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.crm.adminer.settings.onec.sync.table.syncOperationNotify')
    })
    try {
      const response = await api.crm.oneC.syncNow({ namespace: displayableNamespace })
      if (response.success) {
        notif({
            type: 'positive',
            message: $i18n.t('modules.crm.adminer.settings.onec.sync.table.syncSuccessNotify'),
            timeout: 5000
        })
      } else {
        notif({
            type: 'negative',
            message: $i18n.t('modules.crm.adminer.settings.onec.sync.table.syncFailNotify', { error: response.errorMessage }),
            timeout: 5000
        })
      }
    } catch (error) {
      notif({
        type: 'negative',
        message: $i18n.t('modules.crm.adminer.settings.onec.sync.table.syncOperationFailNotify', { error }),
        timeout: 5000
      })
    } finally {
        syncing.value = false
    }

    tableRef.value.requestServerInteraction()
  }

  async function loadData(tableProps: QTableProps) {
      const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
      const page = (tableProps.pagination?.page || 1) - 1
  
      const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.crm.adminer.settings.onec.sync.table.loadOperationNotify')
    })
      dataLoading.value = true
  
      try {
        const response = await api.crm.oneC.getSyncLog({ namespace: displayableNamespace, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.events
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
            message: $i18n.t('modules.crm.adminer.settings.onec.sync.table.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        dataLoading.value = false
      }
  }
  
  onMounted(() => {
    tableRef.value.requestServerInteraction()
  })
</script>
  
<style>
  
</style>