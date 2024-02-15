<template>
  <q-card
    class="bg-secondary"
    style="width: 400px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="updateColumnKanban"
  >
    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="columnKanban.name"
          label="Назва"
        />
      </div>
    </q-card-section>

    <q-card-section class="q-pa-sm">
        <q-input
          autofocus
          outlined
          step="1"
          type="number"
          class="text-body1"
          v-model="columnKanban.idx"
          label="Порядок"
        />
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
        @click="updateColumnKanban"
        :disable="updateColumnKanbanRequest"
        :loading="updateColumnKanbanRequest"
      />
    </q-card-actions>
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import { defineComponent, PropType, ref } from 'vue';
import { useColumnKanbanStore, IColumnKanban } from 'src/stores/column-kanban-store';

export default defineComponent({
  props: {
      editColumnKanban: {
          type: Object as PropType<IColumnKanban>,
          required: true
      }
  },

  beforeMount() {
    this.columnKanban = { ...this.editColumnKanban }
  },

  methods: {
    async updateColumnKanban() {
      this.updateColumnKanbanRequest = true;

      try {
        await this.columnKanbanStore.updateColumnKanban(this.columnKanban);
        this.$emit('changed');
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.updateColumnKanbanRequest = false;
    },


  },
  setup() {
    const updateColumnKanbanRequest = ref(false);

    return {
      columnKanbanStore: useColumnKanbanStore(),
      columnKanban: ref({} as IColumnKanban),
      updateColumnKanbanRequest,
    };
  },
});
</script>
