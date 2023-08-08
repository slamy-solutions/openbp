<template>
    <q-table
          flat
          bordered
          class="col-12 q-ma-none"
          :title="$t('modules.iot.integration.balena.server.header')"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
          :loading="dataLoading"
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
              {{ $t('modules.iot.integration.balena.server.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.iot.integration.balena.server.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:top-right>
          <q-btn
              v-if="props.editable"
              :label="$t('modules.iot.integration.balena.server.createButton')"
              class="q-ma-none"
              unelevated
              outline
              color="positive"
              size="md"
              :disable="addDialog"
              @click="addDialog = true"
          />
        </template>
  
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn color="dark" outline label="" icon="menu">
              <q-menu>
                <q-list style="">
                    <q-item clickable v-close-popup @click="serverUUIDToDelete = props.row.uuid; removingDialog = true;">
                      <q-item-section class="text-dark">{{ $t('modules.iot.integration.balena.server.actionsMenu.changeFleet') }}</q-item-section>
                    </q-item>
                  <q-item clickable v-close-popup @click="serverUUIDToDelete = props.row.uuid; removingDialog = true;">
                    <q-item-section class="text-negative">{{ $t('modules.iot.integration.balena.server.actionsMenu.delete') }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </q-td>
        </template>
  
        <template v-slot:body-cell-managed="props">
          <q-td :props="props">
            <ManagedByComponent dense :managed-by="props.row.managed" />
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
import { computed, onMounted, ref, Ref } from 'vue';
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
    { name: 'uuid', label: $i18n.t('modules.iot.integration.balena.server.uuidColumn'), field: 'uuid', align: 'left', sortable: true },
    { name: 'name', label: $i18n.t('modules.iot.integration.balena.server.nameColumn'), field: 'name', align: 'left', sortable: true },
    { name: 'description', label: $i18n.t('modules.iot.integration.balena.server.descriptionColumn'), field: 'description', align: 'left', sortable: false },
    { name: 'enabled', label: $i18n.t('modules.iot.integration.balena.server.enabledColumn'), field: 'enabled', align: 'left', sortable: false },
  ] as QTableProps['columns']

  if ( v && props.editable) {
    v.push({ name: 'actions', label: $i18n.t('modules.iot.integration.balena.server.actionsColumn'), field: 'actions', align: 'right', sortable: false })
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
  
      const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.iot.integration.balena.server.loadOperationNotify.')
    })
      dataLoading.value = true
  
      try {
        const response = await api.iot.integration.balena.server.list({ namespace: props.namespace, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.servers
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
            message: $i18n.t('modules.iot.integration.balena.server.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        dataLoading.value = false
      }
  }

  async function selectServer(_evt: Event, server: Server) {

  }

  function onServerAdded() {
    addDialog.value = false
    tableRef.value?.requestServerInteraction()
  }

  function onServerRemoved() {
    removingDialog.value = false
    tableRef.value?.requestServerInteraction()
  }

  onMounted(() => {
    tableRef.value?.requestServerInteraction()
  })
</script>