<template>
<q-page class="q-pa-md row">
      <div class="col-2 q-pr-md">
        <MenuComponent selected="roles"/>
      </div>
      <div class="col-10 row">
        <q-table
          class="col-12 q-ma-none"
          :title="$t('modules.accessControl.iam.role.list.header')"
          :columns="tableColumns"
          :rows="tableData"
          row-key="uuid"
          :loading="dataLoading"
          ref="tableRef"
          @request="loadData"
          dense
          v-model:pagination="tablePagination"
          @row-click="selectRole"
      >
          <template v-slot:loading>
              <q-inner-loading showing color="secondary" />
          </template>
  
        <template v-slot:no-data="{}">
          <div class="full-width row flex-center q-gutter-sm">
            <span v-if="loadingError === ''">
              {{ $t('modules.accessControl.iam.role.list.noData') }}
            </span>
            <span v-else class="text-negative">
              {{ $t('modules.accessControl.iam.role.list.failedToLoad', { error: loadingError }) }}
            </span>
          </div>
        </template>
  
        <template v-slot:top-right>
          <q-btn
              :label="$t('modules.accessControl.iam.role.list.createButton')"
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
                  <q-item clickable v-close-popup @click="roleUUIDToDelete = props.row.uuid; deletionDialog = true;">
                    <q-item-section class="text-negative">{{ $t('modules.accessControl.iam.role.list.actionsMenu.delete') }}</q-item-section>
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
          <RoleCreateModal :namespace="displayableNamespace" @created="onRoleCreated" />
        </q-dialog>

        <q-dialog v-model="deletionDialog">
          <RoleDeleteModal :namespace="displayableNamespace" :uuid="roleUUIDToDelete" @deleted="onRoleDeleted" />
        </q-dialog>

        <q-card class="col-12 q-mt-md">
          <RoleViewComponent :namespace="displayableNamespace" :uuid="selectedUUID" update-possible />
        </q-card>
      </div>
</q-page>
</template>

<script setup lang="ts">
  import { QTableProps, useQuasar } from 'quasar';
import { Role } from '../../../../boot/api/accessControl/role';
import { onMounted, ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import MenuComponent from '../../MenuComponent.vue'
import api from '../../../../boot/api';

import ManagedByComponent from 'src/components/managedItem/ManagedByComponent.vue'
import RoleViewComponent from './RoleViewComponent.vue'
import RoleCreateModal from './RoleCreateModal.vue'
import RoleDeleteModal from './RoleDeleteModal.vue'

  const $i18n = useI18n()
  const $q = useQuasar()
  
  const displayableNamespace = ref("")

  const tableColumns: Ref<QTableProps['columns']> = ref([
      {name: 'uuid', required: true, label: $i18n.t('modules.accessControl.iam.role.list.uuidColumn'), align: 'left', sortable: false, field: 'uuid'},
      {name: 'name', required: true, label: $i18n.t('modules.accessControl.iam.role.list.nameColumn'), align: 'left', sortable: false, field: 'name'},
      {name: 'description', required: true, label: $i18n.t('modules.accessControl.iam.role.list.descriptionColumn'), align: 'left', sortable: false, field: 'description'},
      {name: 'managed', required: true, label: $i18n.t('modules.accessControl.iam.role.list.managedColumn'), align: 'left', sortable: false, field: 'managed'},
      {name: 'actions', required: false, label: $i18n.t('modules.accessControl.iam.role.list.actionsColumn'), align: 'right', sortable: false, field: 'actions'}
  ])
  const tableData = ref([] as Array<Role>)
  const dataLoading = ref(false)
  const loadingError = ref("")
  const tableRef = ref()
  
  const tablePagination = ref({
    page: 1,
    rowsPerPage: 10,
    rowsNumber: 0
  })
  
  const creationDialog = ref(false)
  const roleUUIDToDelete = ref('')
  const deletionDialog = ref(false)
  
  const selectedUUID = ref('')

  async function loadData(tableProps: QTableProps) {
      const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
      const page = (tableProps.pagination?.page || 1) - 1
  
      const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.accessControl.iam.role.list.loadOperationNotify')
    })
      dataLoading.value = true
  
      try {
        const response = await api.accessControl.role.list({ namespace: displayableNamespace.value, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.roles
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
            message: $i18n.t('modules.accessControl.iam.role.list.loadFailNotify', { error }),
            timeout: 5000
        })
        loadingError.value = String(error) 
      } finally {
        dataLoading.value = false
      }
  }
  
  async function onRoleCreated() {
    creationDialog.value = false
    tableRef.value.requestServerInteraction()
  }
  
  async function onRoleDeleted() {
    deletionDialog.value = false
    tableRef.value.requestServerInteraction()
  }
  
  onMounted(() => {
    tableRef.value.requestServerInteraction()
  })

  function selectRole(_evt: Event, role: Role) {
    selectedUUID.value = role.uuid
  }
</script>