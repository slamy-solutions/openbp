<template>
    <q-card flat square bordered>
      <q-card-section
        class="q-pa-xs q-ma-none row text-weight-bold text-primary"
        style="font-size: 17px"
      >
        <div class="col">
          <div class="row q-gutter-x-xs">
            <!-- <q-badge class="q-py-none text-subtitle2" :color="`light-blue-${task.columnId + 7}`">{{ task.id }}</q-badge> -->
            <div
              style="word-break: break-word;"
              v-text="task.name"
              @click="$router.push({ query: { id: task.id } })"
            />
          </div>
          <!-- <div v-if="task.name !== task.description" style="word-wrap: break-word" class="ellipsis text-subtitle1" v-html="task.description" /> -->
          <div
            v-if="!kanbanStore.entitiesFilter.clientId"
            style="word-wrap: break-word"
            class="ellipsis text-subtitle1 client"
            v-text="task.client.name"
            @click="
              kanbanStore.entitiesFilter.clientId = !kanbanStore.entitiesFilter
                .clientId
                ? task.client.id
                : ''
            "
          />
          <div
            style="word-wrap: break-word"
            class="ellipsis text-subtitle2"
            v-text="task.contact.name"
            @click="
              kanbanStore.entitiesFilter.clientId = !kanbanStore.entitiesFilter
                .clientId
                ? task.client.id
                : ''
            "
          />
          <div
            v-if="!kanbanStore.entitiesFilter.departmentId"
            style="word-wrap: break-word"
            class="ellipsis text-subtitle1 client"
            v-text="task.department.name"
          />
          <div
            v-if="task.storypoints + task.fact"
            class="row justify-between text-subtitle2"
          >
            <div v-text="`План: ${task.storypoints}`" />
            <div
              @click="
                kanbanStore.entitiesFilter.outOfPlan = kanbanStore.entitiesFilter
                  .outOfPlan
                  ? false
                  : true
              "
            >
              {{ 'Факт:' }}
              <span
                :class="task.storypoints < task.fact && 'text-negative'"
                v-text="task.fact"
              />
            </div>
            <div v-text="`В роботі: ${task.leadTime}`" />
          </div>
          <div class="row justify-between text-subtitle2">
            <div v-text="formatDateWithTime(task.createdDate)" />
            <div
              v-text="formatDate(task.startDate)"
              :class="
                kanbanStore.entitiesFilter.type === 'plan' &&
                task.columnId === 1 &&
                new Date(task.startDate) < new Date() &&
                'text-negative'
              "
            />
          </div>
          <div class="row justify-between items-center">
            <div class="row justify-between items-center">
              <q-btn
                size="12px"
                flat
                round
                icon="arrow_left"
                @click="changeColumnIdx(-1)"
              />
              <!-- <q-btn round size="10px"> -->
              <q-avatar
                v-if="!kanbanStore.entitiesFilter.performerId"
                size="28px"
                @click="
                  kanbanStore.entitiesFilter.performerId = !kanbanStore
                    .entitiesFilter.performerId
                    ? task.performer.id
                    : ''
                "
              >
                <q-img
                  :src="
                    performerStore.performerById(this.task.performer.id)?.avatar
                      .url
                  "
                />
              </q-avatar>
              <!-- </q-btn> -->
            </div>
  
            <div class="row justify-between items-center">
              <q-btn
                size="12px"
                class="isDayTask"
                flat
                round
                icon="mdi-calendar-outline"
                @click="addDayTask"
              />
              <q-icon
                v-if="task.isDayTask"
                size="32px"
                flat
                color="info"
                round
                name="mdi-calendar-outline"
                @click="
                  kanbanStore.entitiesFilter.isDayTask =
                    !kanbanStore.entitiesFilter.isDayTask
                "
              >
                <q-badge v-if="task.numberDayTasks > 1" color="info" floating>
                  {{task.numberDayTasks}} 
                </q-badge>
              </q-icon>
              <q-icon
                v-if="task.isStartTiming"
                size="32px"
                flat
                color="positive"
                round
                name="mdi-timer-outline"
                @click="
                  kanbanStore.entitiesFilter.isStartTiming =
                    !kanbanStore.entitiesFilter.isStartTiming
                "
              />
  
              <!-- <div
                v-if="type !== 'closed' && commentStore.getUnseenCommentsCount(task.id)"
                class="column justify-end items-center"
              >
                <q-badge
                  rounded
                  class="q-py-none text-subtitle2"
                  :color="`light-blue-${6 + (task.columnId > 4 ? 4 : task.columnId)}`"
                  >{{ commentStore.getUnseenCommentsCount(task.id) }}
                </q-badge>
              </div> -->
              <q-btn
                size="12px"
                flat
                round
                icon="arrow_right"
                @click="changeColumnIdx(1)"
              />
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </template>
  
  <style lang="scss">
  .client {
    &:hover {
      color: $info;
    }
  }
  .isDayTask {
    color: white;
  }
  .isDayTask:hover {
    color: $info;
  }
  </style>
  
  <script lang="ts">
  import { defineComponent, PropType } from 'vue';
  import { Ticket } from 'src/boot/api/crm/kanban'; 
  import { useKanbanStore } from '../stores/kanban-store';
  import { usePerformerStore } from '../stores/performer-store';
  import { formatDate, formatDateWithTime } from '../common';
  //import { useDayTaskStore } from '../stores/day-task-store';
  
  export default defineComponent({
    props: {
      type: {
        type: String as PropType<string>,
      },
  
      task: {
        type: Object as PropType<Ticket>,
        required: true,
      },
    },
  
    // mounted() {
    //   void this.getComments();
    // },
  
    methods: {
      // async getComments() {
      //   try {
      //     await this.commentStore.getComments(this.task.id, '');
      //   } catch (e) {
      //     console.log(e);
      //   }
      // },
  
      changeColumnIdx(idx: number) {
        const changeTask = {
          id: this.task.UUID,
          columnId: this.task.columnId + idx,
        };
        this.kanbanStore.updateTask(changeTask);
      },
  
      async addDayTask() {
        //await this.dayTaskStore.updateDayTask({ id: '', taskId: this.task.id });
        await this.kanbanStore.updateTask({ id: this.task.UUID });
      },
    },
    setup() {
      return {
        //commentStore: useCommentStore(),
        performerStore: usePerformerStore(),
        kanbanStore: useKanbanStore(),
        //dayTaskStore: useDayTaskStore(),
        formatDate,
        formatDateWithTime,
      };
    },
  });
  </script>
  