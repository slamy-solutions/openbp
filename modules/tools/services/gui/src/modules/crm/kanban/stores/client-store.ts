import { AxiosResponse } from 'axios';
import { defineStore } from 'pinia';
import { api } from 'src/boot/api';
import { Client, UpdateClientRequest, CreateClientRequest } from 'src/boot/api/crm/client';

export interface IReceivable {
  clientId: string
  docId: string
  docType: string
  docNumber: string
  docDate: string
  amount: string
  overdueAmount: string
  docAmount: string
  daysDebt: number
  paymentDate: Date
}

export interface IReceivableList {
  clientId: string
  docDate: string
  docId: string
  docType: string
  docNumber: string
  amountBeginning: string
  amountIncome: string
  amountExpense: string
  amountEnd: string
}

export interface IClientsFilter {
  telephone?: string
  activityDateBefore?: Date
  activityDateAfter?: Date
  includesText?: string
}

export interface IClientFilter {
  id?: string
}

export const useClientStore = defineStore('crm_clients', {
  state: () => ({
    clients: [] as Client[],
    receivables: [] as IReceivable[],
    receivableLists: [] as IReceivableList[],
    entitiesFilter: {} as IClientsFilter
  }),

  getters: {
    clientsFiltered: (state) => {
      let clients = state.clients
      // console.log(clients)
      if (state.entitiesFilter.includesText) clients = clients.filter(client => client.name.toLowerCase().includes((state.entitiesFilter.includesText as string).toLowerCase()))
      if (state.entitiesFilter.activityDateBefore) clients = clients.filter(client => client.activityDate < (state.entitiesFilter.activityDateBefore as Date))
      if (state.entitiesFilter.activityDateAfter) clients = clients.filter(client => client.activityDate > (state.entitiesFilter.activityDateAfter as Date))
      if (state.entitiesFilter.telephone) clients = clients.filter(client => client.telephone === state.entitiesFilter.telephone)
      return clients
    },

    clientFiltered: (state) => (filters: IClientFilter) => {
      return state.clients.find(client => client.uuid === filters.id)
    }
  },

  actions: {
    async getClients(namespace: string) {
      const res = await api.crm.client.getAllClients({ namespace  })
      this.clients = res.clients
    },

    async getClientById(namespace: string, uuid: string) {
      const res = await api.crm.client.getClient({ namespace, uuid })
      return res.client
    },

    async getReceivable(clientId: string) {
      throw new Error('Not implemented');
      
        // const res: AxiosResponse<IReceivable[]> = await api.get(`/receivable/${useUserStore().token}`, { params: {clientId}})
      // this.receivables = res.data
    },

    async getReceivableList(clientId: string) {
        throw new Error('Not implemented');

      //const res: AxiosResponse<IReceivableList[]> = await api.get(`/receivable-list/${useUserStore().token}`, { params: {clientId}})
      //this.receivableLists = res.data
    },

    async createClient (data: CreateClientRequest) {
      const res = await api.crm.client.createClient(data)
      // useKanbanStore().entitiesFilter.clientId = res.data.id      
      this.clients.push(res.client)
    },

    async updateClient (data: UpdateClientRequest) {
        const res = await api.crm.client.updateClient(data)
      const client = this.clients.find(client => client.uuid === res.client.uuid)
      if (client) Object.assign(client, res.client)
      else this.clients.push(res.client)
      return res.client
      // useKanbanStore().entitiesFilter.clientId = res.data.id      
    },

    async deleteClient (namespace: string, uuid: string) {
        await api.crm.client.deleteClient({ namespace, uuid });
        this.clients = this.clients.filter(client => client.uuid !== uuid)
    }

  }
});
