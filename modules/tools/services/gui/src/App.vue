<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useQuasar } from 'quasar';
import api from './boot/api';

const $q = useQuasar()

// Get information about loaded modules so the application can understand which functions are available
import { useModulesStore } from './stores/modules-store';
const modulesStore = useModulesStore()
async function getLoadedModules() {
  try {
    const response = await api.modules.getStatus()
    modulesStore.updateModulesState(response)
  } catch (error) {
    // Persist information until window reload
    $q.notify({
      message: "Failed to get information about loaded application modules. Please, refresh the page. Error: " + String(error),
      type: 'negative',
      timeout: 30000
    })
  }
}

onMounted(async () => {
  await getLoadedModules()
})

</script>
