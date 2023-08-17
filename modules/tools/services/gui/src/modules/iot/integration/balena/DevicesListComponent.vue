<template>
    <q-table
          flat
          bordered
          class="col-12 q-ma-none"
          :title="$t('modules.iot.integration.balena.device.list.header')"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
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
              {{ $t('modules.iot.integration.balena.device.list.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.iot.integration.balena.device.list.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:body-cell-isOnline="props">
          <q-td :props="props">
              <q-icon name="wifi_off" v-if="!props.row.balenaData.isOnline" color="negative" size="xs" />
              <q-icon name="wifi" v-if="props.row.balenaData.isOnline" color="positive" size="xs" />
          </q-td>
        </template>

        <template v-slot:body-cell-cpuUsage="props">
          <q-td :props="props">
                <q-linear-progress
                    stripe
                    rounded
                    size="xl"
                    :value="props.row.balenaData.cpuUsage / 100"
                    :color="props.row.balenaData.cpuUsage < 70 ? 'info': (props.row.balenaData.cpuUsage < 90 ? 'warning': 'negative')"
                >
                <q-tooltip>
                    {{ String(props.row.balenaData.cpuUsage) + ' %' }}
                </q-tooltip>
                </q-linear-progress>
          </q-td>
        </template>

        <template v-slot:body-cell-ramUsage="props">
          <q-td :props="props">
                <q-linear-progress
                    stripe
                    rounded
                    size="xl"
                    :value="props.row.balenaData.memoryUsage / props.row.balenaData.memoryTotal"
                    :color="props.row.balenaData.memoryUsage / props.row.balenaData.memoryTotal < 0.7 ? 'info': (props.row.balenaData.memoryUsage / props.row.balenaData.memoryTotal < 0.9 ? 'warning': 'negative')"
                >
                <q-tooltip>
                    {{ String(props.row.balenaData.memoryUsage) + 'MB / ' + String(props.row.balenaData.memoryTotal) + 'MB' }}
                </q-tooltip>
                </q-linear-progress>
          </q-td>
        </template>

        <template v-slot:body-cell-lastConnectivity="props">
          <q-td :props="props">
              {{ props.row.balenaData.lastConnectivityEvent.toLocaleString('en-GB') }}
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn color="dark" outline label="" icon="menu" size="xs">
              <q-menu>
                <q-list style="">
                  <q-item clickable v-close-popup @click="deviceUUIDToBind = props.row.uuid; bindDialog = true;" v-if="props.row.bindedDeviceUUID === ''">
                    <q-item-section class="text-dark">{{ $t('modules.iot.integration.balena.device.list.actionsMenu.bind') }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </q-td>
        </template>
    </q-table>

    <q-dialog v-model="bindDialog">
        <DeviceBindModal
            :balena-device-u-u-i-d="deviceUUIDToBind"
            :namespace="props.namespace"
            @on-binded="bindDialog = false; refreshData();"
        />
    </q-dialog>
</template>

<script setup lang="ts">
import { QTableProps, useQuasar, QTable } from 'quasar';
import api from 'src/boot/api';
import { Device } from 'src/boot/api/iot/integration/balena/models';
import { computed, onBeforeUnmount, onMounted, ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';

import DeviceBindModal from './DeviceBindModal.vue';

const $i18n = useI18n()
const $q = useQuasar()

const dataLoading = ref(false)
const loadingError = ref('')
const tableData = ref([] as Array<Device>)
const tableColumns = computed(() => {
  const v = [
    //{ name: 'uuid', label: $i18n.t('modules.iot.integration.balena.device.list.uuidColumn'), field: 'uuid', align: 'left', sortable: false },
    { name: 'server', label: $i18n.t('modules.iot.integration.balena.device.list.serverColumn'), field: (r) => r.balenaServerUUID, align: 'left', sortable: false },
    { name: 'balenaUUID', label: $i18n.t('modules.iot.integration.balena.device.list.balenaUUIDColumn'), field: (r) => r.balenaData.uuid, align: 'left', sortable: false },
    { name: 'name', label: $i18n.t('modules.iot.integration.balena.device.list.nameColumn'), field: (r) => r.balenaData.deviceName, align: 'left', sortable: false },

    { name: 'bindedDeviceNamespace', label: $i18n.t('modules.iot.integration.balena.device.list.bindedDeviceNamespaceColumn'), field: "bindedDeviceNamespace", align: 'left', sortable: false },
    { name: 'bindedDeviceColumn', label: $i18n.t('modules.iot.integration.balena.device.list.bindedDeviceColumn'), field: "bindedDeviceUUID", align: 'left', sortable: false },

    { name: 'status', label: $i18n.t('modules.iot.integration.balena.device.list.statusColumn'), field: (r) => r.balenaData.status, align: 'left', sortable: false },
    { name: 'isOnline', label: $i18n.t('modules.iot.integration.balena.device.list.isOnlineColumn'), field: (r) => r.balenaData.isOnline, align: 'left', sortable: false },
    { name: 'cpuUsage', label: $i18n.t('modules.iot.integration.balena.device.list.cpuUsageColumn'), field: (r) => r.balenaData.cpuUsage, align: 'left', sortable: false },
    { name: 'cpuTemp', label: $i18n.t('modules.iot.integration.balena.device.list.cpuTempColumn'), field: (r) => r.balenaData.cpuTemp, align: 'left', sortable: false },
    { name: 'ramUsage', label: $i18n.t('modules.iot.integration.balena.device.list.ramUsageColumn'), field: (r) => r.balenaData.memoryUsage, align: 'left', sortable: false },
    { name: 'lastConnectivity', label: $i18n.t('modules.iot.integration.balena.device.list.lastConnectivityColumn'), field: (r) => r.balenaData.lastConnectivityEvent, align: 'left', sortable: false },
  ] as QTableProps['columns']

  if ( v && props.editable) {
    v.push({ name: 'actions', label: $i18n.t('modules.iot.integration.balena.device.list.actionsColumn'), field: 'actions', align: 'right', sortable: false })
  }

  return v
})
const tableRef = ref<QTable | null>(null)

const deviceUUIDToBind = ref("")
const bindDialog = ref(false)

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
        const response = await api.iot.integration.balena.device.list({ namespace: props.namespace, skip: rowsPerPage*page, limit: rowsPerPage, bindingFilter: 'all' })
        tableData.value = response.devices
        if (tablePagination.value != undefined) {
          tablePagination.value.page = page + 1
          tablePagination.value.rowsPerPage = rowsPerPage
          tablePagination.value.rowsNumber = response.totalCount
        }
        loadingError.value = ""

      } catch (error) {
        $q.notify({
            type: 'negative',
            message: $i18n.t('modules.iot.integration.balena.device.list.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        dataLoading.value = false
      }
  }

  async function selectServer(_evt: Event, device: Device) {

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