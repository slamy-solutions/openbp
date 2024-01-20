<template>
  <div>
    <q-table
      flat
      class="col-12 q-ma-none bg-transparent"
      :title="$t('modules.crm.adminer.performer.list.header')"
      :columns="tableColumns"
      :rows="performers"
      row-key="uuid"
      :loading="loading"
      ref="tableRef"
      @row-click="
        (_, row) => (props.allowSelection ? emits('selected', row) : null)
      "
      :selected="selectedRow"
      dense
      square
    >
      <template v-slot:loading>
        <q-inner-loading showing color="secondary" />
      </template>

      <template v-slot:no-data="{}">
        <div class="full-width row flex-center q-gutter-sm">
          <span v-if="loadingError === ''">
            {{ $t("modules.crm.adminer.performer.list.noData") }}
          </span>
          <span v-else class="text-negative">
            {{
              $t("modules.crm.adminer.performer.list.failedToLoad", {
                error: loadingError,
              })
            }}
          </span>
        </div>
      </template>

      <template v-slot:top-right>
        <q-btn
          :label="$t('modules.crm.adminer.performer.list.createButton')"
          class="q-ma-none"
          unelevated
          outline
          size="sm"
          :disable="creationDialog"
          @click.stop="creationDialog = true"
        />
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn
            class="q-ma-none"
            unelevated
            outline
            size="sm"
            :disable="deletionDialog"
            icon="delete"
            @click.stop="
              deletionDialogPerformer = props.row;
              deletionDialog = true;
            "
          />
        </q-td>
      </template>
    </q-table>
    <q-dialog v-model="creationDialog">
      <PerformerCreateModal
        :namespace="props.namespace"
        @created="
          async () => {
            creationDialog = false;
            await loadPerformers();
          }
        "
      />
    </q-dialog>
  
    <q-dialog v-model="deletionDialog">
      <PerformerDeleteModal
        :namespace="props.namespace"
        :performer="(deletionDialogPerformer as Performer)"
        @deleted="
          async () => {
            deletionDialog = false;
            await loadPerformers();
          }
        "
      />
    </q-dialog>
  
  </div>
</template>

<script setup lang="ts">
import { QTableProps, useQuasar } from "quasar";
import api from "../../../../boot/api";
import { onMounted, ref, defineProps, defineEmits, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";

import { Performer } from "src/boot/api/crm/performer";

import PerformerCreateModal from "./PerformerCreateModal.vue";
import PerformerDeleteModal from "./PerformerDeleteModal.vue";

const props = defineProps<{
  namespace: string;
  departmentUUID?: string;
  allowSelection?: boolean;
}>();

const emits = defineEmits<{
  (event: "selected", performer: Performer): void;
}>();

const $i18n = useI18n();
const $q = useQuasar();

const loading = ref(false);
const performers = ref([] as Array<Performer>);
const tableColumns = ref<QTableProps["columns"]>([
  {
    name: "name",
    required: true,
    label: $i18n.t("modules.crm.adminer.performer.list.nameColumn"),
    align: "left",
    sortable: false,
    field: "name",
  },
  {
    name: "actions",
    required: true,
    label: $i18n.t("modules.crm.adminer.performer.list.actionsColumn"),
    align: "right",
    sortable: false,
    field: "actions",
  },
]);
const loadingError = ref("");

const creationDialog = ref(false);

const deletionDialog = ref(false);
const deletionDialogPerformer = ref<Performer | null>(null);

const selectedRow = ref<Array<Performer>>([]);

async function loadPerformers() {
  loading.value = true;
  try {
    const response = await api.crm.performer.getAll({
      namespace: props.namespace,
    });
    performers.value = response.performers;
    if (props.departmentUUID) {
      performers.value = performers.value.filter(
        (performer) => performer.departmentUUID === props.departmentUUID
      );
    }
  } catch (error) {
    $q.notify({
      type: "negative",
      message: $i18n.t("modules.crm.adminer.performer.list.loadFailNotify", {
        error,
      }),
      timeout: 5000,
    });
    console.error(error);
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadPerformers();
});

watch([props],
  async () => {
    await loadPerformers();
  }
);

</script>
