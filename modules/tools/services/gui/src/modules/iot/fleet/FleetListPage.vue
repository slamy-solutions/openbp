<template>
    <q-page class="q-pa-md row">
      <div class="col-2 q-pr-md">
        <MenuComponent selected="fleets"/>
      </div>
      <div class="col-10">
      <q-table
          class="row q-ma-none"
          :title="$t('modules.iot.fleet.list.header')"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
          :loading="dataLoading"
          ref="tableRef"
          @request="loadData"
          @row-click="selectFleet"
          :selected="selectedFleet ? [selectedFleet] : []"
          dense
          v-model:pagination="tablePagination"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.iot.fleet.list.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.iot.fleet.list.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:top-right>
          <q-btn
              :label="$t('modules.iot.fleet.list.createButton')"
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
                  <q-item clickable v-close-popup @click="fleetUUIDToDelete = props.row.uuid; deletionDialog = true;">
                    <q-item-section class="text-negative">{{ $t('modules.iot.fleet.list.actionsMenu.delete') }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </q-td>
        </template>
      </q-table>

      <div class="row full-width q-pt-md">
      <q-card class="full-width q-pt-xs">
       
        <FleetViewComponent

          :namespace="displayableNamespace"
          :update-possible="true"
          :uuid="selectedFleet?.uuid || ''"
        />
      </q-card>
    </div>

      </div>
  
      <q-dialog v-model="creationDialog">
          <FleetCreateModal :namespace="displayableNamespace" @created="onFleetCreated" />
      </q-dialog>
  
      <q-dialog v-model="deletionDialog">
          <FleetDeleteModal :namespace="displayableNamespace" :uuid="fleetUUIDToDelete" @deleted="onFleetDeleted" />
      </q-dialog>
    </q-page>
  </template>
  
  <script setup lang="ts">
  import { QPaginationProps, QTableProps, useQuasar } from 'quasar';
  import { onMounted, Ref, ref } from 'vue';
  import { useRoute } from 'vue-router';
  import { useI18n } from 'vue-i18n';
  import api from '../../../boot/api';
  
  import MenuComponent from '../MenuComponent.vue'
  import FleetCreateModal from './FleetCreateModal.vue';
  import FleetDeleteModal from './FleetDeleteModal.vue';
  import FleetViewComponent from './FleetViewComponent.vue';
  import { Fleet } from 'src/boot/api/iot/fleet';
  
  const $i18n = useI18n()
  const $q = useQuasar()
  
  const $route = useRoute()
  const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string
  
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
  const selectedFleet = ref(undefined as Fleet | undefined)


  const tablePagination = ref({
    page: 1,
    rowsPerPage: 10,
    rowsNumber: 0
  })
  
  const creationDialog = ref(false)
  const fleetUUIDToDelete = ref('')
  const deletionDialog = ref(false)
  
  function selectFleet(_evt: Event, fleet: Fleet) {
    selectedFleet.value = fleet
  }
  
  async function loadData(tableProps: QTableProps) {
      const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
      const page = (tableProps.pagination?.page || 1) - 1
  
      const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.iot.fleet.list.loadOperationNotify')
    })
      dataLoading.value = true
  
      try {
        const response = await api.iot.fleet.listFleets({ namespace: displayableNamespace, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.fleets
        if (tablePagination.value != undefined) {
          tablePagination.value.page = page + 1
          tablePagination.value.rowsPerPage = rowsPerPage
          tablePagination.value.rowsNumber = response.totalCount
        }
        notif()
        loadingError.value = ""
        if (response.fleets.length > 0) {
          selectedFleet.value = response.fleets[0]
        }
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
  
  <style>
  
  </style>