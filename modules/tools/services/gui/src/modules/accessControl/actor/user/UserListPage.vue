<template>
  <q-page class="q-pa-md row">
    <div class="col-2 q-pr-md">
      <MenuComponent selected="users"/>
    </div>
    <div class="col-10">
    <q-table
        class="q-ma-none"
        :title="$t('modules.accessControl.iam.actor.user.list.header')"
        :columns="tableColumns"
        :rows="tableData"
        row-key="uuid"
        :loading="dataLoading"
        ref="tableRef"
        @request="loadData"
        dense
        v-model:pagination="tablePagination"
    >
        <template v-slot:loading>
            <q-inner-loading showing color="secondary" />
        </template>

      <template v-slot:no-data="{}">
        <div class="full-width row flex-center q-gutter-sm">
          <span v-if="loadingError === ''">
            {{ $t('modules.accessControl.iam.actor.user.list.noData') }}
          </span>
          <span v-else class="text-negative">
            {{ $t('modules.accessControl.iam.actor.user.list.failedToLoad', { error: loadingError }) }}
          </span>
        </div>
      </template>

      <template v-slot:top-right>
        <q-btn
            :label="$t('modules.accessControl.iam.actor.user.list.createButton')"
            class="q-ma-none"
            unelevated
            outline
            color="positive"
            size="md"
            :disable="creationDialog"
            @click="creationDialog = true"
        />
      </template>

      <template v-slot:body-cell-identity="props">
        <q-td :props="props">
          {{ props.row.identity }}
          <q-btn color="dark" outline round size="sm" icon="info" @click="userIdentityUUIDToShow = props.row.identity; userIdentityInfoDialog=true" />
        </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn color="dark" outline label="" icon="menu">
            <q-menu>
              <q-list style="">
                <q-item clickable v-close-popup @click="userUUIDToDelete = props.row.uuid; deletionDialog = true;">
                  <q-item-section class="text-negative">{{ $t('modules.accessControl.iam.actor.user.list.actionsMenu.delete') }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </template>
    </q-table>
    </div>

    <q-dialog v-model="userIdentityInfoDialog" >
      <div class="bg-primary" style="max-width: 95%; width: 95%;">
        <IdentityViewComponent :namespace="displayableNamespace" :uuid="userIdentityUUIDToShow" />
      </div>
    </q-dialog>

    <q-dialog v-model="creationDialog">
        <UserCreateModal :namespace="displayableNamespace" @created="onUserCreated" />
    </q-dialog>

    <q-dialog v-model="deletionDialog">
        <UserDeleteModal :namespace="displayableNamespace" :uuid="userUUIDToDelete" @deleted="onUserDeleted" />
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { QPaginationProps, QTableProps, useQuasar } from 'quasar';
import { onMounted, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import api from '../../../../boot/api';
import { Identity } from '../../../../boot/api/accessControl/identity';
import { Namespace } from '../../../../boot/api/namespace/models';

import MenuComponent from '../../MenuComponent.vue'
import ManagedByComponent from 'src/components/managedItem/ManagedByComponent.vue';
import { User } from 'src/boot/api/accessControl/actor/user';
import UserCreateModal from './UserCreateModal.vue';
import UserDeleteModal from './UserDeleteModal.vue';
import IdentityViewComponent from '../../iam/identity/IdentityViewComponent.vue';

const $i18n = useI18n()
const $q = useQuasar()

const displayableNamespace = ref("")

const tableColumns: Ref<QTableProps['columns']> = ref([
    {name: 'uuid', required: true, label: $i18n.t('modules.accessControl.iam.actor.user.list.uuidColumn'), align: 'left', sortable: false, field: 'uuid'},
    {name: 'login', required: true, label: $i18n.t('modules.accessControl.iam.actor.user.list.loginColumn'), align: 'left', sortable: false, field: 'login'},
    {name: 'fullName', required: true, label: $i18n.t('modules.accessControl.iam.actor.user.list.fullNameColumn'), align: 'left', sortable: false, field: 'fullName'},
    {name: 'identity', required: true, label: $i18n.t('modules.accessControl.iam.actor.user.list.identityColumn'), align: 'left', sortable: false, field: 'identity'},
    {name: 'actions', required: false, label: $i18n.t('modules.accessControl.iam.actor.user.list.actionsColumn'), align: 'right', sortable: false, field: 'actions'}
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

const userIdentityUUIDToShow = ref('')
const userIdentityInfoDialog = ref(false)
const creationDialog = ref(false)
const userUUIDToDelete = ref('')
const deletionDialog = ref(false)


async function loadData(tableProps: QTableProps) {
    const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
    const page = (tableProps.pagination?.page || 1) - 1

    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.actor.user.list.loadOperationNotify')
  })
    dataLoading.value = true

    try {
      const response = await api.accessControl.actor.user.list({ namespace: displayableNamespace.value, skip: rowsPerPage*page, limit: rowsPerPage })
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
          message: $i18n.t('modules.accessControl.iam.actor.user.list.loadFailNotify', { error }),
          timeout: 5000
      })
      loadingError.value = String(error) 
    } finally {
      dataLoading.value = false
    }
}

async function onUserCreated() {
  creationDialog.value = false
  tableRef.value.requestServerInteraction()
}

async function onUserDeleted() {
  deletionDialog.value = false
  tableRef.value.requestServerInteraction()
}

onMounted(() => {
  tableRef.value.requestServerInteraction()
})
</script>

<style>

</style>