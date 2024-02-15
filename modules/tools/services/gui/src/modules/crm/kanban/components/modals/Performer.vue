<template>
  <q-card
    class="bg-secondary"
    style="width: 400px; max-height: 100%; max-width: 100%"
    @keydown.ctrl.enter="updatePerformer"
  >
    <q-card-section class="q-pa-sm">
      <div class="q-mt-xs">
        <q-input
          autofocus
          outlined
          autogrow
          type="textarea"
          class="text-body1"
          v-model="performer.name"
          label="Призвище, ім’я"
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
        @click="updatePerformer"
        :disable="updatePerformerRequest"
        :loading="updatePerformerRequest"
      />
    </q-card-actions>
  </q-card>
</template>

<script lang="ts">
import { ApiError } from 'src/boot/axios';
import { defineComponent, PropType, ref } from 'vue';
import { usePerformerStore, IPerformer } from 'src/stores/performer-store';

export default defineComponent({
  props: {
      editPerformer: {
          type: Object as PropType<IPerformer>,
          required: true
      }
  },

  beforeMount() {
    this.performer = { ...this.editPerformer }
  },

  methods: {
    async updatePerformer() {
      this.updatePerformerRequest = true;

      try {
        await this.performerStore.updatePerformer(this.performer);
        this.$emit('changed');
      } catch (e) {
        this.$q.notify({
          type: 'negative',
          message: (e as ApiError).message,
        });
      }
      this.updatePerformerRequest = false;
    },


  },
  setup() {
    const updatePerformerRequest = ref(false);

    return {
      performerStore: usePerformerStore(),
      performer: ref({} as IPerformer),
      updatePerformerRequest,
    };
  },
});
</script>
