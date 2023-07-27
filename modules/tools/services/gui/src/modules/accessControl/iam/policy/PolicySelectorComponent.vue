<template>
    <div class="row q-pa-md" style="max-width: 95%">
        <div>
            <h5 class="q-pa-none q-ma-none text-bold">{{ $t('modules.accessControl.iam.policy.select.header') }}</h5>
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
          @row-click="selectPolicy"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
          <template v-slot:body-cell="props">
            <q-td
            :props="props"
            :class="(props.row.uuid==selectedUUID)?'bg-primary':'bg-white text-black'"
            >
            {{props.value}}
            </q-td>
        </template>

        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.accessControl.iam.policy.select.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.accessControl.iam.policy.select.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:body-cell-managed="props">
          <q-td
            :props="props"
            :class="(props.row.uuid==selectedUUID)?'bg-primary':'bg-white text-black'"
        >
            <ManagedByComponent dense :managed-by="props.row.managed" />
          </q-td>
        </template>
      </q-table>
      <div class="col-12 q-mt-md row">
        <q-btn color="negative" class="col-2" @click="emit('canceled')">{{ $t('modules.accessControl.iam.policy.select.cancelButton') }}</q-btn>
        <div class="col-5"></div>
        <q-btn color="positive" class="col-5" :disable="selectedUUID == ''" @click="emit('selected', props.namespace, selectedUUID)">{{ $t('modules.accessControl.iam.policy.select.selectButton') }}</q-btn>
      </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, Ref } from 'vue'
import { QTableProps, useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n';
import { Policy } from 'src/boot/api/accessControl/policy';
import api from 'src/boot/api';

const props = defineProps<{
    namespace: string
}>()
const emit = defineEmits<{
  (e: 'canceled'): void,
  (e: 'selected', namespace: string, uuid: string): void
}>()

const $i18n = useI18n()
const $q = useQuasar()

const loading = ref(false)
const loadingError = ref("")
const tableColumns: Ref<QTableProps['columns']> = ref([
    {name: 'uuid', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.uuidColumn'), align: 'left', sortable: false, field: 'uuid'},
    {name: 'name', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.nameColumn'), align: 'left', sortable: false, field: 'name'},
    {name: 'description', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.descriptionColumn'), align: 'left', sortable: false, field: 'description'},
    {name: 'managed', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.managedColumn'), align: 'left', sortable: false, field: 'managed'},
])
const tableData = ref([] as Array<Policy>)
const tableRef = ref()
const selectedUUID = ref("")

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
        message: $i18n.t('modules.accessControl.iam.policy.select.loadOperationNotify')
    })
      loading.value = true
  
      try {
        const response = await api.accessControl.policy.list({ namespace: props.namespace, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.policies
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
            message: $i18n.t('modules.accessControl.iam.policy.select.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        loading.value = false
      }
}

function selectPolicy(_evt: Event, policy: Policy) {
    selectedUUID.value = policy.uuid
}

onMounted(() => {
    tableRef.value.requestServerInteraction()
})
</script>