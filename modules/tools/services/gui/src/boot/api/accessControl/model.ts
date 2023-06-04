export interface ServiceManaged {
    type: 'service'
    reason: string
	service: string
	managementId: string
}

export interface IdentityManaged {
    type: 'identity'
    identityNamespace: string
	identityUUID: string
}

export interface BuiltInManaged {
    type: 'builtIn'
}

export interface NotManaged {
    type: 'none'
}
