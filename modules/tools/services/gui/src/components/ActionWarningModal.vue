<template>
    <div class="window-height window-width row justify-center items-center">
      <div class="column">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1" style="border: red; border-style: solid; border-width: 10px;">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ props.header }}</h3>
              </q-card-section>
          <q-card-section>
            <span class="q-ma-sm text-bold">{{ props.bodyText }}</span>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="negative" size="lg" class="full-width" :label="props.actionButton" :loading="props.loading" :disabled="props.loading" @click="emit('actionClicked')" />
          </q-card-actions>
          </q-card>
      </div>
      </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  (e: 'actionClicked'): void
}>()

const props = defineProps<{
    header: string,
    bodyText: string,
    actionButton: string,
    loading: boolean
}>()
</script>