<template>
    <q-table
          flat
          bordered
          class="col-12 q-ma-none"
          :title="$t('modules.iot.fleet.deviceList.header')"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
          :loading="dataLoading"
          ref="tableRef"
          @request="loadData"
          dense
          v-model:pagination="tablePagination"
          @row-click="selectDevice"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.iot.fleet.deviceList.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.iot.fleet.deviceList.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:top-right>
          <q-btn
              v-if="props.editable"
              :label="$t('modules.iot.fleet.deviceList.createButton')"
              class="q-ma-none"
              unelevated
              outline
              color="positive"
              size="md"
              :disable="creationDialog"
              @click="creationDialog = true"
          />
        </template>
  
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn color="dark" outline label="" icon="menu">
              <q-menu>
                <q-list style="">
                    <q-item clickable v-close-popup @click="deviceUUIDToDelete = props.row.uuid; deletionDialog = true;">
                      <q-item-section class="text-dark">{{ $t('modules.iot.fleet.deviceList.actionsMenu.changeFleet') }}</q-item-section>
                    </q-item>
                  <q-item clickable v-close-popup @click="deviceUUIDToDelete = props.row.uuid; deletionDialog = true;">
                    <q-item-section class="text-negative">{{ $t('modules.iot.fleet.deviceList.actionsMenu.delete') }}</q-item-section>
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

        <q-dialog v-model="creationDialog">
          <DeviceCreateModal
              :namespace="props.namespace"
              :fleet="props.uuid"
              @created="onDeviceCreated"
          />
        </q-dialog>

        <q-dialog v-model="deletionDialog">
          <DeviceDeleteModal
              :namespace="props.namespace"
              :uuid="deviceUUIDToDelete"
              @deleted="onDeviceDeleted"
          />
        </q-dialog>
</template>

<script setup lang="ts">
import { QTableProps, useQuasar, QTable } from 'quasar';
import api from 'src/boot/api';
import { Device } from 'src/boot/api/iot/device';
import { computed, onMounted, ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';

import DeviceCreateModal from '../device/DeviceCreateModal.vue';
import DeviceDeleteModal from '../device/DeviceDeleteModal.vue';

const $i18n = useI18n()
const $q = useQuasar()

const dataLoading = ref(false)
const loadingError = ref('')
const tableData = ref([] as Array<Device>)
const tableColumns = computed(() => {
  const v = [
    { name: 'uuid', label: $i18n.t('modules.iot.fleet.deviceList.uuidColumn'), field: 'uuid', align: 'left', sortable: true },
    { name: 'name', label: $i18n.t('modules.iot.fleet.deviceList.nameColumn'), field: 'name', align: 'left', sortable: true },
    { name: 'description', label: $i18n.t('modules.iot.fleet.deviceList.descriptionColumn'), field: 'description', align: 'left', sortable: false },
  ] as QTableProps['columns']

  if ( v && props.editable) {
    v.push({ name: 'actions', label: $i18n.t('modules.iot.fleet.deviceList.actionsColumn'), field: 'actions', align: 'right', sortable: false })
  }

  return v
})
const tableRef = ref<QTable | null>(null)


const creationDialog = ref(false)
const deviceUUIDToDelete = ref('')
const deletionDialog = ref(false)

const props = defineProps<{
    namespace: string,
    uuid: string,
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
        message: $i18n.t('modules.iot.fleet.deviceList.loadOperationNotify.')
    })
      dataLoading.value = true
  
      try {
        console.log(props)

        const response = await api.iot.fleet.listDevices({ namespace: props.namespace, uuid: props.uuid, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.devices
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
            message: $i18n.t('modules.iot.fleet.deviceList.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        dataLoading.value = false
      }
  }

  async function selectDevice(_evt: Event, uuid: Device) {

  }

  function onDeviceCreated() {
    creationDialog.value = false
    tableRef.value?.requestServerInteraction()
  }

  function onDeviceDeleted() {
    deletionDialog.value = false
    tableRef.value?.requestServerInteraction()
  }

  onMounted(() => {
    tableRef.value?.requestServerInteraction()
  })
</script>