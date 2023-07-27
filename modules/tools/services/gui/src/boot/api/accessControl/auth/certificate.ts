import { APIModuleBase } from "../../model"

export interface Certificate {
    namespace: string
    uuid: string
    description: string
    disabled: boolean
    created: Date
    updated: Date
    version: number
}


export interface ListForIdentityRequest {
    namespace: string
    identityUUID: string
    skip: number
    limit: number
}
export interface ListForIdentityResponse {
    certificates: Array<Certificate>
    totalCount: number
}

export interface RegisterAndGenerateRequest {
    namespace: string
    identityUUID: string
    description: string
    publicKey: string
}
export interface RegisterAndGenerateResponse {
    certificate: Certificate
    raw: string
}

export interface DisableRequest {
    namespace: string
    certificateUUID: string
}

export interface DeleteRequest {
    namespace: string
    certificateUUID: string
}

export class CertificateAPI extends APIModuleBase {
    async listForIdentity(params: ListForIdentityRequest): Promise<ListForIdentityResponse> {
        const response = await CertificateAPI._axios.get<ListForIdentityResponse>('/accessControl/iam/auth/certificate/listForIdentity', { params })
        const certificates = response.data.certificates.map((i) => {
            i.created = new Date(i.created) // we are receiving string in ISO format
            i.updated = new Date(i.updated)
            return i
        })
        
        return {certificates, totalCount: response.data.totalCount}
    }

    async registerAndGenerate(params: RegisterAndGenerateRequest): Promise<RegisterAndGenerateResponse> {
        const response = await CertificateAPI._axios.post<RegisterAndGenerateResponse>('/accessControl/iam/auth/certificate', params)
        const certificate = response.data.certificate
        certificate.created = new Date(certificate.created) // we are receiving string in ISO format
        certificate.updated = new Date(certificate.updated)

        return { certificate, raw: response.data.raw }
    }

    async disable(params: DisableRequest) {
        await CertificateAPI._axios.patch('/accessControl/iam/auth/certificate/disable', params)
    }

    async deleteCertificate(params: DeleteRequest) {
        await CertificateAPI._axios.delete('/accessControl/iam/auth/certificate', { params })
    }
}