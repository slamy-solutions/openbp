<template>
    <q-card style="width: 400px">
        <q-card-section class="text-h6"> 
            {{dayTask.name}} 
        </q-card-section>
        <q-card-section class="q-pt-none">
            <q-input autogrow outlined v-model="comment" type="text" label="Коментар" />
        </q-card-section>
        <q-card-actions align="right">
            <q-btn color="primary" label="Зупинити засікання" @click="stopTiming" />
        </q-card-actions>
    </q-card>
</template>

<script lang="ts">
import { defineComponent, PropType, ref } from 'vue'
import { useDayTaskStore } from 'src/stores/day-task-store'

export default defineComponent({
    props: {
        dayTask: {
            type: Object as PropType<{ name: string, comment: string }>,
            required: true
        },

        mode: {
            type: Function as PropType<(comment: string | null) => Promise<void>>,
            required: true
        }
    },

    mounted() {
        this.comment = this.dayTask.comment
    },

    methods: {
        async stopTiming() {
            try {
                await this.mode(this.comment)
                await this.dayTaskStore.getAllTasks()
                this.$emit('changed')
            } catch(e) {
                this.$q.notify({
                    message: (e as Error).message,
                    type: 'negative'
                })
            }
        }
    },

    setup() {
        return {
            dayTaskStore: useDayTaskStore(),
            comment: ref('')
        }
    }

})
</script>