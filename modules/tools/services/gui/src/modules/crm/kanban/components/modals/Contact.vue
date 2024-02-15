<template>
  <q-card
    class="bg-secondary"
    style="width: 400px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="updateContact"
  >
    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="contact.name"
          label="Призвище, ім’я"
        />
        <q-input
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="contact.tel1"
          label="Телефон основний"
        />
        <q-input
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="contact.tel2"
          label="Телефон додатковий"
        />
        <q-select
          outlined
          use-input
          clearable
          option-label="name"
          option-value="id"
          v-model="client"
          :options="optionsClient"
          @filter="filterClient"
          label="Клієнт"
          class="text-body1"
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

        <q-toggle
          v-model="contact.notRelevant"
          color="positive"
          size="54px"
          icon="mdi-timer-outline"
          label="Контакт не актуальний"
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
        @click="updateContact"
        :disable="updateContactRequest"
        :loading="updateContactRequest"
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
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import { defineComponent, PropType, ref } from 'vue';
import { useContactStore, IContact } from 'src/stores/contact-store';
import { IClient, useClientStore } from 'src/stores/client-store';
import Client from 'src/components/modals/UpdateClient.vue';

export default defineComponent({
  components: { Client },

  props: {
    editContact: {
      type: Object as PropType<IContact>,
      required: true,
    },
  },

  beforeMount() {
    this.contact = { ...this.editContact };
  },

  computed: {
    client: {
      get() {
        return this.clientStore.clientFiltered({ id: this.contact.clientId })
      },
      async set(value: IClient | null) {
        this.contact.clientId = value?.id as string
        if (value !== null) await this.clientStore.getClientById(this.contact.clientId)
      },
    },
  },

  methods: {
    async updateContact() {
      this.updateContactRequest = true;

      try {
        this.contact = await this.contactStore.updateContact(this.contact);
        this.$emit('changed', this.contact);
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.updateContactRequest = false;
    },

    filterClient(val: string, update: (arg0: { (): void; (): void }) => void) {
      if (val === '') {
        update(() => {
          this.optionsClient = this.clientStore.clients;
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
  },
  setup() {
    const updateContactRequest = ref(false);

    return {
      contactStore: useContactStore(),
      clientStore: useClientStore(),
      contact: ref({} as IContact),
      updateContactRequest,
      optionsClient: ref([] as { id: string; name: string }[]),
      emptyOptionsClient: ref(false),
      editClient: ref({} as IClient),
    };
  },
});
</script>
