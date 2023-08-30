<template>
    <q-page class="q-pa-md row">
      <div class="col-2 q-pr-md">
        <MenuComponent selected="balena"/>
      </div>
      <div class="col-10">
        <q-card class="row bg-transparent" flat >
          <ServersListComponent
            :editable="true"
            :namespace="displayableNamespace"
          />
        </q-card>
        <q-separator class="q-mt-md q-mb-md"/>
        <q-card class="row q-mt-md bg-transparent" flat>
          <DevicesListComponent
            :editable="true"
            :namespace="displayableNamespace"
          />
        </q-card>
      </div>
    </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import MenuComponent from '../../MenuComponent.vue'
import ServersListComponent from './ServersListComponent.vue';
import DevicesListComponent from './DevicesListComponent.vue';

const $route = useRoute()
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string

</script>