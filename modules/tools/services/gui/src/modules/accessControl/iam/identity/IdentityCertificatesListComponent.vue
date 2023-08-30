<template>
    <div class="row">
    <q-table
        class="col-12 q-pl-md q-pr-md bg-transparent"
        flat
        :columns="tableColumns"
        :rows="tableData"
        dense
        row-key="uuid"
        :loading="loading"
        ref="tableRef"
        @request="loadData"
        v-model:pagination="tablePagination"
        @row-click="onCertificateSelected"
    >
        <template v-slot:loading>
            <q-inner-loading showing color="secondary" />
        </template>

      <template v-slot:no-data="{}">
        <div class="full-width row flex-center q-gutter-sm">
          <span v-if="loadingError === ''">
            {{ $t('modules.accessControl.iam.identity.view.certificatesList.noData') }}
          </span>
          <span v-else class="text-negative">
            {{ $t('modules.accessControl.iam.identity.view.certificatesList.failedToLoad', { error: loadingError }) }}
          </span>
        </div>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn color="dark" outline label="" icon="menu" size="sm">
            <q-menu>
              <q-list style="">
                <q-item clickable v-if="editable" v-close-popup @click="deleteCertificate(props.row.uuid)">
                  <q-item-section class="text-negative">{{ $t('modules.accessControl.iam.identity.view.certificatesList.actionsMenu.delete') }}</q-item-section>
                </q-item>
                <q-item clickable v-if="editable" v-close-popup @click="disableCertificate(props.row.uuid)">
                  <q-item-section class="text-negative">{{ $t('modules.accessControl.iam.identity.view.certificatesList.actionsMenu.disable') }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <q-btn v-if="editable" class="col-12 q-mt-sm full-width" :label="$t('modules.accessControl.iam.identity.view.certificatesList.registerAndGenerateButton')" color="dark" @click="creationDialog = true"></q-btn>
    </div>

    <q-dialog v-model="creationDialog">
      <CertificateRegisterModal
        :defaultIdentityUUID="props.identityUUID"
        :defaultNamespace="props.identityNamespace"
        @registered="onCertificateRergistered"
      />
    </q-dialog>
</template>

<script lang="ts" setup>
import { QTableProps, useQuasar } from 'quasar';
import api from 'src/boot/api';
import { Certificate } from 'src/boot/api/accessControl/auth/certificate';
import { Ref, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import CertificateRegisterModal from '../auth/certificate/CertificateRegisterModal.vue'

const $i18n = useI18n()
const $q = useQuasar()

const props = defineProps<{
    identityNamespace: string,
    identityUUID: string,
    editable?: boolean
}>()

const tableRef = ref()
const tableColumns: Ref<QTableProps['columns']> = ref([
    {name: 'uuid', required: true, label: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.uuidColumn'), align: 'left', sortable: false, field: 'uuid'},
    {name: 'description', required: true, label: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.descriptionColumn'), align: 'left', sortable: false, field: 'description'},
    {name: 'disabled', required: true, label: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.disabledColumn'), align: 'left', sortable: false, field: 'disabled'},
    {name: 'actions', required: false, label: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.actionsColumn'), align: 'right', sortable: false, field: 'actions'}
])
const tableData = ref([] as Array<Certificate>)
const tablePagination = ref({
  page: 1,
  rowsPerPage: 10,
  rowsNumber: 0
})


const creationDialog = ref(false)


const loading = ref(false)
const loadingError = ref("")

function onCertificateSelected(_evt: Event, cert: Certificate) {}

async function loadData(tableProps: QTableProps) {
    const rowsPerPage = tableProps.pagination?.rowsPerPage || 100
    const page = (tableProps.pagination?.page || 1) - 1

    loading.value = true
    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.loadOperationNotify')
    })
    try {
        const response = await api.accessControl.auth.certificate.listForIdentity({ namespace: props.identityNamespace, identityUUID: props.identityUUID, skip: rowsPerPage*page, limit: rowsPerPage })
        tableData.value = response.certificates
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
          message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.loadFailNotify', { error }),
          timeout: 5000
        })
        loadingError.value = String(error)
    } finally {
        loading.value = false
    }
}

async function deleteCertificate(uuid: string) {
  const notif = $q.notify({
    type: 'ongoing',
    message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.deleteOperationNotify')
  })

  try {
    await api.accessControl.auth.certificate.deleteCertificate({
      namespace: props.identityNamespace,
      certificateUUID: uuid
    })
    notif({
      type: 'positive',
      message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.deleteSuccessNotify'),
    })
    tableRef.value.requestServerInteraction()
  } catch (error) {
    notif({
      type: 'negative',
      message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.deleteFailNotify', { error }),
      timeout: 5000
    })
  }
}

async function disableCertificate(uuid: string) {
  const notif = $q.notify({
    type: 'ongoing',
    message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.disableOperationNotify')
  })

  try {
    await api.accessControl.auth.certificate.disable({
      namespace: props.identityNamespace,
      certificateUUID: uuid
    })
    notif({
      type: 'positive',
      message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.disableSuccessNotify'),
    })
    tableRef.value.requestServerInteraction()
  } catch (error) {
    notif({
      type: 'negative',
      message: $i18n.t('modules.accessControl.iam.identity.view.certificatesList.disableFailNotify', { error }),
      timeout: 5000
    })
  }
}

function onCertificateRergistered(cert: Certificate) {
  tableData.value.unshift(cert)
  creationDialog.value = false
  tablePagination.value.rowsNumber += 1
}

onMounted(() => {
  tableRef.value.requestServerInteraction()
})

watch(props, ()=>{
  tableRef.value.requestServerInteraction()
})

</script>