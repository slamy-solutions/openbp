<template>
    <div class="row">
      <q-table
          :title="$t('modules.accessControl.iam.actor.user.select.header')"
          class="q-ma-none bg-transparent col-12"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
          :loading="dataLoading"
          ref="tableRef"
          @request="loadData"
          dense
          flat
          v-model:pagination="tablePagination"
          selection="single"
          v-model:selected="selectedRow"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.accessControl.iam.actor.user.select.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.accessControl.iam.actor.user.select.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
      </q-table>

      <div class="col-12 q-mt-md row">
        <q-btn color="negative" class="col-2" @click="emit('cancelled')">{{ $t('modules.accessControl.iam.actor.user.select.cancelButton') }}</q-btn>
        <div class="col-5"></div>
        <q-btn color="positive" class="col-5" :disable="!selectedRow || selectedRow.length === 0" @click="emit('selected', selectedRow[0])">{{ $t('modules.accessControl.iam.actor.user.select.selectButton') }}</q-btn>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { QTableProps, useQuasar } from 'quasar';
  import { onMounted, Ref, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import api from '../../../../boot/api';
  
  import { User } from 'src/boot/api/accessControl/actor/user';
  
  const $i18n = useI18n()
  const $q = useQuasar()
  
  const props = defineProps<{
    namespace: string
  }>()

  const emit = defineEmits<{
    (e: 'selected', user: User): void
    (e: 'cancelled'): void
  }>()
  
  const tableColumns: Ref<QTableProps['columns']> = ref([
    {name: 'login', required: true, label: $i18n.t('modules.accessControl.iam.actor.user.select.loginColumn'), align: 'left', sortable: false, field: 'login'},
      {name: 'fullName', required: true, label: $i18n.t('modules.accessControl.iam.actor.user.select.fullNameColumn'), align: 'left', sortable: false, field: 'fullName'},
  ])
  const tableData = ref([] as Array<User>)
  const dataLoading = ref(false)
  const loadingError = ref("")
  const tableRef = ref()

  const tablePagination = ref({
    page: 1,
    rowsPerPage: 10,
    rowsNumber: 0
  })

  const selectedRow = ref([] as User[])
  
  async function loadData(tableProps: QTableProps) {
      const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
      const page = (tableProps.pagination?.page || 1) - 1
  
      const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.actor.user.select.loadOperationNotify')
    })
      dataLoading.value = true
  
      try {
        const response = await api.accessControl.actor.user.list({ namespace: props.namespace, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.users
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
            message: $i18n.t('modules.accessControl.iam.actor.user.select.loadFailNotify', { error }),
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