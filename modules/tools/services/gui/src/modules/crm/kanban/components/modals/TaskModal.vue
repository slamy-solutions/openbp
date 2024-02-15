<template>
  <q-card
    v-if="task"
    class="q-pt-xl bg-secondary"
    style="width: 1000px; max-width: 100%"
  >
    <q-card-actions
      v-if="$q.platform.is.desktop"
      style="width: 1000px; transform: translate(0, -50px); z-index: 1"
      :class="`fixed row q-pa-sm text-white bg-light-blue-${
        type === 'closed' ? 10 : 6 + (task.columnId > 4 ? 4 : task.columnId)
      } justify-between`"
    >
      <div class="row q-gutter-x-sm">
        <q-btn
          v-if="type !== 'closed'"
          icon="mdi-content-save-edit-outline"
          style="font-size: 16px"
          label="Редагувати"
          :disable="closeTaskRequest"
          @click="updateTaskModal = !updateTaskModal"
        />
        <q-btn
          v-if="type !== 'closed'"
          icon="mdi-check-circle-outline"
          style="font-size: 16px"
          :label="$t('main.task.closeBtn')"
          :disable="closeTaskRequest"
          @click="closeTask"
        />
        <q-btn
          v-if="type !== 'closed'"
          icon="mdi-delete-circle-outline"
          style="font-size: 16px"
          :label="$t('main.task.cancelBtn')"
          :disable="closeTaskRequest"
          @click="cancelTask"
        />
      </div>
      <q-btn
        icon="mdi-close-circle-outline"
        class="q-pa-sm"
        style="font-size: 16px"
        :label="$t('main.task.backBtn')"
        v-close-popup
      />
    </q-card-actions>

    <q-card-actions
      v-if="$q.platform.is.mobile"
      style="transform: translate(0, -50px); z-index: 1; width: 100vw"
      :class="`fixed row items-end text-white bg-light-blue-${
        type === 'closed' ? 10 : 6 + (task.columnId > 4 ? 4 : task.columnId)
      }`"
    >
      <q-btn
        class="col"
        v-if="type !== 'closed'"
        style="font-size: 16px"
        :label="$t('main.task.closeBtn')"
        :disable="closeTaskRequest"
        @click="closeTask"
      />
      <q-btn
        class="col"
        v-if="type !== 'closed'"
        style="font-size: 16px"
        :label="$t('main.task.cancelBtn')"
        :disable="closeTaskRequest"
        @click="cancelTask"
      />
      <q-btn
        class="col"
        style="font-size: 16px"
        :label="$t('main.task.backBtn')"
        v-close-popup
      />
    </q-card-actions>
    <div
      v-if="$q.platform.is.mobile && type != 'closed'"
      :class="`row items-end text-white bg-light-blue-${
        type === 'closed' ? 10 : 6 + (task.columnId > 4 ? 4 : task.columnId)
      }`"
    >
      <q-btn
        size="18px"
        :disable="requestChangeColumn"
        flat
        round
        icon="arrow_left"
        @click="changeColumnIdx(-1)"
      />
      <q-select
        outlined
        :disable="requestChangeColumn"
        map-options
        emit-value
        option-label="name"
        option-value="id"
        v-model="columnId"
        :options="kanbanStore.getColumns"
        label-color="light-blue-6"
        color="light-blue-6"
        class="col"
        style="font-size: 20px; min-width: 200px"
        @update:model-value="changeColumnId(columnId)"
      />
      <q-btn
        size="18px"
        :disable="requestChangeColumn"
        flat
        round
        icon="arrow_right"
        @click="changeColumnIdx(1)"
      />
    </div>

    <q-card-actions v-if="$q.platform.is.desktop">
      <q-btn-group outline class="fit row">
        <q-btn
          class="col"
          v-for="column in kanbanStore.getColumns"
          :key="column.id"
          :label="column.name"
          :color="`light-blue-${
            6 + (column.id > task.columnId ? -4 : column.id > 4 ? 4 : column.id)
          }`"
          @click="changeColumnId(column.id)"
        />
      </q-btn-group>
    </q-card-actions>

    <!-- назва задачі, клієнт, умови співпраці -->
    <div class="row items-start q-gutter-x-sm justify-between q-mt-xs">
      <div :class="$q.platform.is.desktop && 'col'">
        <q-card-section>
          <div class="text-h6 text-primary row" style="font-weight: 600">
            <div class="q-mr-sm">
              <q-icon size="26px" name="mdi-email-outline" />
            </div>

            <div class="col q-gutter-x-sm">
              <div v-text="task.name" />
            </div>

            <div v-if="$q.platform.is.mobile" class="q-mr-sm">
              <q-icon
                size="26px"
                name="mdi-pencil"
                @click="updateTaskModal = true"
              />
            </div>
          </div>
          <q-card-section class="col q-pa-sm q-pl-xl">
            <hr v-if="$q.platform.is.desktop" />
            <div
              :class="
                'text-subtitle1 ' +
                ($q.platform.is.desktop ? 'row justify-between' : '')
              "
            >
              <div :class="$q.platform.is.desktop && 'q-ml-sm'">
                {{
                  $t('main.task.createdDate') +
                  ': ' +
                  formatDateWithTime(task.createdDate)
                }}
              </div>

              <div>
                {{ $t('main.task.startDate') + ': ' }}
                <span
                  :class="!task.startDate && 'text-negative'"
                  v-text="
                    formatDate(task.startDate) || $t('main.task.noStartDate')
                  "
                />
              </div>

              <div class="q-ml-sm">
                {{ $t('main.task.storypoints') + ': ' }}
                <span
                  :class="task.storypoints === 0 && 'text-negative'"
                  v-text="task.storypoints"
                />
              </div>

              <div class="q-ml-sm">
                <!-- {{ $t('main.task.timing') + ': ' }} -->
                {{ 'Факт: ' }}
                <span
                  :class="task.storypoints < task.fact && 'text-negative'"
                  v-text="task.fact"
                />
              </div>
              <div class="q-ml-sm">
                {{ 'В роботі: ' }}
                <span v-text="task.leadTime" />
              </div>
            </div>
            <hr v-if="$q.platform.is.desktop" />
          </q-card-section>
        </q-card-section>

        <q-card-section>
          <div
            class="row items-center q-mr-sm text-h6 cursor-pointer"
            style="word-wrap: break-word"
            @click="clientModal = true"
          >
            <q-icon
              class="q-mr-sm"
              size="26px"
              name="mdi-card-account-details-outline"
            />
            <div class="row justify-between items-center" style="width: 90%">
              {{ task.client.name }}
              <div
                v-if="projectStore.projects.length > 1"
                class="text-subtitle2"
              >
                {{ task.project.name }}
              </div>
            </div>
          </div>
          <div
            class="text-h6 q-pl-xl justify-between"
            style="word-wrap: break-word"
          >
            <div
              class="cursor-pointer"
              clickable
              @click="updateContact(task.contact)"
              v-html="task.contact.name"
            />
            <div
              class="row justify-between items-center"
              v-if="task.contact.tel1.length"
            >
              {{ task.contact.tel1 }}
              <div>
                <q-btn
                  v-if="$q.platform.is.desktop"
                  size="18px"
                  flat
                  round
                  color="info"
                  icon="mdi-cellphone"
                  @click="callBackMobile(task.contact.tel1)"
                />
                <q-btn
                  size="18px"
                  flat
                  round
                  color="info"
                  icon="phone_callback"
                  @click="callBack(task.contact.tel1)"
                />
              </div>
            </div>
            <div
              class="row justify-between items-center"
              v-if="task.contact.tel2.length"
            >
              {{ task.contact.tel2 }}
              <div>
                <q-btn
                  v-if="$q.platform.is.desktop"
                  size="18px"
                  flat
                  round
                  color="info"
                  icon="mdi-cellphone"
                  @click="callBackMobile(task.contact.tel1)"
                />
                <q-btn
                  size="18px"
                  flat
                  round
                  color="info"
                  icon="phone_callback"
                  @click="callBack(task.contact.tel2)"
                />
              </div>
            </div>
          </div>
        </q-card-section>
        <q-card-section
          class="cursor-pointer"
          @click="selectPerfomerModal = true"
        >
          <div
            class="text-h6 row justify-between items-center"
            style="word-wrap: break-word"
          >
            <div>
              <q-icon class="q-mr-sm" size="26px" name="mdi-account-outline" />
              {{ task.performer.name }}
            </div>
            <q-avatar size="36px">
              <q-img
                :src="
                  performerStore.performerById(task.performer.id)?.avatar.url
                "
              />
            </q-avatar>
          </div>
        </q-card-section>
        <q-card-section
          v-if="task.taskTemplateId.length"
          class="cursor-pointer"
          @click="selectTaskTemplateModal = true"
        >
          <div
            class="text-h6 row justify-between items-center"
            style="word-wrap: break-word"
          >
            <div>
              <q-icon
                class="q-mr-sm"
                size="26px"
                name="mdi-information-outline"
              />
              {{ 'Створена на основі шаблону' }}
            </div>
          </div>
        </q-card-section>
        <q-card-section>
          <div
            class="text-subtitle1 row items-center"
            style="word-wrap: break-word"
          >
            <q-icon class="q-mr-md" size="26px" name="mdi-bell-cog-outline" />
            {{ 'Сповіщувати про зміну статусу' }}
            <q-checkbox
              v-model="task.isTracking"
              size="52px"
              class="q-ml-sm"
              checked-icon="mdi-bell-ring"
              unchecked-icon="mdi-bell-outline"
              indeterminate-icon="help"
              @update:model-value="updateTracking(task.isTracking)"
            />
          </div>
        </q-card-section>
      </div>

      <div :class="$q.platform.is.desktop && 'col'">
        <q-card-section>
          <div
            v-if="task.name != task.description.replace(/<[^>]+>/g, '')"
            class="row q-gutter-x-md"
          >
            <q-icon size="26px" name="mdi-text" />
            <div
              class="col q-pa-xs text-body1"
              style="word-wrap: break-word"
              v-html="task.description"
            />
          </div>
          <div v-if="task.files?.length" class="q-pl-md text-h6 text-primary">
            <q-icon size="26px" name="mdi-attachment" />
            {{ $t('main.task.files.title') }}
          </div>
          <div v-if="task.files?.length" class="q-pa-sm q-pl-xl">
            <div class="q-pa-xs full-width" style="min-height: 108.56px">
              <div
                v-for="(file, idx) in task.files"
                :key="idx"
                class="q-ma-xs inline-block"
                style="width: 100%; max-width: calc((100%) / 5 - 8px)"
              >
                <div
                  class="q-py-xs q-px-sm text-white text-weight-bold bg-primary ellipsis"
                  color="white"
                  style="border-radius: 6px 6px 0 0"
                  v-text="file.name"
                  :title="file.name"
                />
                <div
                  :class="`row justify-center items-center ${
                    canPreviewFile(file) && 'cursor-pointer'
                  }`"
                  style="height: 80px"
                  @click="
                    canPreviewFile(file) &&
                      previewFile({ fileId: file.id, fileName: file.name })
                  "
                  :title="file.name"
                >
                  <q-icon
                    :name="getIconFromMime(getMimeFromFileName(file.name))"
                    size="70px"
                  />
                </div>
                <div class="row bg-primary" style="border-radius: 0 0 6px 6px">
                  <q-btn
                    flat
                    dense
                    class="col"
                    text-color="white"
                    :disable="!canPreviewFile(file)"
                    :icon="
                      canPreviewFile(file)
                        ? 'mdi-eye-outline'
                        : 'mdi-eye-off-outline'
                    "
                    @click="
                      previewFile({ fileId: file.id, fileName: file.name })
                    "
                  />
                  <q-btn
                    flat
                    dense
                    class="col"
                    text-color="white"
                    icon="mdi-download-outline"
                    @click="
                      downloadFile({ fileId: file.id, fileName: file.name })
                    "
                  />
                </div>
              </div>
            </div>
            <q-dialog
              v-model="showFile.modal"
              :persistent="
                fileMimeTypeHasString(showFile.name, 'application') ||
                fileMimeTypeHasString(showFile.name, 'video')
              "
            >
              <ShowFileModal
                @changed="showFile.modal = false"
                :showFile="showFile"
              />
            </q-dialog>
          </div>

          <div
            v-if="
              type !== 'closed' || !commentStore.getComments(task.id)?.length
            "
            class="row items-center"
          >
            <div class="col">
              <Comments
                :type="type"
                :taskId="task.id"
                :comments="commentStore.commentsFiltered({ taskId: task.id })"
              />
            </div>
          </div>
        </q-card-section>
      </div>
    </div>

    <q-dialog v-model="updateTaskModal" v-if="task">
      <UpdateTask :editTask="task" @changed="updateTaskModal = false" />
    </q-dialog>
    <q-dialog v-model="clientModal">
      <ClientModal
        :client="task.client"
        :taskId="task.id"
        @changed="clientModal = false"
      />
    </q-dialog>
    <!-- ContactModal -->
    <q-dialog v-model="updateContactModal" v-if="editContact">
      <Contact
        :editContact="editContact"
        @changed="updateContactModal = false"
      />
    </q-dialog>

    <q-dialog v-model="selectPerfomerModal" v-if="task">
      <SelectPerformer
        :editTask="task"
        @changed="selectPerfomerModal = false"
      />
    </q-dialog>
    <q-dialog v-model="selectTaskTemplateModal" v-if="task">
      <TaskTemplate
        :editTaskTemplate="taskTemplateStore.taskTemplates[0]"
        @changed="selectTaskTemplateModal = false"
      />
    </q-dialog>
  </q-card>
</template>

<script lang="ts">
import { defineComponent, ref, PropType } from 'vue';
import downloadFile from 'src/utils/download-file';
import { AxiosError } from 'axios';
import Comments from 'components/Comments.vue';
import ShowFileModal from 'components/modals/ShowFile.vue';
import getIconFromMime from 'src/utils/mime-icon';
import {
  getMimeFromFileName,
  fileMimeTypeHasString,
} from 'src/utils/mime-type';
import { ITask, useKanbanStore } from 'src/stores/kanban-store';
import { useCommentStore } from 'src/stores/comment-store';
import { IFile, useFileStore } from 'src/stores/file-store';
import { usePerformerStore } from 'src/stores/performer-store';
import { ICall, useCallStore } from 'src/stores/call-store';
import UpdateTask from 'src/components/modals/UpdateTask.vue';
import SelectPerformer from 'src/components/modals/SelectPerformer.vue';
import TaskTemplate from 'src/components/modals/TaskTemplate.vue';
import { useTaskTemplateStore } from 'src/stores/task-template-store';
import ClientModal from 'src/components/modals/ClientModal.vue';
import { formatDate, formatDateWithTime } from 'src/boot/axios';
import { useProjectStore } from 'src/stores/project-store';
import Contact from 'src/components/modals/Contact.vue';
import { IContact } from 'src/stores/contact-store';

export default defineComponent({
  props: {
    type: {
      type: String as PropType<string>,
    },
  },

  // components: { Comments, ShowFileModal, UpdateTask },
  components: {
    Comments,
    ShowFileModal,
    UpdateTask,
    SelectPerformer,
    TaskTemplate,
    ClientModal,
    Contact,
  },

  async mounted() {
    //   void this.setCommentsSeen();
    if (!this.$route.query.id) return;
    await this.projectStore.getProjects(
      (this.task as ITask).client.id,
      (this.task as ITask).department.id
    );
    await this.taskTemplateStore.getTaskTemplatesById(
      (this.task as ITask).taskTemplateId
    );
    this.columnId = this.task?.columnId as number;
    await this.getComments();
  },

  async beforeUnmount() {
    void this.$router.push({ query: { id: undefined } });
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
    async getComments() {
      try {
        if (!this.task) return;
        await this.commentStore.getComments(this.task.id, '');
      } catch (e) {
        console.log(e);
      }
    },

    // async setCommentsSeen() {
    //   if (
    //     !this.task ||
    //     [
    //       ...this.commentStore.commentsFiltered({ taskId: this.task.id }),
    //     ].reverse()[0]?.seen !== false
    //   )
    //     return;
    //   try {
    //     await this.commentStore.setCommentSeen(this.task.id);
    //   } catch (e) {
    //     console.log(e);
    //   }
    // },

    async closeTask() {
      this.closeTaskRequest = true;

      try {
        if (!this.task) return;
        await this.kanbanStore.closeTask(this.task);
      } catch (e) {}

      this.closeTaskRequest = false;
    },

    async cancelTask() {
      this.closeTaskRequest = true;

      try {
        if (!this.task) return;
        await this.kanbanStore.deleteTask(this.task.id);
      } catch (e) {
        this.closeTaskRequest = false;
      }
      this.closeTaskRequest = false;
    },

    async loadFile({ fileId, fileName }: { fileId: string; fileName: string }) {
      try {
        if (!this.task) return;
        await this.fileStore.getFile({
          id: fileId,
          fileName,
          sourceId: this.task.id,
        });
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as AxiosError).message,
        });
      }
    },

    canPreviewFile(file: { id: string; name: string }) {
      if (
        fileMimeTypeHasString(file.name, 'application') &&
        getMimeFromFileName(file.name) !== 'application/pdf'
      )
        return false;
      if (
        fileMimeTypeHasString(file.name, 'video') &&
        getMimeFromFileName(file.name) !== 'video/mp4'
      )
        return false;
      return true;
    },

    async previewFile({
      fileId,
      fileName,
    }: {
      fileId: string;
      fileName: string;
    }) {
      if (!this.task) return;
      if (!this.fileStore.getFileById({ fileId, sourceId: this.task.id }))
        await this.loadFile({ fileId, fileName });
      this.showFile.src = (
        this.fileStore.getFileById({ fileId, sourceId: this.task.id }) as IFile
      ).url;
      this.showFile.name = (
        this.fileStore.getFileById({ fileId, sourceId: this.task.id }) as IFile
      ).file.name;
      this.showFile.modal = true;
    },

    async downloadFile({
      fileId,
      fileName,
    }: {
      fileId: string;
      fileName: string;
    }) {
      if (!this.task) return;
      if (!this.fileStore.getFileById({ fileId, sourceId: this.task.id }))
        await this.loadFile({ fileId, fileName });
      downloadFile(
        (
          this.fileStore.getFileById({
            fileId: fileId,
            sourceId: this.task.id,
          }) as IFile
        ).file
      );
    },

    changeColumnId(id: number) {
      if (!this.task) return;
      this.requestChangeColumn = true;
      const changeTask = { id: this.task.id, columnId: id };
      this.kanbanStore.updateTask(changeTask);
      this.requestChangeColumn = false;
    },

    changeColumnIdx(id: number) {
      this.columnId + id < 0 ? 0 : (this.columnId = this.columnId + id);
      this.changeColumnId(this.columnId);
    },

    callBack(tel: string) {
      if (!this.task) return;
      const call = {
        tel: tel,
        taskId: this.task.id,
        mobile: this.$q.platform.is.mobile ? true : false,
      } as ICall;
      console.log(call.mobile);
      this.callStore.callBack(call);
    },

    callBackMobile(tel: string) {
      if (!this.task) return;
      const call = {
        tel: tel,
        taskId: this.task.id,
        mobile: true,
      } as ICall;
      this.callStore.callBack(call);
    },

    updateContact(contact: IContact) {
      this.editContact = contact;
      this.editContact.clientId = this.task?.client.id as string;
      this.editContact.notRelevant = false;
      this.updateContactModal = true;
    },

    updateTracking(isTracking: boolean) {
      const changeTask = { id: this.task?.id, isTracking };
      this.kanbanStore.updateTask(changeTask);
    },
  },

  setup() {
    return {
      kanbanStore: useKanbanStore(),
      commentStore: useCommentStore(),
      fileStore: useFileStore(),
      performerStore: usePerformerStore(),
      callStore: useCallStore(),
      projectStore: useProjectStore(),
      taskTemplateStore: useTaskTemplateStore(),
      showFile: ref({
        modal: false,
        name: '',
        src: '',
      }),
      closeTaskRequest: ref(false),
      updateTaskModal: ref(false),
      clientModal: ref(false),
      requestChangeColumn: ref(false),
      selectPerfomerModal: ref(false),
      selectTaskTemplateModal: ref(false),
      columnId: ref(0),
      getMimeFromFileName,
      getIconFromMime,
      fileMimeTypeHasString,
      formatDate,
      formatDateWithTime,
      editContact: ref({} as IContact),
      updateContactModal: ref(false),
    };
  },
});
</script>
