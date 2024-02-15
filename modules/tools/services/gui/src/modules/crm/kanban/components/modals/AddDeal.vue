<template>
  <q-card
    class="bg-secondary"
    style="width: 800px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="createTask"
  >
    <q-card-section class="q-pa-sm">
      <q-select
        autofocus
        outlined
        clearable
        map-options
        use-input
        option-label="idx"
        option-value="id"
        v-model="task.contact"
        :options="optionsContact"
        @filter="filterContactIdx"
        label-color="light-blue-6"
        label="Пошук"
        color="light-blue-6"
        class="text-weight-bold"
        style="font-size: 20px; min-width: 200px"
      >
        <template v-slot:append>
          <q-btn
            v-if="emptyOptionsContact"
            round
            dense
            flat
            icon="add"
            @click="createContact"
          />
        </template>
      </q-select>
      <div class="row justify-between">
        <q-select
          outlined
          clearable
          map-options
          use-input
          option-label="name"
          option-value="id"
          v-model="task.contact"
          :options="optionsContact"
          @filter="filterContact"
          label-color="light-blue-6"
          label="Контактна особа"
          color="light-blue-6"
          class="col-4 text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        >
          <template v-slot:append>
            <q-btn
              v-if="emptyOptionsContact"
              round
              dense
              flat
              icon="add"
              @click="createContact"
            />
          </template>
        </q-select>
        <q-select
          outlined
          clearable
          map-options
          use-input
          option-label="name"
          option-value="id"
          v-model="client"
          :options="optionsClient"
          @filter="filterClient"
          label-color="light-blue-6"
          label="Клієнт"
          color="light-blue-6"
          class="col text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        >
          <template v-slot:append>
            <q-btn
              v-if="emptyOptionsClient"
              round
              dense
              flat
              icon="add"
              @click="updateClientModal = true"
            />
          </template>
        </q-select>
        <q-select
          v-if="showProject"
          outlined
          map-options
          use-input
          option-label="name"
          option-value="id"
          v-model="task.project"
          :options="optionsProject"
          @filter="filterProject"
          label-color="light-blue-6"
          label="Проект"
          color="light-blue-6"
          class="col-4 text-weight-bold"
          style="font-size: 20px; min-width: 200px"
        >
          <template v-slot:append>
            <q-btn
              v-if="emptyOptionsProject"
              round
              dense
              flat
              icon="add"
              @click="createProject"
            />
          </template>
        </q-select>
      </div>
      <div class="q-mt-xs">
        <q-input
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="task.description"
          :label="$t('main.task.description')"
          @paste="onPasteFile"
        />
      </div>
      <q-select
        outlined
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
        <q-input v-model="startDate" filled type="date" label="Почати" />
      </div>
      <q-select
        outlined
        map-options
        emit-value
        option-label="name"
        option-value="id"
        v-model="task.columnId"
        :options="kanbanStore.getColumns"
        label-color="light-blue-6"
        label="Статус"
        color="light-blue-6"
        class="text-weight-bold"
        style="font-size: 20px; min-width: 200px"
      />

      <!-- files -->
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
        @click="createTask"
        :disable="createTaskRequest"
        :loading="createTaskRequest"
      />
    </q-card-actions>
    <q-dialog v-model="updateClientModal" v-if="editClient.name">
      <Client
        :editClient="editClient"
        @changed="
          (client) => {
            updateClientModal = false;
            this.client = client;
          }
        "
      />
    </q-dialog>
    <q-dialog v-model="updateContactModal" v-if="editContact">
      <Contact
        :editContact="editContact"
        @changed="
          (contact) => {
            updateContactModal = false;
            this.task.contact = contact;
          }
        "
      />
    </q-dialog>
    <q-dialog v-model="updateProjectModal" v-if="editProject">
      <Project
        :editProject="editProject"
        @changed="
          (project) => {
            updateProjectModal = false;
            this.task.project = project;
          }
        "
      />
    </q-dialog>
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import { useUserStore } from 'src/stores/user-store';
import {
  fileMimeTypeHasString,
  getMimeFromFileName,
} from 'src/utils/mime-type';
import { defineComponent, ref } from 'vue';
import getIconFromMime from '../../utils/mime-icon';
import ShowFileModal from './ShowFile.vue';
import { useKanbanStore, ITask } from '../../stores/kanban-store';
import { useCommentStore } from '../../stores/comment-store';
import { useFileStore } from '../../stores/file-store';
import { usePerformerStore } from '../../stores/performer-store';
import { date } from 'quasar';
import { IProject, useProjectStore } from '../../stores/projects-store';
import { IContact, useContactStore } from '../../stores/contact-store';
import { IClient, useClientStore } from '../../stores/client-store';
import Client from './UpdateClient.vue';
import Contact from './Contact.vue';
import Project from './Project.vue';

type FileInterface = File & { __img: { src: string } };

export default defineComponent({
  components: { ShowFileModal, Client, Contact, Project },

  computed: {
    client: {
      get() {
        return this.task.client;
      },
      set(value: IClient) {
        this.task.client = value;
      },
    },
    contact: {
      get() {
        return this.task.contact;
      },
      set(value: IContact) {
        this.task.contact = value;
      },
    },
  },

  watch: {
    async client(newValue: IClient) {
      if (newValue && newValue.id != '') {
        await this.projectStore.getProjects(
          newValue.id,
          this.kanbanStore.entitiesFilter.departmentId as string
        );
        this.showProject = this.projectStore.projects.length !== 1;
        this.task.project = this.projectStore.projects[0];
        await this.contactStore.getContacts(newValue.id);
        if (
          this.task.contact === null ||
          newValue.id != (this.task.contact as IContact).clientId
        ) {
          let contactDefault = {} as IContact;
          if (this.projectStore.projects.length) {
            contactDefault = this.contactStore.contactById(
              this.projectStore.projects[0].contactId
            ) as IContact;
          }
          if (contactDefault) {
            this.task.contact = contactDefault as IContact;
          } else {
            this.task.contact = { id: '', name: '' };
          }
        }
      } else {
        this.showProject = false;
        this.task.project = { id: '', name: '' };
        this.task.contact = null;
        await this.contactStore.getContacts('');
      }
    },
    async contact(newValue: IContact) {
      if (newValue && newValue.id != '') {
        this.task.client = await this.clientStore.getClientById(
          newValue.clientId
        );
      } else {
        if (this.task.client === null || this.task.client.id == '') {
          await this.contactStore.getContacts('');
        } else {
          await this.contactStore.getContacts(this.task.client.id);
        }
      }
    },
  },

  async mounted() {
    if (this.kanbanStore.entitiesFilter.clientId) {
      if (this.kanbanStore.entitiesFilter.clientId) {
        this.task.client = this.kanbanStore.getClients.filter(
          (client) => client.id === this.kanbanStore.entitiesFilter.clientId
        )[0];
        await this.contactStore.getContacts(
          this.kanbanStore.entitiesFilter.clientId
        );
      }
      await this.projectStore.getProjects(
        this.task.client.id,
        this.kanbanStore.entitiesFilter.departmentId as string
      );
      this.task.project = this.projectStore.projects[0];
      this.showProject = this.projectStore.projects.length !== 1;
    } else {
      await this.contactStore.getContacts('');
    }
    if (!this.kanbanStore.entitiesFilter.performerId) {
      this.kanbanStore.entitiesFilter.performerId = useUserStore().token
    }
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

    async createTask() {
      this.createTaskRequest = true;
      this.task.performer.id = this.kanbanStore.entitiesFilter
        .performerId as string;
      this.task.startDate = date.formatDate(this.startDate, 'YYYY-MM-DD');

      try {
        await this.kanbanStore.createTask(this.task as ITask);
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

        this.$q.notify({
          type: 'positive',
          message: this.$t('main.task.messages.success'),
        });
        this.$emit('changed');
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.createTaskRequest = false;
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

    filterClient(val: string, update: (arg0: { (): void; (): void }) => void) {
      if (val === '') {
        update(() => {
          this.optionsClient = this.kanbanStore.getClients;
          this.emptyOptionsClient = false;
        });
        return;
      }

      update(() => {
        const needle = val.toLowerCase();
        this.optionsClient = this.clientStore.clients.filter(
          (v) => v.name.toLowerCase().indexOf(needle) > -1
        );
        this.emptyOptionsClient = this.optionsClient.length == 0;
        this.editClient.name = val;
      });
    },

    filterContact(val: string, update: (arg0: { (): void; (): void }) => void) {
      if (val === '') {
        update(() => {
          this.optionsContact = this.contactStore.contacts;
          this.emptyOptionsContact = false;
        });
        return;
      }

      update(() => {
        const needle = val.toLowerCase();
        this.optionsContact = this.contactStore.contacts.filter(
          (v) => v.name.toLowerCase().indexOf(needle) > -1
        );
        this.emptyOptionsContact = this.optionsContact.length == 0;
        this.editContact.name = val;
      });
    },

    filterContactIdx(
      val: string,
      update: (arg0: { (): void; (): void }) => void
    ) {
      if (val === '') {
        update(() => {
          this.optionsContact = this.contactStore.contacts;
          this.emptyOptionsContact = false;
        });
        return;
      }

      update(() => {
        const needle = val.toLowerCase();
        this.optionsContact = this.contactStore.contacts.filter(
          (v) => v.idx.toLowerCase().indexOf(needle) > -1
        );
        this.emptyOptionsContact = this.optionsContact.length == 0;
        this.editContact.tel1 = val;
      });
    },

    filterProject(val: string, update: (arg0: { (): void; (): void }) => void) {
      if (val === '') {
        update(() => {
          this.optionsProject = this.projectStore.projects;
          this.emptyOptionsProject = false;
        });
        return;
      }

      update(() => {
        const needle = val.toLowerCase();
        this.optionsProject = this.projectStore.projects.filter(
          (v) => v.name.toLowerCase().indexOf(needle) > -1
        );
        this.emptyOptionsProject = this.optionsProject.length == 0;
        this.editProject.name = val;
      });
    },

    createClient() {
      this.task.client.id = '';
      this.updateClientModal = true;
    },

    createContact() {
      this.editContact.clientId = this.client ? this.client.id : '';
      this.editContact.notRelevant = false;
      this.updateContactModal = true;
    },

    createProject() {
      this.editProject.clientId = this.client.id;
      this.editProject.department = {
        id: this.kanbanStore.entitiesFilter.departmentId as string,
        name: '',
      };
      this.editProject.notRelevant = false;
      this.updateProjectModal = true;
    },
  },
  setup() {
    const showFile = ref({
      modal: false,
      name: '',
      src: '',
    });
    const task = ref({
      id: '',
      name: '',
      description: '',
      storypoints: 0.0,
      startDate: '',
      client: {
        id: '',
        name: '',
      },
      project: {
        id: '',
        name: '',
      },
      contact: null as null | { id: string; name: string },
      performer: {
        id: '',
        name: '',
      },
      columnId: useKanbanStore().addTaskColumn.id,
    });
    const startDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'));
    const files = ref<{ name: string; file: FileInterface }[]>([]);
    const createTaskRequest = ref(false);

    return {
      kanbanStore: useKanbanStore(),
      clientStore: useClientStore(),
      contactStore: useContactStore(),
      commentStore: useCommentStore(),
      performerStore: usePerformerStore(),
      fileStore: useFileStore(),
      projectStore: useProjectStore(),
      showFile,
      task,
      startDate,
      files,
      createTaskRequest,
      getMimeFromFileName,
      getIconFromMime,
      fileMimeTypeHasString,
      showProject: ref(false),
      optionsClient: ref([] as { id: string; name: string }[]),
      emptyOptionsClient: ref(false),
      editClient: ref({} as IClient),
      updateClientModal: ref(false),
      optionsContact: ref([] as { id: string; name: string }[]),
      emptyOptionsContact: ref(false),
      editContact: ref({} as IContact),
      updateContactModal: ref(false),
      editProject: ref({} as IProject),
      updateProjectModal: ref(false),
      optionsProject: ref([] as { id: string; name: string }[]),
      emptyOptionsProject: ref(false),
    };
  },
});
</script>
