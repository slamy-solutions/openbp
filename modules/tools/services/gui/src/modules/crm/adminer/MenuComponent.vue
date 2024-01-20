<template>
  <q-card class="full-width full-height bg-transparent" flat >
    <q-card-section>
        <div class="text-overline">General</div>
        <q-list>
            <q-separator />
            <div v-for="view in mainViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName"/>
                <q-separator />
            </div>
        </q-list>

        <div class="text-overline q-mt-sm">Departments</div>
        <q-list>
            <q-separator />
            <div v-for="view in departmentsViews" :key="view.name">
                <MenuListItem :name="view.name" :title="view.title" :icon="view.icon" :selected="view.name === props.selected" :route="view.routerPageName" :query="view.query"/>
                <q-separator />
            </div>
        </q-list>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { defineProps, ref, onMounted, defineExpose } from 'vue';
import { useRoute } from 'vue-router';
import MenuListItem from './MenuListItem.vue';
import api from 'src/boot/api';

interface Viewdata {
    name: string
    title: string
    icon: string
    routerPageName: string
    query?: Record<string, string>
}

const props = defineProps<{
    selected: string
}>()

const $route = useRoute()
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string

const mainViews = [
    {
        name: 'dashboard',
        icon: 'fa-solid fa-chart-line',
        title: 'Dashboard',
        routerPageName: 'crm_adminer_dashboard'
    },
    {
        name: 'settings',
        icon: 'settings',
        title: 'Settings',
        routerPageName: 'crm_adminer_settings'
    },
    {
        name: 'departments',
        icon: 'group',
        title: 'Departments',
        routerPageName: 'crm_adminer_departments'
    },
] as Viewdata[] 


const departmentsLoading = ref(false)
const departmentsLoadingError = ref('')
const departmentsViews = ref([
    {
        name: 'department',
        icon: 'group',
        title: 'Departments',
        routerPageName: 'crm_adminer_department',
        query: {
            departmentUUID: 'departments'
        }
    },
] as Viewdata[])


async function loadDepartments() {
    departmentsLoading.value = true
    try {
        const response = await api.crm.department.getAll({ namespace: displayableNamespace })
        departmentsViews.value = response.departments.map((d) => {
            return {
                name: d.uuid,
                icon: 'group',
                title: d.name,
                routerPageName: 'crm_adminer_department',
                query: {
                    departmentUUID: d.uuid
                }
            } as Viewdata
        })
        departmentsLoadingError.value = ''
    } catch (e) {
        departmentsLoadingError.value = String(e)
    } finally {
        departmentsLoading.value = false
    }
}
onMounted(async () => {
    await loadDepartments()
})

defineExpose({
    loadDepartments
})

</script>

<style>
.no-minimums {
    min-width: 0px;
    min-height: 0px;
}
</style>