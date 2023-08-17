<template>
    <q-table
          flat
          bordered
          class="col-12 q-ma-none"
          :title="$t('modules.iot.integration.balena.server.list.header')"
          :columns="tableColumns"
          :rows="tableData"
          :row-key="(r) => r.server.uuid"
          :loading="dataLoading && tableData.length === 0"
          ref="tableRef"
          @request="loadData"
          dense
          v-model:pagination="tablePagination"
          @row-click="selectServer"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.iot.integration.balena.server.list.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.iot.integration.balena.server.list.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:top-right>
          <q-btn
              v-if="props.editable"
              :label="$t('modules.iot.integration.balena.server.list.createButton')"
              class="q-ma-none"
              unelevated
              outline
              color="positive"
              size="md"
              :disable="addDialog"
              @click="addDialog = true"
          />
        </template>
  
        <template v-slot:body-cell-enabled="props">
          <q-td :props="props">
            <q-toggle size="xs" v-model="props.row.server.enabled" @update:model-value="async (newVal) => await setServerEnabled(props.row.server, newVal)" />
          </q-td>
        </template>

        <template v-slot:body-cell-syncStatus="props">
          <q-td :props="props">
            <div v-if="props.row.lastSyncLog === null || props.row.lastSyncLog === undefined">-</div>
            <div v-else>
              {{ props.row.lastSyncLog.status + (props.row.lastSyncLog.message ? ` (${props.row.lastSyncLog.message})` : '') }}
            </div>
          </q-td>
        </template>

        <template v-slot:body-cell-syncDevices="props">
          <q-td :props="props">
            <div v-if="props.row.lastSyncLog === null || props.row.lastSyncLog === undefined">-</div>
            <div v-else>
              {{ props.row.lastSyncLog.stats.foundedDevicesOnServer }}
            </div>
          </q-td>
        </template>

        <template v-slot:body-cell-syncTime="props">
          <q-td :props="props">
            <div v-if="props.row.lastSyncLog === null || props.row.lastSyncLog === undefined">-</div>
            <div v-else>
              {{ String((props.row.lastSyncLog.stats.executionTime / 1000).toFixed(1)) + " s" }}
            </div>
          </q-td>
        </template>

        <template v-slot:body-cell-syncTimestamp="props">
          <q-td :props="props">
            <div v-if="props.row.lastSyncLog === null || props.row.lastSyncLog === undefined">-</div>
            <div v-else>
              {{ props.row.lastSyncLog.timestamp.toLocaleString('en-GB') }}
            </div>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn color="dark" outline label="" icon="menu" size="xs">
              <q-menu>
                <q-list style="">
                  <q-item clickable v-close-popup @click="serverUUIDToDelete = props.row.uuid; removingDialog = true;">
                    <q-item-section class="text-negative">{{ $t('modules.iot.integration.balena.server.list.actionsMenu.delete') }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </q-td>
        </template>
        </q-table>

        <q-dialog v-model="addDialog">
          <ServerAddModal
              :namespace="props.namespace"
              @added="onServerAdded"
          />
        </q-dialog>

        <q-dialog v-model="removingDialog">
          <ServerRemoveModal
              :namespace="props.namespace"
              :uuid="serverUUIDToDelete"
              @removed="onServerRemoved"
          />
        </q-dialog>
</template>

<script setup lang="ts">
import { QTableProps, useQuasar, QTable } from 'quasar';
import api from 'src/boot/api';
import { Server, SyncLogEntry } from 'src/boot/api/iot/integration/balena/models';
import { computed, onMounted, ref, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';

import ServerAddModal from './ServerAddModal.vue';

import ServerRemoveModal from './ServerRemoveModal.vue';

const $i18n = useI18n()
const $q = useQuasar()

const dataLoading = ref(false)
const loadingError = ref('')
const tableData = ref([] as Array<{
  server: Server
  lastSyncLog?: SyncLogEntry
}>)
const tableColumns = computed(() => {
  const v = [
    { name: 'uuid', label: $i18n.t('modules.iot.integration.balena.server.list.uuidColumn'), field: (r) => r.server.uuid, align: 'left', sortable: true },
    { name: 'name', label: $i18n.t('modules.iot.integration.balena.server.list.nameColumn'), field: (r) => r.server.name, align: 'left', sortable: true },
    { name: 'description', label: $i18n.t('modules.iot.integration.balena.server.list.descriptionColumn'), field: (r) => r.server.description, align: 'left', sortable: false },
    { name: 'enabled', label: $i18n.t('modules.iot.integration.balena.server.list.enabledColumn'), align: 'left', sortable: false },
    { name: 'syncStatus', label: $i18n.t('modules.iot.integration.balena.server.list.syncStatusColumn'), align: 'left', sortable: false },
    { name: 'syncDevices', label: $i18n.t('modules.iot.integration.balena.server.list.syncDevicesColumn'), align: 'left', sortable: false },
    { name: 'syncTime', label: $i18n.t('modules.iot.integration.balena.server.list.syncTimeColumn'), align: 'left', sortable: false },
    { name: 'syncTimestamp', label: $i18n.t('modules.iot.integration.balena.server.list.syncTimestampColumn'), align: 'left', sortable: false },
  ] as QTableProps['columns']

  if ( v && props.editable) {
    v.push({ name: 'actions', label: $i18n.t('modules.iot.integration.balena.server.list.actionsColumn'), field: 'actions', align: 'right', sortable: false })
  }

  return v
})
const tableRef = ref<QTable | null>(null)


const addDialog = ref(false)
const serverUUIDToDelete = ref('')
const removingDialog = ref(false)

const props = defineProps<{
    namespace: string,
    editable?: boolean
}>()

const tablePagination = ref({
    page: 1,
    rowsPerPage: 10,
    rowsNumber: 0
  })


  async function loadData(tableProps: QTableProps) {
      const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
      const page = (tableProps.pagination?.page || 1) - 1
  
      dataLoading.value = true
  
      try {
        const response = await api.iot.integration.balena.server.list({ namespace: props.namespace, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.servers
        if (tablePagination.value != undefined) {
          tablePagination.value.page = page + 1
          tablePagination.value.rowsPerPage = rowsPerPage
          tablePagination.value.rowsNumber = response.totalCount
        }
        loadingError.value = ""

      } catch (error) {
        $q.notify({
            type: 'negative',
            message: $i18n.t('modules.iot.integration.balena.server.list.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        dataLoading.value = false
      }
  }

  async function selectServer(_evt: Event, server: Server) {

  }

  async function setServerEnabled(server: Server, enabled: boolean) {
    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.iot.integration.balena.server.list.setEnabledOperationNotify')
    })

    try {
      await api.iot.integration.balena.server.setEnabled({ uuid: server.uuid, enabled })
      notif({
        type: 'positive',
        message: $i18n.t('modules.iot.integration.balena.server.list.setEnabledSuccessNotify'),
      })
      loadingError.value = ""
    } catch (error) {
      notif({
        type: 'negative',
        message: $i18n.t('modules.iot.integration.balena.server.list.setEnabledFailNotify', { error }),
        timeout: 5000
      })
      loadingError.value = String(error)
    } finally {
      tableRef.value?.requestServerInteraction()
    }
  }

  function onServerAdded() {
    addDialog.value = false
    tableRef.value?.requestServerInteraction()
  }

  function onServerRemoved() {
    removingDialog.value = false
    tableRef.value?.requestServerInteraction()
  }

  //TODO: should this be changed to websocket?
  const refreshInterval = ref(null as null | NodeJS.Timeout)

function refreshData() {
    tableRef.value?.requestServerInteraction()
}

onMounted(() => {
    refreshData()
    refreshInterval.value = setInterval(refreshData, 5000)
})

onBeforeUnmount(()=>{
    if (refreshInterval.value) {
        clearInterval(refreshInterval.value)
    }
})
</script>