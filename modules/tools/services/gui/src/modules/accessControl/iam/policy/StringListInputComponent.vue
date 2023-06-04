<template>
  <q-list bordered separator>
    <q-item
        v-for="item in props.value" :key="item.id"
        class="q-pa-xs"
    >
        <q-input square filled class="full-width" dense v-model="item.value" :disable="props.disable"></q-input>
        <q-btn color="negative" icon="remove" round class="q-ml-sm" @click="remove(item.id)" v-if="!props.disable"></q-btn>
    </q-item>


    <q-item class="justify-center items-center" v-if="!props.disable">
        <q-btn dense round color="positive" icon="add" @click="addNew"/>
    </q-item>
  </q-list>
</template>

<script setup lang="ts">
    interface ValueItem {
        id: string
        value: string
    }

    const props = defineProps<{
        value: ValueItem[],
        disable?: boolean
    }>()

    function addNew() {
        props.value.push({
            id: Date.now().toString(),
            value: ""
        })
    }

    function remove(id: string) {
        props.value.splice(props.value.findIndex((s)=> s.id === id), 1)
    }
</script>

<style>

</style>