import { APIModuleBase } from '../model'

export interface ContactPerson {
    namespace: string
    uuid: string
    clientUUID: string
    name: string
    email: string
    phone: Array<string>
    notRelevant: boolean
    comment: string
}

export interface Client {
    namespace: string
    uuid: string
    name: string
    contactPersons: Array<ContactPerson>
    createdAt: Date
    updatedAt: Date
    version: number
}

export interface CreateClientRequest {
    namespace: string
    name: string
}

export interface CreateClientResponse {
    client: Client
}

export interface GetAllClientsRequest {
    namespace: string
}

export interface GetAllClientsResponse {
    clients: Array<Client>
}

export interface GetClientRequest {
    namespace: string
    uuid: string
}

export interface GetClientResponse {
    client: Client
}

export interface UpdateClientRequest {
    namespace: string
    uuid: string
    name: string
}

export interface UpdateClientResponse {
    client: Client
}

export interface DeleteClientRequest {
    namespace: string
    uuid: string
}

export interface AddContactPersonRequest {
    namespace: string
    clientUUID: string
    name: string
    email: string
    phone: Array<string>
    comment: string
}

export interface AddContactPersonResponse {
    contactPerson: ContactPerson
}

export interface GetContactPersonsForClientRequest {
    namespace: string
    clientUUID: string
}

export interface GetContactPersonsForClientResponse {
    contactPersons: Array<ContactPerson>
}

export interface UpdateContactPersonRequest {
    namespace: string
    uuid: string
    name: string
    email: string
    phone: Array<string>
    notRelevant: boolean
    comment: string
}

export interface UpdateContactPersonResponse {
    contactPerson: ContactPerson
}

export interface DeleteContactPersonRequest {
    namespace: string
    uuid: string
}

export class ClientAPI extends APIModuleBase {
    async createClient(params: CreateClientRequest): Promise<CreateClientResponse> {
        const response = await ClientAPI._axios.post<CreateClientResponse>('/crm/clients/client', params)
        return response.data
    }

    async getAllClients(params: GetAllClientsRequest): Promise<GetAllClientsResponse> {
        const response = await ClientAPI._axios.get<GetAllClientsResponse>('/crm/clients', { params })
        return response.data
    }

    async getClient(params: GetClientRequest): Promise<GetClientResponse> {
        const response = await ClientAPI._axios.get<GetClientResponse>('/crm/clients/client', { params })
        return response.data
    }

    async updateClient(params: UpdateClientRequest): Promise<UpdateClientResponse> {
        const response = await ClientAPI._axios.patch<UpdateClientResponse>('/crm/clients/client', params)
        return response.data
    }

    async deleteClient(params: DeleteClientRequest): Promise<void> {
        await ClientAPI._axios.delete('/crm/clients/client', { params })
    }

    async addContactPerson(params: AddContactPersonRequest): Promise<AddContactPersonResponse> {
        const response = await ClientAPI._axios.post<AddContactPersonResponse>('/crm/clients/client/contacts/contact', params)
        return response.data
    }

    async getContactPersonsForClient(params: GetContactPersonsForClientRequest): Promise<GetContactPersonsForClientResponse> {
        const response = await ClientAPI._axios.get<GetContactPersonsForClientResponse>('/crm/clients/client/contacts', { params })
        return response.data
    }

    async updateContactPerson(params: UpdateContactPersonRequest): Promise<UpdateContactPersonResponse> {
        const response = await ClientAPI._axios.patch<UpdateContactPersonResponse>('/crm/clients/client/contacts/contact', params)
        return response.data
    }

    async deleteContactPerson(params: DeleteContactPersonRequest): Promise<void> {
        await ClientAPI._axios.delete('/crm/clients/client/contacts/contact', { params })
    }
}