<template>
    <q-page class="q-pa-md row">
      <div class="col-2 q-pr-md">
        <MenuComponent selected="policies"/>
      </div>
      <div class="col-10 row">
      <q-table
          class="col-12 q-ma-none bg-transparent"
          flat
          :title="$t('modules.accessControl.iam.policy.list.header')"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
          :loading="dataLoading"
          ref="tableRef"
          @request="loadData"
          dense
          v-model:pagination="tablePagination"
          @row-click="selectPolicy"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.accessControl.iam.policy.list.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.accessControl.iam.policy.list.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:top-right>
          <q-btn
              :label="$t('modules.accessControl.iam.policy.list.createButton')"
              class="q-ma-none"
              unelevated
              outline
              size="md"
              :disable="creationDialog"
              @click="creationDialog = true"
          />
        </template>
  
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn color="dark" outline label="" icon="menu" size="sm">
              <q-menu>
                <q-list style="">
                  <q-item clickable v-close-popup @click="identityUUIDToDelete = props.row.uuid; deletionDialog = true;">
                    <q-item-section class="text-negative">{{ $t('modules.accessControl.iam.policy.list.actionsMenu.delete') }}</q-item-section>
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

      <q-card
        class="col-12 q-mt-md bg-transparent" flat
      >
        <PolicyViewComponent :namespace="displayableNamespace" :uuid="selectedUUID" update-possible></PolicyViewComponent>
      </q-card>

      </div>

      <q-dialog v-model="creationDialog">
          <PolicyCreateModal :namespace="displayableNamespace" @created="onIdentityCreated"></PolicyCreateModal>
      </q-dialog>
  
      <q-dialog v-model="deletionDialog">
          <PolicyDeleteModal :namespace="displayableNamespace" :uuid="identityUUIDToDelete" @deleted="onIdentityDeleted"></PolicyDeleteModal>
      </q-dialog>
    </q-page>
  </template>
  
  <script setup lang="ts">
  import { QPaginationProps, QTableProps, useQuasar } from 'quasar';
  import { onMounted, Ref, ref } from 'vue';
  import { useRoute } from 'vue-router';
  import { useI18n } from 'vue-i18n';
  import api from '../../../../boot/api';
  import { Namespace } from '../../../../boot/api/namespace/models';
  
  import PolicyCreateModal from './PolicyCreateModal.vue'
  import PolicyDeleteModal from './PolicyDeleteModal.vue'
  import PolicyViewComponent from './PolicyViewComponent.vue'
  import MenuComponent from '../../MenuComponent.vue'
  import { Policy } from 'src/boot/api/accessControl/policy';
  import ManagedByComponent from '../../../../components/managedItem/ManagedByComponent.vue'
  
  const $i18n = useI18n()
  const $q = useQuasar()
  const $route = useRoute()
  const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string
  
  const tableColumns: Ref<QTableProps['columns']> = ref([
      {name: 'uuid', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.uuidColumn'), align: 'left', sortable: false, field: 'uuid'},
      {name: 'name', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.nameColumn'), align: 'left', sortable: false, field: 'name'},
      {name: 'description', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.descriptionColumn'), align: 'left', sortable: false, field: 'description'},
      {name: 'managed', required: true, label: $i18n.t('modules.accessControl.iam.policy.list.managedColumn'), align: 'left', sortable: false, field: 'managed'},
      {name: 'actions', required: false, label: $i18n.t('modules.accessControl.iam.policy.list.actionsColumn'), align: 'right', sortable: false, field: 'actions'}
  ])
  const tableData = ref([] as Array<Policy>)
  const dataLoading = ref(false)
  const loadingError = ref("")
  const tableRef = ref()
  
  const tablePagination = ref({
    page: 1,
    rowsPerPage: 10,
    rowsNumber: 0
  })
  
  const creationDialog = ref(false)
  const identityUUIDToDelete = ref('')
  const deletionDialog = ref(false)
  
  const selectedUUID = ref('')

  async function loadData(tableProps: QTableProps) {
      const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
      const page = (tableProps.pagination?.page || 1) - 1
  
      const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.policy.list.loadOperationNotify')
    })
      dataLoading.value = true
  
      try {
        const response = await api.accessControl.policy.list({ namespace: displayableNamespace, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.policies
        if (tablePagination.value != undefined) {
          tablePagination.value.page = page + 1
          tablePagination.value.rowsPerPage = rowsPerPage
          tablePagination.value.rowsNumber = response.totalCount
        }
        notif()
        loadingError.value = ""

        if (selectedUUID.value === '' && tableData.value.length > 0) {
          selectedUUID.value = tableData.value[0].uuid
        }

      } catch (error) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.accessControl.iam.policy.list.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        dataLoading.value = false
      }
  }
  
  async function onIdentityCreated() {
    creationDialog.value = false
    tableRef.value.requestServerInteraction()
  }
  
  async function onIdentityDeleted() {
    deletionDialog.value = false
    tableRef.value.requestServerInteraction()
  }
  
  onMounted(() => {
    tableRef.value.requestServerInteraction()
  })

  function selectPolicy(_evt: Event, policy: Policy) {
    selectedUUID.value = policy.uuid
  }
  
  </script>
  
  <style>
  
  </style>