import { AxiosResponse } from 'axios';
import { defineStore } from 'pinia';
import { api } from 'src/boot/axios';

import { Ticket, CreateTicketRequest, UpdateTicketBasicInfoRequest } from 'src/boot/api/crm/kanban'

// import { formatDate, formatDateWithTime } from 'src/boot/axios';
import { 
    
 } from 'src/router';

export interface ITasksFilter {
  clientId?: string;
  performerId?: string;
  departmentId?: string;
  isDayTask?: boolean;
  isStartTiming?: boolean;
  outOfPlan?: boolean;
  createdDateBefore?: string;
  includesText?: string;
  type?: 'work' | 'plan';
}

export const useKanbanStore = defineStore('crm_kanban', {
  state: () => ({
    kanban: [] as Ticket[],
    entitiesFilter: {} as ITasksFilter,
    typesKanban: [
      { id: 'work', name: 'В роботі' },
      { id: 'plan', name: 'В планах' },
    ],
    addTaskModal: false,
    addDealModal: false,
    addTaskColumn: { id: 0, name: 'НОВІ' },
  }),

  getters: {
    getColumnById: (state) => (columnId: string) => {
      let kanban = state.kanban.filter((task) => task.ticketStageUUID === columnId);
      const filters = state.entitiesFilter;
      if (filters.clientId)
        kanban = kanban.filter((task) => task.clientUUID === filters.clientId);
      if (filters.performerId)
        kanban = kanban.filter(
          (task) => task.performerUUID === filters.performerId
        );
      if (filters.departmentId)
        kanban = kanban.filter(
          (task) => task.departmentUUID === filters.departmentId
        );
      if (filters.isDayTask)
        kanban = kanban.filter((task) => task.isDayTask === filters.isDayTask);
      if (filters.isStartTiming)
        kanban = kanban.filter(
          (task) => task.isStartTiming === filters.isStartTiming
        );
      if (filters.outOfPlan)
        kanban = kanban.filter((task) => task.fact > task.storypoints);
      if (filters.createdDateBefore)
        kanban = kanban.filter(
          (task) =>
            new Date(task.createdDate) <
            new Date(filters.createdDateBefore as string)
        );
      // console.log((filters.includesText as string).toLowerCase())
      if (filters.includesText)
        kanban = kanban.filter((task) =>
          task.name
            .toLowerCase()
            .includes((filters.includesText as string).toLowerCase())
        );

      return kanban;
      // return kanban.map(task => {
      //   return {
      //     ...task,
      //     createdDate: formatDateWithTime(task.createdDate),
      //     startDate: task.startDate === '0001-01-01T00:00:00Z' ? '' : formatDate(task.startDate),
      //   }
      // })
    },

    getTaskById: (state) => (id: string) => {
      const task = state.kanban.find((task) => task.UUID === id);
      if (!task) return undefined;
      return task;

      // return {
      //   ...task,
      //   createdDate: formatDateWithTime(task.createdDate),
      //   startDate: task.startDate === '0001-01-01T00:00:00Z' ? '' : formatDate(task.startDate)
      // }
    },

    getColumns: (state) => {
      let kanban = state.kanban;
      const filters = state.entitiesFilter;
      if (filters.clientId)
        kanban = kanban.filter((task) => task.clientUUID === filters.clientId);
      if (filters.performerId)
        kanban = kanban.filter(
          (task) => task.performerUUID === filters.performerId
        );
      if (filters.departmentId)
        kanban = kanban.filter(
          (task) => task.departmentUUID === filters.departmentId
        );
      if (filters.isDayTask)
        kanban = kanban.filter((task) => task.isDayTask === filters.isDayTask);
      if (filters.isStartTiming)
        kanban = kanban.filter(
          (task) => task.isStartTiming === filters.isStartTiming
        );
      if (filters.outOfPlan)
        kanban = kanban.filter((task) => task.fact > task.storypoints);
      if (filters.createdDateBefore)
        kanban = kanban.filter(
          (task) =>
            new Date(task.createdDate) <
            new Date(filters.createdDateBefore as string)
        );
      if (filters.includesText)
        kanban = kanban.filter(
          (task) =>
            task.name
              .toLowerCase()
              .includes((filters.includesText as string).toLowerCase()) ||
            task.description
              .toLowerCase()
              .includes((filters.includesText as string).toLowerCase())
        );
      const columns = kanban
        .map((task) => ({ id: task.columnId, name: task.columnName }))
        .sort((a, b) => a.id - b.id);
      const columnsToReturn = [] as { id: number; name: string }[];
      columns.forEach(({ id, name }) => {
        if (columnsToReturn.map((column) => column.id).includes(id)) return;
        columnsToReturn.push({ id, name });
      });
      return columnsToReturn;
    },

    getClients: (state) => {
      let kanban = state.kanban;
      const filters = state.entitiesFilter;
      if (filters.performerId)
        kanban = kanban.filter(
          (task) => task.performerUUID === filters.performerId
        );
      if (filters.departmentId)
        kanban = kanban.filter(
          (task) => task.departmentUUID === filters.departmentId
        );
      if (filters.isDayTask)
        kanban = kanban.filter((task) => task.isDayTask === filters.isDayTask);
      if (filters.isStartTiming)
        kanban = kanban.filter(
          (task) => task.isStartTiming === filters.isStartTiming
        );
      const clients = kanban.map((task) => ({
        id: task.clientUUID,
        name: task.client?.name || '~~~CLIENT DELETED~~~',
      }));
      const clientsToReturn = [] as { id: string; name: string }[];
      clients.forEach(({ id, name }) => {
        if (clientsToReturn.map((client) => client.id).includes(id)) return;
        clientsToReturn.push({ id, name });
      });
      return clientsToReturn.sort((a, b) => {
        if (a.name.toLowerCase() > b.name.toLowerCase()) return 1;
        if (a.name.toLowerCase() < b.name.toLowerCase()) return -1;
        return 0;
      });
    },
  },

  actions: {
    async getKanban(namespace: string) {
      const res = await api.crm.kanban.getTickets({ namespace, departmentUUID: this.entitiesFilter.departmentId })
      this.kanban = res.tickets
    },

    async createTask(task: CreateTicketRequest) {
      const res = await api.crm.kanban.createTicket(task)
      this.kanban.push(res.ticket)
    },

    async deleteTask(namespace: string, uuid: string) {
      await api.crm.kanban.deleteTicket({ namespace, uuid })
      this.kanban = this.kanban.filter((task) => task.UUID !== uuid);
    },

    async closeTask(namespace: string, uuid: string) {
      await api.crm.kanban.closeTicket({ namespace, uuid })
      this.kanban = this.kanban.filter((task) => task.UUID !== uuid);
    },

    async updateTask(data: UpdateTicketBasicInfoRequest) {
      const res = await api.crm.kanban.updateTicketBasicInfo(data)
      const task = this.kanban.find((task) => task.UUID === res.ticket.UUID);
      if (task) Object.assign(task, res.ticket);
      else this.kanban.push(res.ticket);
    },

    addTaskFile(data: { id: string; sourceId: string; name: string }) {
      this.kanban = this.kanban.map((task) => {
        if (task.id === data.sourceId)
          return {
            ...task,
            files: task.files.concat({ id: data.id, name: data.name }),
          };
        return task;
      });
    },
  },
});
