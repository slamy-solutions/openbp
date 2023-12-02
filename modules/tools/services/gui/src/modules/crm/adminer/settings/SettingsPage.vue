<template>
    <q-page padding class="row">
        <div class="col-2">
            <MenuComponent selected="settings"/>
        </div>
        <div class="col-10">
            <BackendComponent
                :namespace-name="displayableNamespace"
            />
        </div>
    </q-page>
</template>

<script lang="ts" setup>

import { useRoute } from 'vue-router';
import BackendComponent from './BackendComponent.vue';
import MenuComponent from '../MenuComponent.vue';

const $route = useRoute()
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string

</script>