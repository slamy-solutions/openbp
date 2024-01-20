<template>
    <q-page padding class="row">
        <div class="col-2">
            <MenuComponent :selected="departmentUUID"/>
        </div>
        <div class="col-10">
            <!-- <div class="text-h6">Kanban settings</div> -->
            <PerformerListComponent
                :namespace="displayableNamespace"
                :department-u-u-i-d="departmentUUID"
            />

            <q-separator class="q-mt-md q-mb-md"/>

            <StagesListComponent
                :namespace="displayableNamespace"
                :department-u-u-i-d="departmentUUID"
            />
        
        </div>
    </q-page>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import MenuComponent from '../MenuComponent.vue'

import PerformerListComponent from '../performer/PerformersListComponent.vue'
import StagesListComponent from '../kanban/StagesListComponent .vue';

const $route = useRoute()
const departmentUUID = computed(() => $route.params.departmentUUID as string)
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string


</script>