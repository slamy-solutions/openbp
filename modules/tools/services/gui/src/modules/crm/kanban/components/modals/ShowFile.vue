<template>
    <div v-if="fileMimeTypeHasString(showFile.name, 'image')" style="max-width: 90vw; max-height: 100%">
        <img :src="showFile.src" style="width: 100%" />
    </div>
    <div v-if="fileMimeTypeHasString(showFile.name, 'application')" class="row bg-accent no-scroll" style="max-width: 90vw; max-height: 100%">
        <div class="column">
            <q-btn class="col bg-white q-px-sm no-border-radius" icon="mdi-arrow-left" @click="onChange" />
        </div>
        <iframe class="col" :src="showFile.src" :type="getMimeFromFileName(showFile.name)" style="width: 90vw; height: 90vh;" />
    </div>
    <div v-if="fileMimeTypeHasString(showFile.name, 'video')" class="row bg-accent no-scroll" style="max-width: 90vw; max-height: 100%">
        <div class="column">
            <q-btn class="col bg-white q-px-sm no-border-radius" icon="mdi-arrow-left" @click="onChange" />
        </div>
        <video controls class="col" style="min-height: inherit">
            <source :src="showFile.src" :type="getMimeFromFileName(showFile.name)">
        </video>
    </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import getIconFromMime from '../../utils/mime-icon'
import { getMimeFromFileName, fileMimeTypeHasString } from '../../utils/mime-type'

export default defineComponent({
    props: {
        showFile: {
            type: Object as PropType<{
                modal: boolean
                name: string
                src: string
            }>,
            required: true
        }
    },

    emits: ['changed'],

    methods: {
        onChange() {
            this.$emit('changed')
        }
    },

    setup() {
        return {
            getMimeFromFileName,
            getIconFromMime,
            fileMimeTypeHasString
        }
    }
});
</script>
