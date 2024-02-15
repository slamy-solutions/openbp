<template>
  <q-card
    v-if="task"
    class="q-pt-xl bg-secondary"
    style="width: 700px; max-width: 100%"
  >
    <q-card-actions
      v-if="$q.platform.is.desktop"
      style="width: 700px; transform: translate(0, -50px); z-index: 1"
      :class="`fixed row q-pa-sm text-white bg-light-blue-${
        6 + (task.columnId > 4 ? 4 : task.columnId)
      } justify-between`"
    >
      <q-btn
        icon="mdi-arrow-left-circle-outline"
        class="q-pa-sm"
        style="font-size: 16px"
        :label="$t('main.task.backBtn')"
        v-close-popup
      />
      <div class="row q-gutter-x-sm">
        <q-btn
          icon="mdi-content-save-edit-outline"
          style="font-size: 16px"
          label="ВИБРАТИ КОМАНДУ БЕЗ ВИКОНАВЦЯ"
          @click="updateTaskDepartment(departmentId)"
        />
      </div>
    </q-card-actions>

    <q-card-actions
      v-if="$q.platform.is.mobile"
      style="transform: translate(0, -50px); z-index: 1; width: 89vw"
      :class="`fixed row items-end text-white bg-light-blue-${
        6 + (task.columnId > 4 ? 4 : task.columnId)
      }`"
    >
      <q-btn
        class="col"
        style="font-size: 16px"
        :label="$t('main.task.backBtn')"
        v-close-popup
      />
      <q-btn
        class="col"
        style="font-size: 16px"
        :disable="updateTaskRequest"
        label="КОМАНДУ"
        @click="updateTaskDepartment(departmentId)"
      />
    </q-card-actions>

    <div class="row q-gutter-x-sm q-mt-xs">
      <div v-if="$q.platform.is.mobile">
        <q-card-section>
          <q-select
            outlined
            map-options
            use-input
            emit-value
            option-label="name"
            option-value="id"
            v-model="departmentId"
            :options="departments"
            label-color="light-blue-6"
            label="Команда"
            color="light-blue-6"
            class="text-weight-bold"
            style="font-size: 20px; min-width: 200px"
          />
          <!-- <q-list>
            <q-item
              v-for="department in departments"
              :key="department.id"
              clickable
              v-ripple
              :active="departmentId === department.id"
              @click="departmentId = department.id"
              active-class="department-link"
            >
              <q-item-section>{{department.name}}</q-item-section>
            </q-item>
          </q-list> -->
        </q-card-section>
      </div>

      <div class="col">
        <q-card-section>
          <q-list>
            <q-item
              v-for="performer in performerStore.performerDepartmentById(
                departmentId
              )"
              :key="performer.id"
              clickable
              v-ripple
              class="text-h6 row items-center"
              @click="updateTaskPerformer(performer.id)"
            >
              <q-avatar size="36px" class="q-mr-md">
                <q-img :src="performer.avatar?.url" />
              </q-avatar>
              <q-item-section>{{ performer.name }}</q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </div>

      <div v-if="$q.platform.is.desktop" class="col-5">
        <q-card-section>
          <q-list>
            <q-item
              v-for="department in departments"
              :key="department.id"
              clickable
              v-ripple
              :active="departmentId === department.id"
              @click="departmentId = department.id"
              active-class="department-link"
            >
              <q-item-section avatar>
                <q-icon name="inbox" />
              </q-item-section>

              <q-item-section>{{ department.name }}</q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </div>
    </div>
  </q-card>
</template>

<style lang="sass">
.department-link
  color: white
  background: $light-blue-6
</style>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { ITask, useKanbanStore } from 'src/stores/kanban-store';
import { usePerformerStore } from 'src/stores/performer-store';
import { IDepartment, useDepartmentStore } from 'src/stores/department-store';

export default defineComponent({
  async mounted() {
    this.departmentId = this.task?.department.id as string;
    this.departments = await this.departmentStore.getDepartmentsByClientId(
      this.task?.client.id as string
    );
  },

  computed: {
    task() {
      if (!this.$route.query.id) return;
      return this.kanbanStore.getTaskById(this.$route.query.id as string) as
        | ITask
        | undefined;
    },
  },

  methods: {
    async updateTaskPerformer(performerId: string) {
      this.updateTaskRequest = true;

      try {
        if (!this.task) return;
        await this.kanbanStore.updateTask({
          id: this.task.id,
          performer: { id: performerId },
        });
        this.$emit('changed');
      } catch (e) {}

      this.updateTaskRequest = false;
    },
    async updateTaskDepartment(departmentId: string) {
      this.updateTaskRequest = true;

      try {
        if (!this.task) return;
        await this.kanbanStore.updateTask({
          id: this.task.id,
          department: { id: departmentId },
        });
        this.$emit('changed');
      } catch (e) {}

      this.updateTaskRequest = false;
    },
  },

  setup() {
    return {
      kanbanStore: useKanbanStore(),
      performerStore: usePerformerStore(),
      departmentStore: useDepartmentStore(),
      departmentId: ref(''),
      departments: ref<IDepartment[]>([]),
      updateTaskRequest: ref(false),
    };
  },
});
</script>
