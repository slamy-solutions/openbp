<template>
    <div class="row q-pa-md">
        <div>
            <h5 class="q-pa-none q-ma-none text-bold">{{ $t('modules.iot.fleet.select.header') }}</h5>
        </div>
      <q-table
          class="col-12 q-ma-none q-mt-md"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
          :loading="loading"
          ref="tableRef"
          @request="loadData"
          dense
          v-model:pagination="tablePagination"
          @row-click="selectFleet"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
          <template v-slot:body-cell="props">
            <q-td
            :props="props"
            :class="(props.row.uuid==selectedFleet?.uuid)?'bg-primary':'bg-white text-black'"
            >
            {{props.value}}
            </q-td>
        </template>

        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.iot.fleet.select.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.iot.fleet.select.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:body-cell-managed="props">
          <q-td
            :props="props"
            :class="(props.row.uuid==selectedFleet?.uuid)?'bg-primary':'bg-white text-black'"
        >
            <ManagedByComponent dense :managed-by="props.row.managed" />
          </q-td>
        </template>
      </q-table>
      <div class="col-12 q-mt-md row">
        <q-btn color="negative" class="col-2" @click="emit('canceled')">{{ $t('modules.iot.fleet.select.cancelButton') }}</q-btn>
        <div class="col-5"></div>
        <q-btn color="positive" class="col-5" :disable="selectedFleet == null" @click="emit('selected', selectedFleet as Fleet)">{{ $t('modules.iot.fleet.select.selectButton') }}</q-btn>
      </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, Ref } from 'vue'
import { QTableProps, useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n';
import api from 'src/boot/api';
import { Fleet } from 'src/boot/api/iot/fleet';

const props = defineProps<{
    namespace: string
}>()
const emit = defineEmits<{
  (e: 'canceled'): void,
  (e: 'selected', fleet: Fleet): void
}>()

const $i18n = useI18n()
const $q = useQuasar()

const loading = ref(false)
const loadingError = ref("")
const tableColumns: Ref<QTableProps['columns']> = ref([
    { name: 'uuid', label: $i18n.t('modules.iot.fleet.select.uuidColumn'), field: 'uuid', align: 'left', sortable: true },
    { name: 'name', label: $i18n.t('modules.iot.fleet.select.nameColumn'), field: 'name', align: 'left', sortable: true },
    { name: 'description', label: $i18n.t('modules.iot.fleet.select.descriptionColumn'), field: 'description', align: 'left', sortable: false },
])
const tableData = ref([] as Array<Fleet>)
const tableRef = ref()
const selectedFleet = ref(null as Fleet | null)

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
        message: $i18n.t('modules.iot.fleet.select.loadOperationNotify')
    })
      loading.value = true
  
      try {
        const response = await api.iot.fleet.listFleets({ namespace: props.namespace,   skip: rowsPerPage*page, limit: rowsPerPage })
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
            message: $i18n.t('modules.iot.fleet.select.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        loading.value = false
      }
}

function selectFleet(_evt: Event, fleet: Fleet) {
    selectedFleet.value = fleet
}

onMounted(() => {
    tableRef.value.requestServerInteraction()
})
</script>