import { AxiosResponse } from 'axios';
import { defineStore } from 'pinia';
import { api } from 'src/boot/api';
import { Project, CreateProjectRequest, UpdateProjectRequest } from 'src/boot/api/crm/project';

export const useProjectStore = defineStore('crm_projects', {
  state: () => ({
    projects: [] as Project[],
  }),

  getters: {
    projectsByClienId: (state) => (clientId: string) => {
      return state.projects.filter((project) => project.clientUUID === clientId);
    },
  },

  actions: {
    async getProjects(namespace: string, clientUUID: string, departmentUUID: string) {
      const res = await api.crm.project.getAll({ namespace, clientUUID, departmentUUID  })
      this.projects = res.projects;
    },

    async createProject(data: CreateProjectRequest) {
      const res = await api.crm.project.create(data)
      this.projects.push(res.project);
    },

    async updateProject(data: UpdateProjectRequest) {
      const res = await api.crm.project.update(data)
      const project = this.projects.find(
        (project) => project.uuid === res.project.uuid
      );
      if (project) Object.assign(project, res.project);
      else this.projects.push(res.project);
      return res.project;
    },

    async deleteProject(namespace: string, uuid: string) {
        await api.crm.project.delete({ namespace, uuid });
      this.projects = this.projects.filter((project) => project.uuid !== uuid);
    },
  },
});
