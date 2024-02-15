<template>
  <q-card
    class="bg-secondary"
    style="width: 400px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="updateProject"
  >
    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="project.name"
          label="Назва проекту"
        />
        <q-select
          outlined
          map-options
          option-label="name"
          option-value="id"
          v-model="project.department"
          :options="departmentStore.departments"
          label-color="light-blue-6"
          label="Команда"
          color="light-blue-6"
          class="text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        />
        <q-toggle
          v-model="project.notRelevant"
          color="positive"
          size="54px"
          icon="mdi-timer-outline"
          label="Проект не актуальний"
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
        @click="updateProject"
        :disable="updateProjectRequest"
        :loading="updateProjectRequest"
      />
    </q-card-actions>
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import { defineComponent, PropType, ref } from 'vue';
import { useProjectStore, IProject } from 'src/stores/project-store';
import { useDepartmentStore } from 'src/stores/department-store';

export default defineComponent({
  props: {
    editProject: {
      type: Object as PropType<IProject>,
      required: true,
    },
  },

  beforeMount() {
    this.project = { ...this.editProject };
    this.departmentStore.getDepartments();
  },

  methods: {
    async updateProject() {
      this.updateProjectRequest = true;

      try {
        this.project = await this.projectStore.updateProject(this.project);
        this.$emit('changed', this.project);
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.updateProjectRequest = false;
    },
  },
  setup() {
    const updateProjectRequest = ref(false);

    return {
      projectStore: useProjectStore(),
      departmentStore: useDepartmentStore(),
      project: ref({} as IProject),
      updateProjectRequest,
    };
  },
});
</script>
