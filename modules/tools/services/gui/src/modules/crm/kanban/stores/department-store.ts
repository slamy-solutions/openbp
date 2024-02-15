import { AxiosResponse } from 'axios';
import { defineStore } from 'pinia';
import { api } from 'src/boot/api';
import { Department, CreateDepartmentRequest, UpdateDepartmentRequest } from 'src/boot/api/crm/department';


export interface IBaseEntityFilter {
  id?: string;
}

export const useDepartmentStore = defineStore('crm_departments', {
  state: () => ({
    departments: [] as Department[],
  }),

  getters: {
    departmentFiltered: (state) => (filters: IBaseEntityFilter) => {
      return state.departments.find(
        (department) => department.uuid === filters.id
      );
    },
  },

  actions: {
    async getDepartments(namespace: string) {
        const res = await api.crm.department.getAll({ namespace })
        this.departments = res.departments.sort((a, b) => {
            if (a.name.toLowerCase() > b.name.toLowerCase()) return 1;
            if (a.name.toLowerCase() < b.name.toLowerCase()) return -1;
            return 0;
        });
        return this.departments
    },

    async getDepartmentsByClientId(clientId: string) {
      throw new Error('Not implemented');
      /*  const res: AxiosResponse<Department[]> = await api.get(
        `/department/${useUserStore().token}`,
        { params: { clientId } }
      );
      return res.data;*/
    },

    async createDepartment(data: CreateDepartmentRequest) {
        const res = await api.crm.department.create(data)
        this.departments.push(res.department)
    },

    async updateDepartment(data: UpdateDepartmentRequest) {
      const res = await api.crm.department.update(data)
      const department = this.departments.find(
        (department) => department.uuid === res.department.uuid
      );
      if (department) Object.assign(department, res.department);
      else this.departments.push(res.department);
    },

    async deleteDepartment(namespace: string, uuid: string) {
      await api.crm.department.delete({ namespace, uuid });
      this.departments = this.departments.filter(
        (department) => department.uuid !== uuid
      );
    },
  },
});
