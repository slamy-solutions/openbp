<template>
  <q-page class="q-pa-md">
    <q-table
        :title="$t('modules.namespace.list.table.header')"
        :columns="tableColumns"
        :rows="tableData"
        row-key="name"
        :loading="dataLoading"
        :filter="filter"
    >
        <template v-slot:loading>
            <q-inner-loading showing color="secondary" />
        </template>

      <template v-slot:no-data="{}">
        <div class="full-width row flex-center q-gutter-sm">
          <span v-if="loadingError === ''">
            {{ $t('modules.namespace.list.table.noData') }}
          </span>
          <span v-else class="text-negative">
            {{ $t('modules.namespace.list.table.failedToLoad', { error: loadingError }) }}
          </span>
        </div>
      </template>

      <template v-slot:top-right>
        <q-btn
            :label="$t('modules.namespace.list.table.createButton')"
            class="q-ma-none q-mr-lg"
            unelevated
            outline
            color="positive"
            size="lg"
            :disable="creationDialog"
            @click="creationDialog = true"
        />
        <q-input outlined debounce="100" v-model="filter" :placeholder="$t('modules.namespace.list.table.search')">
          <template v-slot:append>
            <q-icon name="search" size="sm"/>
          </template>
        </q-input>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn color="dark" outline label="" icon="menu">
            <q-menu>
              <q-list style="">
                <q-item clickable v-close-popup @click="namespaceToDelete = props.row.name; deletionDialog = true;">
                  <q-item-section class="text-negative">{{ $t('modules.namespace.list.table.actionsMenu.delete') }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <q-dialog v-model="creationDialog">
        <NamespaceCreateModal @created="onNamespaceCreated"></NamespaceCreateModal>
    </q-dialog>

    <q-dialog v-model="deletionDialog">
        <NamespaceDeleteModal :namespace-name="namespaceToDelete" @deleted="onNamespaceDeleted"></NamespaceDeleteModal>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { QTableProps, useQuasar } from 'quasar';
import { onMounted, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import api from '../../boot/api';
import { Namespace } from '../../boot/api/namespace/models';

import NamespaceCreateModal from './NamespaceCreateModal.vue'
import NamespaceDeleteModal from './NamespaceDeleteModal.vue'

const $i18n = useI18n()
const $q = useQuasar()

const tableColumns: Ref<QTableProps['columns']> = ref([
    {name: 'name', required: true, label: $i18n.t('modules.namespace.list.table.nameColumn'), align: 'left', sortable: true, field: 'name'},
    {name: 'fullName', required: false, label: $i18n.t('modules.namespace.list.table.fullNameColumn'), align: 'left', sortable: true, field: 'fullName'},
    {name: 'description', required: false, label: $i18n.t('modules.namespace.list.table.descriptionColumn'), align: 'left', sortable: false, field: 'description'},
    {name: 'actions', required: false, label: $i18n.t('modules.namespace.list.table.actionsColumn'), align: 'right', sortable: false, field: 'actions'}
])
const tableData = ref([] as Array<Namespace>)
const dataLoading = ref(false)
const loadingError = ref("")

const filter = ref("")

const creationDialog = ref(false)
const namespaceToDelete = ref('')
const deletionDialog = ref(false)

async function loadData() {
    if (dataLoading.value) {
        return
    }

    dataLoading.value = true
    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.namespace.list.table.loadOperationNotify')
    })

    try {
        tableData.value = await api.namespace.list.list()
        loadingError.value = ""
        notif()
    } catch (e) {
        loadingError.value = String(e)
        notif({
          type: 'negative',
          message: $i18n.t('modules.namespace.list.table.loadFailNotify'),
          timeout: 5000
        })
    } finally {
        dataLoading.value = false
    }
}

function onNamespaceCreated(namespace: Namespace) {
    creationDialog.value = false
    tableData.value.unshift(namespace)
}

function onNamespaceDeleted(name: string) {
  deletionDialog.value = false
  tableData.value = tableData.value.filter((item) => item.name != name)
}

onMounted(loadData)

</script>

<style>

</style>