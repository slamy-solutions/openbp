<template>
    <q-page padding>
        <BackendComponent
            :namespace-name="displayableNamespace"
        />
    </q-page>
</template>

<script lang="ts" setup>

import { useRoute } from 'vue-router';
import BackendComponent from './BackendComponent.vue';

const $route = useRoute()
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string

</script>