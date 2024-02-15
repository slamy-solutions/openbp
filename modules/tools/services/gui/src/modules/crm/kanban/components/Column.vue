<template>
    <q-card flat class="transparent" @drop.prevent="dragTask" @dragenter.prevent @dragover.prevent>
      <q-item
        :class="`q-pa-sm text-white bg-blue-${6 + (columnId > 4 ? 4 : columnId)}`"
      >
        <q-item-section avatar style="min-width: 25px">
          <q-icon
            name="mdi-format-list-bulleted-square"
            class="q-pa-none q-ma-none"
          />
        </q-item-section>
  
        <q-item-section style="word-break: break-word;font-size: 20px" class="text-weight-bold">
          {{ columnName }}
        </q-item-section>
      </q-item>
  
      <div
        :class="`row justify-between items-center q-pl-sm q-pr-sm text-white bg-blue-${6 + (columnId > 4 ? 4 : columnId)}`"
      >
        <div>
          План: {{kanbanStore.getColumnById(columnId).reduce((total, task) => total + task.storypoints, 0).toFixed(1)}}
        </div>
        <q-btn
          flat
          style="width:30%; height:12px"
          size="12px"
          icon="mdi-plus-circle-outline"
          @click="addTask"
        />
        <div>
          Факт: {{kanbanStore.getColumnById(columnId).reduce((total, task) => total + task.fact, 0).toFixed(1)}}
        </div>
      </div>
  
      <div v-if="$q.platform.is.desktop">
        <!-- <Draggble
          v-if="$q.platform.is.desktop"
          v-model="tasks"
          @end="sortPriority"
        > -->
        <Task
          :class="`q-my-xs cursor-pointer border-${columnId > 4 ? 4 : columnId}`"
          v-for="task in kanbanStore.getColumnById(columnId)"
          :key="task.id"
          :draggable="true"
          @dragstart="e => addDragging(e, task)"
          @dragleave.prevent
          @dragend="removeDragging"
          :task="task"
          :move="move"
        />
        <!-- </Draggble> -->
      </div>
      <div v-else>
        <Task
          :class="`q-my-xs cursor-pointer non-selectable border-${
            columnId > 4 ? 4 : columnId
          }`"
          v-for="task in kanbanStore.getColumnById(columnId)"
          :key="task.id"
          :task="task"
        />
      </div>
    </q-card>
  </template>
  
  <style lang="scss">
  .dragging {
    color: transparent;
    background: rgba($color: #aaaaaa, $alpha: 0.4);
  }
  
  .border-0 {
    border-left: 4px solid $light-blue-6;
  }
  
  .border-1 {
    border-left: 4px solid $light-blue-7;
  }
  
  .border-2 {
    border-left: 4px solid $light-blue-8;
  }
  
  .border-3 {
    border-left: 4px solid $light-blue-9;
  }
  
  .border-4 {
    border-left: 4px solid $light-blue-10;
  }
  </style>
  
  <script lang="ts">
  import { defineComponent, ref, PropType } from 'vue';
  import Task from './Task.vue';
  import { ITask, useKanbanStore } from 'src/stores/kanban-store';
  
  export default defineComponent({
    components: { Task },
    // components: { Draggble, Task },
  
    props: {
      columnId: {
        type: Number as PropType<number>,
        required: true,
      },
      columnName: {
        type: String as PropType<string>,
        required: true,
      },
    },
  
    mounted() {
      const colorNumber = 10 - this.columnId;
      if (this.columnId < 6) {
        this.columnColor = 'light-blue-' + colorNumber.toString();
      }
    },
  
    methods: {
      dragTask(e: DragEvent) {
        const taskId = e?.dataTransfer?.getData('text')
        const changeTask = {'id': taskId, 'columnId': this.columnId}
        this.kanbanStore.updateTask(changeTask)
      },
  
      addDragging(e: DragEvent, task: ITask) {
        e.dataTransfer?.setData('text/plain', task.id)
        const t = e.target as HTMLBaseElement;
        setTimeout(() => t.firstElementChild?.classList.add('dragging'), 0);
      },
  
      removeDragging(e: Event) {
        const t = e.target as HTMLBaseElement;
        t.firstElementChild?.classList.remove('dragging');
      },
  
      async sortPriority() {
        const changedTasks = this.tasks
          .map((t, idx) => {
            if (t.priority !== idx) return { id: t.id, priority: idx };
          })
          .filter((t) => t) as ITask[];
  
        if (!changedTasks?.length) return;
  
        try {
          // await this.kanbaStore.editTask(changedTasks)
        } catch {
          return this.$q.notify({
            type: 'negative',
            message: this.$t('pages.messages.loadDataError'),
          });
        }
  
        this.tasks = this.tasks
          .map((t, idx) => {
            return { ...t, priority: idx };
          })
          .filter((t) => t);
      },
  
      async move(taskId: string, side: number) {
        this.tasks.some((t, idx) => {
          if (
            t.id === taskId &&
            (side === 1 ? idx !== 0 : idx !== this.tasks.length - 1)
          ) {
            [this.tasks[idx - side], this.tasks[idx]] = [
              this.tasks[idx],
              this.tasks[idx - side],
            ];
            return true;
          }
        });
        await this.sortPriority();
      },
      addTask() {
        this.kanbanStore.addTaskColumn = {id:this.columnId, name:this.columnName}
        this.kanbanStore.addTaskModal = true
      }
    },
  
    setup() {
      return {
        kanbanStore: useKanbanStore(),
        tasks: ref<ITask[]>([]),
        columnColor: ref('light-blue-8'),
      };
    },
  });
  </script>
  