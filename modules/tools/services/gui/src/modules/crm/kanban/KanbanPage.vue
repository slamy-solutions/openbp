<template>
    <q-page class="q-pa-sm bg-secondary">
      <div
        :class="
          $q.platform.is.desktop ? 'row q-gutter-x-sm' : 'row justify-between'
        "
      >
        <q-btn
          outline
          icon="mdi-plus-circle-outline"
          color="light-blue-6"
          class="text-weight-bold"
          @click="kanbanStore.addTaskModal = true"
          :label="$q.platform.is.desktop && 'задачу'"
        />
        <q-btn
          outline
          icon="mdi-briefcase-plus-outline"
          color="light-blue-6"
          class="text-weight-bold"
          @click="kanbanStore.addDealModal = true"
          :label="$q.platform.is.desktop && 'справу'"
        />
        <q-select
          :disable="!department || requestKanban"
          outlined
          map-options
          emit-value
          option-label="name"
          option-value="id"
          v-model="kanbanStore.entitiesFilter.type"
          :options="kanbanStore.typesKanban"
          label-color="light-blue-6"
          label="Задачи"
          color="light-blue-6"
          class="text-weight-bold"
          style="font-size: 20px; min-width: 150px"
          @update:model-value="onChangeDepartment"
        />
  
        <q-btn
          outline
          v-if="$q.platform.is.mobile"
          icon="mdi-filter-variant"
          color="light-blue-6"
          class="text-weight-bold"
          @click="addFilterModal = true"
        />
  
        <q-btn
          outline
          v-if="$q.platform.is.mobile"
          icon="mdi-filter-variant-remove"
          color="light-blue-6"
          class="text-weight-bold"
          @click="deleteFilter"
        />
  
        <q-select
          v-if="$q.platform.is.desktop"
          :disable="requestKanban"
          outlined
          :clearable="kanbanStore.entitiesFilter.type == 'plan'"
          map-options
          emit-value
          option-label="name"
          option-value="id"
          v-model="department"
          :options="departmentStore.departments"
          label-color="light-blue-6"
          label="Команда"
          color="light-blue-6"
          class="text-weight-bold"
          style="font-size: 20px; min-width: 200px"
          @update:model-value="onChangeDepartment"
        />
  
        <q-select
          v-if="$q.platform.is.desktop"
          :disable="requestKanban"
          outlined
          clearable
          map-options
          emit-value
          option-label="name"
          option-value="id"
          v-model="kanbanStore.entitiesFilter.performerId"
          :options="performerStore.performersFiltered"
          label-color="light-blue-6"
          label="Виконавець"
          color="light-blue-6"
          class="text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        />
  
        <q-select
          v-if="$q.platform.is.desktop"
          :disable="requestKanban"
          outlined
          clearable
          map-options
          emit-value
          option-label="name"
          option-value="id"
          v-model="kanbanStore.entitiesFilter.clientId"
          :options="kanbanStore.getClients"
          label-color="light-blue-6"
          label="Клієнт"
          color="light-blue-6"
          class="text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        />
  
        <q-input
          v-if="$q.platform.is.desktop"
          :disable="requestKanban"
          type="date"
          outlined
          v-model="kanbanStore.entitiesFilter.createdDateBefore"
          label-color="light-blue-6"
          label="Створена до"
          color="light-blue-6"
          class="text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        />
  
        <q-input
          v-if="$q.platform.is.desktop"
          :disable="requestKanban"
          outlined
          clearable
          v-model="kanbanStore.entitiesFilter.includesText"
          label-color="light-blue-6"
          label="Містить"
          color="light-blue-6"
          class="text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        />
  
        <q-toggle
          v-if="$q.platform.is.desktop"
          v-model="kanbanStore.entitiesFilter.isStartTiming"
          color="positive"
          size="54px"
          icon="mdi-timer-outline"
          title="Активний хронометраж"
        />
  
        <q-toggle
          v-if="$q.platform.is.desktop"
          v-model="kanbanStore.entitiesFilter.isDayTask"
          color="info"
          size="54px"
          icon="mdi-calendar-outline"
          title="Заплановані в задачах на день"
        />
  
        <q-toggle
          v-if="$q.platform.is.desktop"
          v-model="kanbanStore.entitiesFilter.outOfPlan"
          color="info"
          size="54px"
          icon="mdi-calendar-alert"
          title="Факт перевищує план"
        />
      </div>
  
      <q-card-section
        :style="`${
          $q.platform.is.desktop ? 'max-height: 87vh' : 'max-height: 84vh'
        }`"
        class="q-pa-none scroll"
      >
        <div>
          <div
            v-if="$q.platform.is.desktop"
            class="row items-start q-gutter-x-sm justify-between q-mt-sm"
          >
            <Column
              class="col"
              v-for="column in kanbanStore.getColumns"
              :key="column.id"
              :columnId="column.id"
              :columnName="column.name"
            />
          </div>
          <div v-else class="q-mt-xs">
            <q-pull-to-refresh @refresh="loadData">
              <q-tab-panels
                class="transparent"
                v-model="pageId"
                animated
                infinite
                swipeable
              >
                <q-tab-panel
                  v-for="column in kanbanStore.getColumns"
                  :key="column.id"
                  :name="column.name"
                  class="no-padding"
                >
                  <q-page>
                    <Column :columnId="column.id" :columnName="column.name" />
                  </q-page>
                </q-tab-panel>
              </q-tab-panels>
            </q-pull-to-refresh>
          </div>
        </div>
      </q-card-section>
      <!-- </q-card> -->
      <q-dialog v-model="addFilterModal">
        <q-card
          class="col q-gutter-y-sm q-pa-sm"
          @changed="addFilterModal = false"
        >
          <q-select
            outlined
            map-options
            emit-value
            option-label="name"
            option-value="id"
            v-model="department"
            :options="departmentStore.departments"
            label-color="light-blue-6"
            label="Команда"
            color="light-blue-6"
            class="text-weight-bold"
            style="font-size: 20px; min-width: 200px"
            @update:model-value="onChangeDepartment"
          />
  
          <q-select
            outlined
            clearable
            map-options
            emit-value
            option-label="name"
            option-value="id"
            v-model="kanbanStore.entitiesFilter.performerId"
            :options="performerStore.performersFiltered"
            label-color="light-blue-6"
            label="Виконавець"
            color="light-blue-6"
            class="text-weight-bold"
            style="font-size: 20px; min-width: 200px"
          />
  
          <q-select
            outlined
            clearable
            use-input
            map-options
            emit-value
            option-label="name"
            option-value="id"
            v-model="kanbanStore.entitiesFilter.clientId"
            :options="optionsClient"
            @filter="filterClient"
            label-color="light-blue-6"
            label="Клієнт"
            color="light-blue-6"
            class="text-weight-bold"
            style="font-size: 20px; min-width: 200px"
          />
  
          <q-input
            :disable="requestKanban"
            outlined
            clearable
            v-model="kanbanStore.entitiesFilter.includesText"
            label-color="light-blue-6"
            label="Містить"
            color="light-blue-6"
            class="text-weight-bold"
            style="font-size: 20px; min-width: 200px"
          />
  
          <q-toggle
            v-model="kanbanStore.entitiesFilter.isStartTiming"
            color="positive"
            size="54px"
            icon="mdi-timer-outline"
          />
  
          <q-toggle
            v-model="kanbanStore.entitiesFilter.isDayTask"
            color="info"
            size="54px"
            icon="mdi-calendar-outline"
          />
        </q-card>
      </q-dialog>
  
      <q-dialog v-model="kanbanStore.addTaskModal">
        <AddTask @changed="kanbanStore.addTaskModal = false" />
      </q-dialog>
      <q-dialog v-model="kanbanStore.addDealModal">
        <AddDeal @changed="kanbanStore.addDealModal = false" />
      </q-dialog>
      <q-dialog
        :model-value="!!$route.query.id?.length"
        :maximized="$q.platform.is.mobile"
      >
        <TaskModal />
      </q-dialog>
  
      <footer
        v-if="$q.platform.is.mobile"
        class="fixed-bottom text-white bg-primary"
      >
        <q-tabs
          v-model="pageId"
          inverted
          swipeable
          class="shadow-3"
          indicator-color="white"
        >
          <q-tab
            v-for="column in kanbanStore.getColumns"
            :key="column.id"
            :name="column.name"
            :label="column.name"
            no-caps
          />
        </q-tabs>
      </footer>
    </q-page>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref } from 'vue';
  import Column from './components/Column.vue';
  import AddTask from './components/modals/AddTask.vue';
  import AddDeal from './components/modals/AddDeal.vue';
  import TaskModal from './components/modals/TaskModal.vue';
  import { ITasksFilter, useKanbanStore } from './stores/kanban-store';
  import { useDepartmentStore } from './stores/department-store';
  import { usePerformerStore } from './stores/performer-store';
  import { useClientStore } from './stores/client-store';
  
  export default defineComponent({
    components: { Column, AddTask, AddDeal, TaskModal },
  
    async beforeMount() {
        const namespace = this.$route.params.currentNamespace === "_global" ? "" : this.$route.params.currentNamespace as string

      await this.performerStore.getPerformers(namespace);
      await this.departmentStore.getDepartments(namespace);
      await this.clientStore.getClients(namespace);
      //const departmentId = this.performerStore.performerFiltered?.departmentId;
      //if (!this.kanbanStore.entitiesFilter.departmentId)
      //  this.kanbanStore.entitiesFilter.departmentId = departmentId;
      if (!this.kanbanStore.entitiesFilter.type)
        this.kanbanStore.entitiesFilter.type = 'work';
      //this.performerStore.entitiesFilter.departmentId = departmentId;
      // if (!this.kanbanStore.kanban.length) await this.loadData();
      await this.loadData();
    },
  
    computed: {
      department: {
        get() {
          return this.kanbanStore.entitiesFilter.departmentId;
        },
        set(value: string) {
          this.kanbanStore.entitiesFilter.performerId = undefined;
          this.kanbanStore.entitiesFilter.departmentId = value;
          this.performerStore.entitiesFilter.departmentId = value;
        },
      },
    },
  
    methods: {
      async loadData(done?: () => void) {
        const namespace = this.$route.params.currentNamespace === "_global" ? "" : this.$route.params.currentNamespace as string

        this.requestKanban = true;
        try {
          await this.kanbanStore.getKanban(namespace);
          this.pageId = this.kanbanStore.getColumns[0].name;
        } catch (e) {
          this.$q.notify({
            type: 'negative',
            message: this.$t('pages.messages.loadDataError'),
          });
        }
        this.requestKanban = false;
        done && done();
      },
  
      onChangeDepartment() {
        this.loadData();
      },
  
      deleteFilter() {
        let newFilter = {} as ITasksFilter;
        newFilter.departmentId = this.kanbanStore.entitiesFilter.departmentId;
        newFilter.type = this.kanbanStore.entitiesFilter.type;
        this.kanbanStore.entitiesFilter = newFilter;
      },
  
      filterClient(
        val: { id: string; name: string },
        update: (arg0: { (): void; (): void }) => void
      ) {
        console.log(val);
        if (val.name === '') {
          update(() => {
            this.optionsClient = this.kanbanStore.getClients;
          });
          return;
        }
  
        update(() => {
          const needle = val.name.toLowerCase();
          this.optionsClient = this.kanbanStore.getClients.filter(
            (v) => v.name.toLowerCase().indexOf(needle) > -1
          );
        });
      },
    },
    setup() {
      return {
        kanbanStore: useKanbanStore(),
        performerStore: usePerformerStore(),
        departmentStore: useDepartmentStore(),
        clientStore: useClientStore(),
        pageId: ref(''),
        addFilterModal: ref(false),
        requestKanban: ref(false),
        optionsClient: ref([] as { id: string; name: string }[]),
      };
    },
  });
  </script>
  