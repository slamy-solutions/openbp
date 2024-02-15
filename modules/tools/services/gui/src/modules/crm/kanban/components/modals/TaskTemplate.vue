<template>
  <q-card
    class="q-pt-xl bg-secondary"
    style="width: 320px; max-width: 100%"
    @keydown.ctrl.enter="updateTaskTemplate"
  >
    <q-card-actions
      v-if="$q.platform.is.desktop"
      style="width: 320px; transform: translate(0, -50px); z-index: 1"
      :class="`fixed row q-pa-sm text-white bg-light-blue justify-between`"
    >
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
        @click="updateTaskTemplate"
        :disable="selectTaskTemplateModal"
        :loading="selectTaskTemplateModal"
      />
    </q-card-actions>

    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="taskTemplate.name"
          label="Назва шаблону"
        />
      </div>

      <div class="q-mt-xs">
        <q-input
          outlined
          type="number"
          class="col text-h6"
          input-style="line-height: 24px"
          v-model="taskTemplate.storyPointPlan"
          label="СП план"
        />
      </div>

      <div class="q-mt-xs">
        <q-input
          outlined
          type="number"
          class="col text-h6"
          input-style="line-height: 24px"
          v-model="taskTemplate.storyPointClient"
          label="СП клієнт"
        />
      </div>

      <div class="q-mt-xs">
        <q-input
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="taskTemplate.type"
          label="Вид робіт"
          readonly
        />
      </div>

      <div class="q-mt-xs">
        <q-field color="indigo-3" outlined label="Розклад" stack-label>
          <div class="q-mt-xs">
            <q-field bg-color="grey-2" filled label="Дні тижня">
              <div class="q-gutter-xs">
                <q-checkbox
                  size="xs"
                  v-model="taskTemplate.shedule.weekDays"
                  val="1"
                  label="Пн"
                />
                <q-checkbox
                  size="sm"
                  v-model="taskTemplate.shedule.weekDays"
                  val="2"
                  label="Вт"
                />
                <q-checkbox
                  size="sm"
                  v-model="taskTemplate.shedule.weekDays"
                  val="3"
                  label="Ср"
                />
                <q-checkbox
                  size="sm"
                  v-model="taskTemplate.shedule.weekDays"
                  val="4"
                  label="Чт"
                />
                <q-checkbox
                  size="sm"
                  v-model="taskTemplate.shedule.weekDays"
                  val="5"
                  label="Пт"
                />
                <q-checkbox
                  size="sm"
                  v-model="taskTemplate.shedule.weekDays"
                  val="6"
                  label="Сб"
                />
                <q-checkbox
                  size="sm"
                  v-model="taskTemplate.shedule.weekDays"
                  val="7"
                  label="Нд"
                />
              </div>
            </q-field>
          </div>
          <q-field bg-color="grey-2" filled>
            <div class="q-gutter-xs">
              <div class="q-mt-xs">
                <q-input
                  outlined
                  class="col text-h6"
                  input-style="line-height: 24px"
                  v-model="taskTemplate.shedule.beginDate"
                  filled
                  type="date"
                  label="Дата початку"
                />
              </div>
              <div class="q-mt-xs">
                <q-input
                  outlined
                  class="col text-h6"
                  input-style="line-height: 24px"
                  v-model="taskTemplate.shedule.endDate"
                  filled
                  type="date"
                  label="Дата закінчення"
                />
              </div>
              <div class="q-mt-xs">
                {{ taskTemplate.shedule.beginTime }}
                <q-input
                  outlined
                  class="col text-h6"
                  input-style="line-height: 24px"
                  v-model="taskTemplate.shedule.beginTime"
                  filled
                  type="time"
                  label="Час початку"
                />
              </div>
              <div class="q-mt-xs">
                {{ taskTemplate.shedule.endTime }}
                <q-input
                  outlined
                  class="col text-h6"
                  input-style="line-height: 24px"
                  v-model="taskTemplate.shedule.endTime"
                  filled
                  type="time"
                  label="Час закінчення"
                />
              </div>
            </div>
          </q-field>
        </q-field>
      </div>
    </q-card-section>
  </q-card>
</template>

<script lang="ts">
import { date } from 'quasar';
import { ApiError, formatTimeWithHoursMinutes } from 'src/boot/axios';
import { defineComponent, PropType, ref } from 'vue';
import {
  useTaskTemplateStore,
  ITaskTemplate,
} from 'src/stores/task-template-store';

export default defineComponent({
  props: {
    editTaskTemplate: {
      type: Object as PropType<ITaskTemplate>,
      required: true,
    },
  },

  beforeMount() {
    this.taskTemplate = { ...this.editTaskTemplate };
    this.weekDays = this.taskTemplate.shedule.weekDays as [];
    this.taskTemplate.shedule.beginDate = date.formatDate(
      this.taskTemplate.shedule.beginDate,
      'YYYY-MM-DD'
    );
    this.taskTemplate.shedule.endDate = date.formatDate(
      this.taskTemplate.shedule.endDate,
      'YYYY-MM-DD'
    );
    this.taskTemplate.shedule.beginTime = formatTimeWithHoursMinutes(
      this.taskTemplate.shedule.beginTime
    );
    this.taskTemplate.shedule.endTime = formatTimeWithHoursMinutes(
      this.taskTemplate.shedule.endTime
    );
  },

  methods: {
    async updateTaskTemplate() {
      this.selectTaskTemplateModal = true;

      try {
        await this.taskTemplateStore.updateTaskTemplate(this.taskTemplate);
        this.$emit('changed');
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.selectTaskTemplateModal = false;
    },
  },
  setup() {
    return {
      taskTemplateStore: useTaskTemplateStore(),
      taskTemplate: ref({} as ITaskTemplate),
      selectTaskTemplateModal: ref(false),
      weekDays: ref([]),
    };
  },
});
</script>
