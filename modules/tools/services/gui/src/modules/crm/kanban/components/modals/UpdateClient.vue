<template>
  <q-card
    class="bg-secondary"
    style="width: 400px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="updateClient"
  >
    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="client.name"
          label="Назва"
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
        @click="updateClient"
        :disable="updateClientRequest"
        :loading="updateClientRequest"
      />
    </q-card-actions>
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import { defineComponent, PropType, ref } from 'vue';
import { useClientStore, IClient } from 'src/stores/client-store';
// import { useKanbanStore } from 'src/stores/kanban-store';

export default defineComponent({
  props: {
      editClient: {
          type: Object as PropType<IClient>,
          required: true
      }
  },

  beforeMount() {
    this.client = { ...this.editClient }
  },

  methods: {
    async updateClient() {
      this.updateClientRequest = true;

      try {
        this.client = await this.clientStore.updateClient(this.client);
        // this.kanbanStore.entitiesFilter.clientId = this.client.id
        this.$emit('changed', this.client);
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.updateClientRequest = false;
    },


  },
  setup() {
    const updateClientRequest = ref(false);

    return {
      clientStore: useClientStore(),
      // kanbanStore: useKanbanStore(),
      client: ref({} as IClient),
      updateClientRequest,
    };
  },
});
</script>
