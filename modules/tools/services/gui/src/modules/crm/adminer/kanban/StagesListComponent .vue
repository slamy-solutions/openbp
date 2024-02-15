<template>
    <div>
      <q-table
        flat
        class="col-12 q-ma-none bg-transparent"
        :title="$t('modules.crm.adminer.kanban.stage.list.header')"
        :columns="tableColumns"
        :rows="stages"
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
              {{ $t("modules.crm.adminer.kanban.stage.list.noData") }}
            </span>
            <span v-else class="text-negative">
              {{
                $t("modules.crm.adminer.kanban.stage.list.failedToLoad", {
                  error: loadingError,
                })
              }}
            </span>
          </div>
        </template>
  
        <template v-slot:top-right>
          <q-btn
            :label="$t('modules.crm.adminer.kanban.stage.list.createButton')"
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
              v-if="shouldRenderUpButton(props.row)"
              class="q-ma-none q-mr-xs"
              unelevated
              round
              outline
              size="sm"
              :disable="deletionDialog || loading"
              icon="arrow_upward"
              @click.stop="() => moveStageUp(props.row)"
            />
            <q-btn
              v-if="shouldRenderDownButton(props.row)"
              class="q-ma-none q-mr-xs"
              unelevated
              round
              outline
              size="sm"
              :disable="deletionDialog || loading"
              icon="arrow_downward"
              @click.stop="() => moveStageDown(props.row)"
            />
            <q-btn
              class="q-ma-none"
              unelevated
              outline
              size="sm"
              :disable="deletionDialog || loading"
              icon="delete"
              @click.stop="
                deletionDialogStage = props.row;
                deletionDialog = true;
              "
            />
          </q-td>
        </template>
      </q-table>
      <q-dialog v-model="creationDialog">
        <StageCreateModal
          :namespace="props.namespace"
          :department="props.departmentUUID"
          @created="
            async () => {
              creationDialog = false;
              await loadStages();
            }
          "
        />
      </q-dialog>
    
      <q-dialog v-model="deletionDialog">
        <StageDeleteModal
          :namespace="props.namespace"
          :name="deletionDialogStage?.name || ''"
          :stage-u-u-i-d="deletionDialogStage?.uuid || ''"
          @deleted="
            async () => {
              deletionDialog = false;
              await loadStages();
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
  import { TicketStage } from "src/boot/api/crm/kanban"
  
  import StageCreateModal from "./StageCreateModal.vue";
  import StageDeleteModal from "./StageDeleteModal.vue";
  
  const props = defineProps<{
    namespace: string;
    departmentUUID: string;
    allowSelection?: boolean;
  }>();
  
  const emits = defineEmits<{
    (event: "selected", stage: TicketStage): void;
  }>();
  
  const $i18n = useI18n();
  const $q = useQuasar();
  
  const loading = ref(false);
  const stages = ref([] as Array<TicketStage>);
  const tableColumns = ref<QTableProps["columns"]>([
    {
      name: "name",
      required: true,
      label: $i18n.t("modules.crm.adminer.kanban.stage.list.nameColumn"),
      align: "left",
      sortable: false,
      field: "name",
    },
    {
      name: "actions",
      required: true,
      label: $i18n.t("modules.crm.adminer.kanban.stage.list.actionsColumn"),
      align: "right",
      sortable: false,
      field: "actions",
    },
  ]);
  const loadingError = ref("");
  
  const creationDialog = ref(false);
  
  const deletionDialog = ref(false);
  const deletionDialogStage = ref<TicketStage | null>(null);
  
  const selectedRow = ref<Array<TicketStage>>([]);
  
  async function loadStages() {
    if (loading.value) return;

    loading.value = true;
    try {
      const response = await api.crm.kanban.getTicketStages({
        departmentUUID: props.departmentUUID,
        namespace: props.namespace,
      });
      stages.value = response.stages;
    } catch (error) {
      $q.notify({
        type: "negative",
        message: $i18n.t("modules.crm.adminer.kanban.stage.list.loadFailNotify", {
          error,
        }),
        timeout: 5000,
      });
      console.error(error);
    } finally {
      loading.value = false;
    }
  }
  
  function shouldRenderUpButton(stage: TicketStage) {
    return stages.value.indexOf(stage) !== 0;
  }
  function shouldRenderDownButton(stage: TicketStage) {
    return stages.value.indexOf(stage) !== stages.value.length - 1;
  }
  async function moveStageUp(stage: TicketStage) {
    const index = stages.value.indexOf(stage);
    if (index === 0) return;
    loading.value = true
    try {
      await api.crm.kanban.swapTicketStagesPriorities({
        namespace: props.namespace,
        uuid1: stages.value[index - 1].uuid,
        uuid2: stage.uuid,
      })
    } catch (error) {
      $q.notify({
        type: "negative",
        message: $i18n.t("modules.crm.adminer.kanban.stage.list.moveFailNotify", {
          error,
        }),
        timeout: 5000,
      });
      console.error(error);
    } finally {
      loading.value = false
    }

    await loadStages();
  }
  async function moveStageDown(stage: TicketStage) {
    const index = stages.value.indexOf(stage);
    if (index === stages.value.length - 1) return;
    loading.value = true
    try {
      await api.crm.kanban.swapTicketStagesPriorities({
        namespace: props.namespace,
        uuid1: stage.uuid,
        uuid2: stages.value[index + 1].uuid,
      })
    } catch (error) {
      $q.notify({
        type: "negative",
        message: $i18n.t("modules.crm.adminer.kanban.stage.list.moveFailNotify", {
          error,
        }),
        timeout: 5000,
      });
      console.error(error);
    } finally {
      loading.value = false
    }

    await loadStages();
  }

  onMounted(async () => {
    await loadStages();
  });
  
  watch([props],
    async () => {
      await loadStages();
    }, { deep: true }
  );
  
  </script>
  