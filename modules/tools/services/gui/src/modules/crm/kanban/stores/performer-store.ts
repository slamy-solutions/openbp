import { AxiosResponse } from 'axios';
import { defineStore } from 'pinia';
import { api } from 'src/boot/api';
import { Performer, CreatePerformerRequest } from 'src/boot/api/crm/performer';
// import { getMimeFromFileName } from 'src/utils/mime-type';


export interface IPerformersFilter {
  departmentId?: string;
}

export const usePerformerStore = defineStore('crm_performers', {
  state: () => ({
    performers: [] as Performer[],
    entitiesFilter: {} as IPerformersFilter,
  }),

  getters: {
    performersFiltered: (state) => {
      return state.performers.filter(
        (performer) =>
          performer.departmentUUID === state.entitiesFilter.departmentId
      );
    },

    performerById: (state) => (id: string) => {
      return state.performers.find((performer) => performer.uuid === id);
    },

    performerDepartmentById: (state) => (id: string) => {
      return state.performers.filter(
        (performer) => performer.departmentUUID === id
      );
    },
  },

  actions: {
    async getPerformers(namespace: string) {
        const res = await api.crm.performer.getAll({ namespace })
      this.performers = res.performers.sort((a, b) => {
        if (a.name.toLowerCase() > b.name.toLowerCase()) return 1;
        if (a.name.toLowerCase() < b.name.toLowerCase()) return -1;
        return 0;
      });
      await this.getAvatars();
    },

    async getAvatars() {
      /*this.performers.forEach(async (performer) => {
        if (performer.avatar.id.length != 0) {
          // const res: AxiosResponse<File> = await api.get('/files/' + useUserStore().token, { params: { id: performer.avatar.id }, responseType: 'blob' })
          // performer.avatar.url = URL.createObjectURL(new Blob([res.data], { type: getMimeFromFileName(performer.avatar.name) }))
          performer.avatar.url = `${process.env.API}/files/${
            useUserStore().token
          }?id=${performer.avatar.id}`;
        } else {
          performer.avatar.url = '';
        }
      });*/
    },
  },
});
