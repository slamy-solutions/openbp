<template>
  <div class="row">
    <q-table
      flat
      class="col-12 q-ma-none bg-transparent"
      :title="$t('modules.crm.adminer.department.select.header')"
      :columns="tableColumns"
      :rows="departments"
      row-key="uuid"
      :loading="loading"
      ref="tableRef"
      dense
      selection="single"
      v-model:selected="selectedDepartment"
    >
      <template v-slot:loading>
        <q-inner-loading showing color="secondary" />
      </template>

      <template v-slot:no-data="{}">
        <div class="full-width row flex-center q-gutter-sm">
          <span v-if="loadingError === ''">
            {{ $t("modules.crm.adminer.department.select.noData") }}
          </span>
          <span v-else class="text-negative">
            {{
              $t("modules.crm.adminer.department.select.failedToLoad", {
                error: loadingError,
              })
            }}
          </span>
        </div>
      </template>
    </q-table>

    <div class="col-12 q-mt-md row">
      <q-btn color="negative" class="col-2" @click="emit('cancelled')">
        {{ $t("modules.crm.adminer.department.select.cancelButton") }}
      </q-btn>
      <div class="col-5"></div>
      <q-btn
        color="positive"
        class="col-5"
        :disable="!selectedDepartment || selectedDepartment.length === 0"
        @click="emit('selected', selectedDepartment[0])"
      >
        {{ $t("modules.crm.adminer.department.select.selectButton") }}
      </q-btn>
    </div>
  </div>
</template>

<script setup lang="ts">
import { QTableProps, useQuasar } from "quasar";
import api from "../../../../boot/api";
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

import { Department } from "../../../../boot/api/crm/department";

const $i18n = useI18n();
const $q = useQuasar();

const props = defineProps<{
  namespace: string;
}>();

const emit = defineEmits<{
  (e: "selected", department: Department): void;
  (e: "cancelled"): void;
}>();

const loading = ref(false);
const departments = ref([] as Array<Department>);
const tableColumns = ref<QTableProps["columns"]>([
  {
    name: "name",
    required: true,
    label: $i18n.t("modules.crm.adminer.department.select.nameColumn"),
    align: "left",
    sortable: false,
    field: "name",
  },
]);
const loadingError = ref("");

const selectedDepartment = ref([] as Department[]);

async function loadDepartments() {
  loading.value = true;
  try {
    const response = await api.crm.department.getAll({
      namespace: props.namespace,
    });
    departments.value = response.departments;
  } catch (error) {
    $q.notify({
      type: "negative",
      message: $i18n.t("modules.crm.adminer.department.select.loadFailNotify", {
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
  await loadDepartments();
});
</script>
