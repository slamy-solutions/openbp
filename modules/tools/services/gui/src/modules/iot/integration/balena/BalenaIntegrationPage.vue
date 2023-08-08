<template>
    <q-page class="q-pa-md row">
      <div class="col-2 q-pr-md">
        <MenuComponent selected="balena"/>
      </div>
      <div class="col-10">
        <q-card>
          <ServersListComponent
            :editable="true"
            :namespace="displayableNamespace"
          />
        </q-card>
      </div>
    </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import MenuComponent from '../../MenuComponent.vue'
import ServersListComponent from './ServersListComponent.vue';

const displayableNamespace = ref('')

</script>