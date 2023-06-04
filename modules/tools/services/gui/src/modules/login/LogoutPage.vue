<template>
    
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useLoginStore } from "../../stores/login-store";

onMounted(async () => {
    const loginStore = useLoginStore()
    const $router = useRouter()

    loginStore.logout()

    await $router.push({
        name: 'login'
    })
})
</script>