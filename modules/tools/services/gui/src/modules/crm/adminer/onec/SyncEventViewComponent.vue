<template>
    <q-card
        flat
        class="bg-transparent"
    >
        <q-card-section>
            <div class="text-h6">{{ $t('modules.crm.adminer.settings.onec.sync.eventView.title') + " " + props.event.uuid }}</div>
        </q-card-section>

        <q-card-section>
            <div><span class="text-subtitle2">{{ $t('modules.crm.adminer.settings.onec.sync.eventView.timestamp') }}:</span> {{ props.event.timestamp.toLocaleString() }} </div>
            <div><span class="text-subtitle2">{{ $t('modules.crm.adminer.settings.onec.sync.eventView.success') }}:</span> {{ props.event.success }} </div>
            <div><span class="text-subtitle2">{{ $t('modules.crm.adminer.settings.onec.sync.eventView.errorMessage') }}:</span> {{ props.event.errorMessage }} </div>
        </q-card-section>

        <q-card-section style="border: solid #ccc 1px;" class="q-ma-sm">
            <VueJsonPretty :data="logData" :deep="2" :virtual="true"/>
        </q-card-section>
    </q-card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { SyncEvent } from '../../../../boot/api/crm/onec';

const props = defineProps<{
    event: SyncEvent
}>()

const logData = computed(() => {
    return props.event.log
        .trim()
        .split('\n')
        .filter((l) => l.length != 0)
        .map(line => line.trim())
        .map(line => JSON.parse(line))
})

</script>