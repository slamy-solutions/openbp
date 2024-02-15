<template>
  <q-card
    v-if="client"
    class="q-pt-xl bg-secondary"
    style="width: 1600px; max-width: 100%"
  >
    <q-card-actions
      v-if="$q.platform.is.desktop"
      style="width: 1600px; transform: translate(0, -50px); z-index: 1"
      class="fixed row q-pa-sm text-white bg-light-blue-10 justify-between"
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
          label="Редагувати"
          :disable="closeTaskRequest"
          @click="updateTaskModal = !updateTaskModal"
        />
        <q-btn
          icon="mdi-delete-circle-outline"
          style="font-size: 16px"
          :label="$t('main.task.cancelBtn')"
          :disable="closeTaskRequest"
          @click="cancelTask"
        />
      </div>
    </q-card-actions>

    <q-card-actions
      v-if="$q.platform.is.mobile"
      style="transform: translate(0, -50px); z-index: 1; width: 100vw"
      class="fixed row items-end text-white bg-light-blue-10"
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
        :label="$t('main.task.cancelBtn')"
        :disable="closeTaskRequest"
        @click="cancelTask"
      />
    </q-card-actions>

    <!-- назва клієнта, умови співпраці -->
    <div class="row items-start q-gutter-x-sm q-mt-xs">
      <!-- client, contact, project -->
      <div class="col-4">
        <q-card-section>
          <div class="text-h6 text-primary row" style="font-weight: 600">
            <div class="q-mr-sm">
              <q-icon size="26px" name="mdi-account-cowboy-hat" />
            </div>

            <div class="col q-gutter-x-sm">
              <div v-text="client.name" />
            </div>
          </div>
        </q-card-section>
        <!-- Контакти -->
        <q-list>
          <q-item v-for="contact in contactStore.contacts" :key="contact.id">
            <q-item-section>
              <div class="row q-gutter-x-sm items-center">
                <q-icon size="26px" name="mdi-account" />
                <div
                  v-text="contact.name"
                  class="cursor-pointer"
                  clickable
                  @click="updateContact(contact)"
                />
              </div>
            </q-item-section>

            <q-item-section avatar>
              <div
                class="row justify-between items-center"
                v-if="contact.tel2.length"
              >
                <div>
                  <q-btn-dropdown
                    v-if="$q.platform.is.desktop"
                    size="18px"
                    style="width: 54px"
                    flat
                    color="info"
                    dropdown-icon="mdi-cellphone"
                  >
                    <q-list>
                      <q-item
                        v-for="tel in [contact.tel1, contact.tel2]"
                        :key="tel"
                        clickable
                        v-close-popup
                        @click="callBackMobile(tel, contact.id)"
                      >
                        <q-item-section>
                          <q-item-label>{{ tel }}</q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-btn-dropdown>
                  <q-btn-dropdown
                    size="18px"
                    style="width: 54px"
                    flat
                    color="info"
                    dropdown-icon="phone_callback"
                  >
                    <q-list>
                      <q-item
                        v-for="tel in [contact.tel1, contact.tel2]"
                        :key="tel"
                        clickable
                        v-close-popup
                        @click="callBack(tel, contact.id)"
                      >
                        <q-item-section>
                          <q-item-label>{{ tel }}</q-item-label>
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </q-btn-dropdown>
                </div>
              </div>
              <div
                v-else-if="contact.tel1.length"
                class="row justify-between items-center"
              >
                <q-btn
                  v-if="$q.platform.is.desktop"
                  size="18px"
                  flat
                  round
                  color="info"
                  icon="mdi-cellphone"
                  @click="callBackMobile(contact.tel1, contact.id)"
                />
                <q-btn
                  size="18px"
                  flat
                  round
                  color="info"
                  icon="phone_callback"
                  @click="callBack(contact.tel1, contact.id)"
                />
              </div>
            </q-item-section>
          </q-item>
        </q-list>
        <!-- Проекти -->
        <q-list>
          <q-item
            v-for="project in projectStore.projectsByClienId(client.id)"
            :key="project.id"
            class="text-weight row justify-between items-center"
            @click="updateProject(project)"
          >
            <q-item-section>
              <div class="row q-gutter-x-sm items-center">
                <q-icon size="26px" name="mdi-briefcase-account-outline" />
                <div v-text="project.name" />
              </div>
            </q-item-section>

            <q-item-section avatar>
              <div v-text="project.department.name" />
            </q-item-section>
          </q-item>
        </q-list>
      </div>

      <!-- comments -->
      <div :class="$q.platform.is.desktop && 'col'">
        <q-card-section>
          <div class="row items-center">
            <div class="col">
              <Comments
                type="clientModal"
                :taskId="taskId"
                :comments="
                  commentStore.commentsFiltered({
                    taskId: '',
                    clientId: client.id,
                  })
                "
              />
            </div>
          </div>
        </q-card-section>
      </div>

      <!-- debit -->
      <div
        v-if="clientStore.receivables.length"
        :class="$q.platform.is.desktop && 'col-5'"
      >
        <div v-if="clientStore.receivables.length">
          <q-card-section v-if="clientStore.receivables.length">
            <q-table
              :rows="clientStore.receivables"
              :columns="columnsReceivables"
              row-key="docNumber"
            />
          </q-card-section>
        </div>

        <div v-if="clientStore.receivables.length">
          <q-card-section v-if="clientStore.receivables.length">
            <q-table
              :rows="clientStore.receivableLists"
              :columns="columnsReceivableLists"
              row-key="docNumber"
            />
          </q-card-section>
        </div>
      </div>
      <!-- ContactModal -->
      <q-dialog v-model="updateContactModal" v-if="editContact">
        <Contact
          :editContact="editContact"
          @changed="updateContactModal = false"
        />
      </q-dialog>
    </div>
  </q-card>
</template>

<script lang="ts">
import { defineComponent, ref, PropType } from 'vue';
import Comments from 'components/Comments.vue';
import { useKanbanStore } from 'src/stores/kanban-store';
import { useCommentStore } from 'src/stores/comment-store';
import { useFileStore } from 'src/stores/file-store';
import { usePerformerStore } from 'src/stores/performer-store';
import { ICall, useCallStore } from 'src/stores/call-store';
import { formatDate, formatDateWithTime } from 'src/boot/axios';
import { useProjectStore } from 'src/stores/project-store';
import { IClient, useClientStore } from 'src/stores/client-store';
import { IContact, useContactStore } from 'src/stores/contact-store';
import Contact from 'src/components/modals/Contact.vue';

export default defineComponent({
  components: { Comments, Contact },
  props: {
    client: {
      type: Object as PropType<IClient>,
      required: true,
    },
    taskId: {
      type: String as PropType<string>,
      required: true,
    },
  },

  async mounted() {
    await this.projectStore.getProjects(this.client.id, '');
    await this.contactStore.getContacts(this.client.id);
    await this.commentStore.getComments('', this.client.id);
    await this.clientStore.getReceivable(this.client.id);
    await this.clientStore.getReceivableList(this.client.id);
  },

  methods: {
    async cancelTask() {
      this.closeTaskRequest = true;

      try {
        await this.kanbanStore.deleteTask(this.client.id);
      } catch (e) {
        this.closeTaskRequest = false;
      }
      this.closeTaskRequest = false;
    },

    callBack(tel: string, contactId: string) {
      const call = {
        tel: tel,
        contactId: contactId,
        taskId: this.taskId,
        mobile: this.$q.platform.is.mobile ? true : false,
      } as ICall;
      this.callStore.callBack(call);
    },

    callBackMobile(tel: string, contactId: string) {
      const call = {
        tel: tel,
        contactId: contactId,
        taskId: this.taskId,
        mobile: true,
      } as ICall;
      this.callStore.callBack(call);
    },

    updateContact(contact: IContact) {
      this.editContact = contact;
      // if (!this.editContact.clientId) this.editContact.clientId = this.clientId;
      this.updateContactModal = true;
    },
  },

  setup() {
    return {
      kanbanStore: useKanbanStore(),
      commentStore: useCommentStore(),
      fileStore: useFileStore(),
      performerStore: usePerformerStore(),
      contactStore: useContactStore(),
      callStore: useCallStore(),
      projectStore: useProjectStore(),
      clientStore: useClientStore(),
      showFile: ref({
        modal: false,
        name: '',
        src: '',
      }),
      closeTaskRequest: ref(false),
      updateTaskModal: ref(false),
      requestChangeColumn: ref(false),
      selectPerfomerModal: ref(false),
      columnId: ref(0),
      formatDate,
      formatDateWithTime,
      columnsReceivables: [
        {
          name: 'docType',
          label: 'Документ',
          field: 'docType',
          align: 'left',
        },
        {
          name: 'docNumber',
          label: 'Номер',
          field: 'docNumber',
          align: 'center',
        },
        { name: 'docDate', label: 'Дата', field: 'docDate', align: 'center' },
        { name: 'amount', label: 'Борг', field: 'amount', align: 'right' },
        {
          name: 'overdueAmount',
          label: 'Протерм.',
          field: 'overdueAmount',
          align: 'right',
        },
        {
          name: 'docAmount',
          label: 'Сума док.',
          field: 'docAmount',
          align: 'right',
        },
        {
          name: 'daysDebt',
          label: 'Дн.',
          field: 'daysDebt',
          align: 'right',
        },
      ],
      columnsReceivableLists: [
        {
          name: 'docType',
          label: 'Документ',
          field: 'docType',
          align: 'left',
        },
        {
          name: 'docNumber',
          label: 'Номер',
          field: 'docNumber',
          align: 'center',
        },
        { name: 'docDate', label: 'Дата', field: 'docDate', align: 'center' },
        {
          name: 'amountBeginning',
          label: 'На поч.',
          field: 'amountBeginning',
          align: 'right',
        },
        {
          name: 'amountIncome',
          label: 'Прихід',
          field: 'amountIncome',
          align: 'right',
        },
        {
          name: 'amountExpense',
          label: 'Розхід',
          field: 'amountExpense',
          align: 'right',
        },
        {
          name: 'amountEnd',
          label: 'На кін.',
          field: 'amountEnd',
          align: 'right',
        },
      ],
      editContact: ref({} as IContact),
      updateContactModal: ref(false),
    };
  },
});
</script>
