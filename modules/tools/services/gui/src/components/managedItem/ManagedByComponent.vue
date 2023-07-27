<template>
    <ServiceManaged v-if="props.managedBy.type === 'service'" :manage-data="props.managedBy" :dense="props.dense"/>
    <IdentityManaged v-if="props.managedBy.type === 'identity'" :manage-data="props.managedBy" :dense="props.dense"/>
    <BuiltInManaged v-if="props.managedBy.type === 'builtIn'" :manage-data="props.managedBy" :dense="props.dense"/>
    <NotManaged v-if="props.managedBy.type === 'none'" :dense="props.dense"/>
</template>

<script setup lang="ts">
import { ManagedBy } from './model';

import ServiceManaged from './ServiceManaged.vue'
import IdentityManaged from './IdentityManaged.vue'
import BuiltInManaged from './BuiltInManaged.vue'
import NotManaged from './NotManaged.vue'

const props = defineProps<{
    managedBy: ManagedBy,
    dense?: boolean
}>()

</script>