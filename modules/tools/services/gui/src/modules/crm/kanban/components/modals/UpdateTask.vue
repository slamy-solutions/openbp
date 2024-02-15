<template>
  <q-card
    class="bg-secondary"
    style="width: 800px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="updateTask"
  >
    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="task.name"
          :label="$t('main.task.name')"
          @paste="onPasteFile"
        />
      </div>
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="task.description"
          :label="$t('main.task.description')"
          @paste="onPasteFile"
        />
      </div>
      <div class="row">
        <q-input
          outlined
          type="number"
          class="col text-h6"
          maxlength="5"
          input-style="line-height: 24px"
          v-model="task.storypoints"
          label="План трудовитрат, год.:"
        />
        <q-input v-model="task.startDate" filled type="date" label="Почати" />
      </div>
    </q-card-section>

    <q-card-section class="q-pa-sm">
      <q-uploader
        @paste="onPasteFile"
        class="fit"
        :label="$t('main.task.files.title')"
        multiple
        hide-upload-btn
        @added="addFiles"
      >
        <template v-slot:list>
          <div
            class="row justify-center items-center full-width q-py-xs text-primary"
            v-if="!files[0] && $q.screen.width >= 768"
          ></div>

          <div v-else class="row q-col-gutter-xs">
            <div v-for="(file, idx) in files" :key="idx" class="col-2">
              <q-input
                class="text-weight-bold bg-primary"
                color="white"
                input-class="q-px-sm text-white"
                style="border-radius: 6px 6px 0 0"
                borderless
                dense
                v-model="file.name"
              />

              <div
                class="row justify-center items-center cursor-pointer"
                style="height: 80px"
              >
                <q-icon
                  :name="getIconFromMime(file.file.type)"
                  size="70px"
                  @click="previewFile(idx)"
                />
              </div>

              <div class="row bg-primary" style="border-radius: 0 0 6px 6px">
                <q-btn
                  flat
                  dense
                  class="col"
                  text-color="white"
                  icon="mdi-eye-outline"
                  @click="previewFile(idx)"
                />
                <q-btn
                  flat
                  dense
                  class="col"
                  text-color="white"
                  icon="mdi-trash-can-outline"
                  @click="removeFile(idx)"
                />
              </div>
            </div>
          </div>
        </template>
      </q-uploader>

      <q-dialog
        v-model="showFile.modal"
        :persistent="
          fileMimeTypeHasString(showFile.name, 'application') ||
          fileMimeTypeHasString(showFile.name, 'video')
        "
      >
        <ShowFileModal @changed="showFile.modal = false" :showFile="showFile" />
      </q-dialog>
      <q-select
        outlined
        map-options
        use-input
        emit-value
        option-label="name"
        option-value="id"
        v-model="client"
        :options="kanbanStore.getClients"
        label-color="light-blue-6"
        label="Клієнт"
        color="light-blue-6"
        class="text-weight-bold"
        style="font-size: 20px; min-width: 200px"
      />
      <q-select
        outlined
        clearable
        map-options
        use-input
        option-label="name"
        option-value="id"
        v-model="task.contact"
        :options="contactStore.contacts"
        label-color="light-blue-6"
        label="Контактна особа"
        color="light-blue-6"
        class="text-weight-bold"
        style="font-size: 20px; min-width: 200px"
      />
      <q-select
        outlined
        map-options
        use-input
        emit-value
        option-label="name"
        option-value="id"
        v-model="task.performer.id"
        :options="performerStore.performersFiltered"
        label-color="light-blue-6"
        label="Виконавець"
        color="light-blue-6"
        class="text-weight-bold"
        style="font-size: 20px; min-width: 200px"
      />
      <q-select
        v-if="showProject"
        outlined
        map-options
        use-input
        option-label="name"
        option-value="id"
        v-model="task.project"
        :options="projectStore.projects"
        label-color="light-blue-6"
        label="Проект"
        color="light-blue-6"
        class="text-weight-bold"
        style="font-size: 20px; min-width: 200px"
      />
    </q-card-section>

    <q-card-actions class="q-pt-none row justify-between">
      <q-btn
        icon="mdi-arrow-left-circle-outline"
        style="font-size: 16px"
        label="Відмінити"
        v-close-popup
      />
      <q-btn
        icon="mdi-content-save-outline"
        style="font-size: 16px"
        label="Зберегти"
        @click="updateTask"
        :disable="updateTaskRequest"
        :loading="updateTaskRequest"
      />
    </q-card-actions>
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import {
  fileMimeTypeHasString,
  getMimeFromFileName,
} from 'src/utils/mime-type';
import { defineComponent, PropType, ref } from 'vue';
import getIconFromMime from '../../utils/mime-icon';
import ShowFileModal from 'components/modals/ShowFile.vue';
import { useKanbanStore, ITask } from 'src/stores/kanban-store';
import { useCommentStore } from 'src/stores/comment-store';
import { useFileStore } from 'src/stores/file-store';
import { usePerformerStore } from 'src/stores/performer-store';
import { useProjectStore } from 'src/stores/project-store';

import { date } from 'quasar';
import { useContactStore } from 'src/stores/contact-store';

type FileInterface = File & { __img: { src: string } };

export default defineComponent({
  components: { ShowFileModal },

  props: {
    editTask: {
      type: Object as PropType<ITask>,
      required: true,
    },
  },

  computed: {
    client: {
      get() {
        return this.editTask.client.id;
      },
      set(value: string) {
        this.task.client.id = value;
      },
    },
  },

  watch: {
    async client(newValue: string) {
      await this.projectStore.getProjects(
        newValue,
        this.kanbanStore.entitiesFilter.departmentId as string
      );
      this.showProject = this.projectStore.projects.length > 1;
      this.task.project = this.projectStore.projects[0];
      await this.contactStore.getContacts(newValue);
      // this.task.contact = this.contactStore.contacts[0]
    },
  },

  beforeMount() {
    this.task = { ...this.editTask };
    this.task.startDate = date.formatDate(this.task.startDate, 'YYYY-MM-DD');
    this.projectStore.getProjects(this.task.client.id, this.task.department.id);
    this.contactStore.getContacts(this.task.client.id);
    this.showProject = this.projectStore.projects.length > 1;
  },

  methods: {
    onPasteFile(event: ClipboardEvent) {
      const items = event.clipboardData?.items;
      if (!items) return;
      const files = [];
      for (let i = 0; i < items.length; i++) {
        const item = items[i].getAsFile() as FileInterface | null;
        if (item == null) continue;
        files.push(item);
      }
      this.addFiles(files);
    },

    addFiles(files: FileInterface[]) {
      let alreadyExistErrorCount = 0;
      let maxSizeErrorCount = 0;
      const uniqueFiles = files
        .filter((file) => {
          return (
            !this.files.some((storeFile) => {
              if (
                storeFile.file.lastModified === file.lastModified &&
                storeFile.file.name === file.name &&
                storeFile.file.size === file.size
              ) {
                alreadyExistErrorCount++;
                return true;
              }
            }) ||
            !files.some((storeFile) => {
              if (
                storeFile.lastModified === file.lastModified &&
                storeFile.name === file.name &&
                storeFile.size === file.size
              ) {
                alreadyExistErrorCount++;
                return true;
              }
            })
          );
        })
        .filter((file) => {
          if (file.size > 20971520) {
            maxSizeErrorCount++;
            return false;
          }
          return true;
        })
        .map((file) => {
          return {
            name: file.name,
            file,
          };
        });

      if (alreadyExistErrorCount > 0)
        this.onAlreadyExistError(alreadyExistErrorCount);
      if (maxSizeErrorCount > 0) this.onMaxSizeError(maxSizeErrorCount);

      this.files = this.files.concat(uniqueFiles);
    },

    removeFile(idx: number) {
      this.files.splice(idx, 1);
    },

    onAlreadyExistError(count: number) {
      this.$q.notify({
        type: 'negative',
        message: this.$t('main.task.files.messages.alreadyExistError', {
          count,
        }),
      });
    },

    onMaxSizeError(count: number) {
      this.$q.notify({
        type: 'negative',
        message: this.$t('main.task.files.messages.maxSizeError', { count }),
      });
    },

    async updateTask() {
      this.updateTaskRequest = true;

      try {
        await this.kanbanStore.updateTask(this.task);
        await Promise.all(
          this.files.map((file) =>
            this.fileStore.uploadFile({
              file: file.file,
              type: 'task',
              sourceId: this.kanbanStore.$state.kanban.filter(
                (task) => task.columnId === 0
              )[
                this.kanbanStore.$state.kanban.filter(
                  (task) => task.columnId === 0
                ).length - 1
              ].id,
              filename: file.name,
            })
          )
        );
        this.$emit('changed');
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.updateTaskRequest = false;
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

    async previewFile(index: number) {
      this.showFile.src = URL.createObjectURL(
        new Blob([await this.files[index].file.arrayBuffer()], {
          type: this.files[index].file.type,
        })
      );

      this.showFile.name = this.files[index].name;
      this.showFile.modal = true;
    },
  },
  setup() {
    const showFile = ref({
      modal: false,
      name: '',
      src: '',
    });
    const files = ref<{ name: string; file: FileInterface }[]>([]);
    const updateTaskRequest = ref(false);

    return {
      kanbanStore: useKanbanStore(),
      commentStore: useCommentStore(),
      performerStore: usePerformerStore(),
      contactStore: useContactStore(),
      fileStore: useFileStore(),
      projectStore: useProjectStore(),
      task: ref({} as ITask),
      showFile,
      files,
      updateTaskRequest,
      getMimeFromFileName,
      getIconFromMime,
      fileMimeTypeHasString,
      showProject: ref(false),
    };
  },
});
</script>
