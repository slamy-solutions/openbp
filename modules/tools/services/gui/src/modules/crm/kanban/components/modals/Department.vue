<template>
  <q-card
    class="bg-secondary"
    style="width: 400px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="updateDepartment"
  >
    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="department.name"
          label="Назва"
        />
      </div>
    </q-card-section>

    <q-card-actions class="q-pt-none row justify-between">
      <q-btn
        icon="mdi-arrow-left-circle-outline"
        style="font-size: 16px"
        label="До списку"
        v-close-popup
      />
      <q-btn
        icon="mdi-content-save-outline"
        style="font-size: 16px"
        label="Зберегти"
        @click="updateDepartment"
        :disable="updateDepartmentRequest"
        :loading="updateDepartmentRequest"
      />
    </q-card-actions>
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import { defineComponent, PropType, ref } from 'vue';
import { useDepartmentStore, IDepartment } from 'src/stores/department-store';

export default defineComponent({
  props: {
    editDepartment: {
      type: Object as PropType<IDepartment>,
      required: true,
    },
  },

  beforeMount() {
    this.department = { ...this.editDepartment };
  },

  methods: {
    async updateDepartment() {
      this.updateDepartmentRequest = true;

      try {
        await this.departmentStore.updateDepartment(this.department);
        this.$emit('changed');
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.updateDepartmentRequest = false;
    },
  },
  setup() {
    const updateDepartmentRequest = ref(false);

    return {
      departmentStore: useDepartmentStore(),
      department: ref({} as IDepartment),
      updateDepartmentRequest,
    };
  },
});
</script>
