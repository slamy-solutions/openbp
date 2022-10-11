export interface AuthCase {
    has: Array<{
        namespace: string
        resources: Array<string>
        actions: Array<string>
    }>
    requests: Array<{
        namespace: string
        resources: Array<string>
        actions: Array<string>
    }>
    authorized: boolean
    name: string
}

export function makeCases (namespace: string) {
    return [{
        has: [],
        requests: [],
        authorized: true,
        name: "Empty request with zero assigned policies is OK"
    },
    {
        has: [],
        requests: [{
            namespace,
            resources: ["someresource"],
            actions: ["read"]
        }],
        authorized: false,
        name: "Requesting scopes with zero assigned policies is not OK"
    },
    {
        has: [{
            namespace,
            resources: ["someresource"],
            actions: ["read"]
        }],
        requests: [{
            namespace,
            resources: ["someresource"],
            actions: ["read"]
        }],
        authorized: true,
        name: "Simple request and response with one valid policy and scope"
    },
    {
        has: [{
            namespace,
            resources: ["*"],
            actions: ["*"]
        }],
        requests: [{
            namespace,
            resources: ["someresource.asddas.afdsdf"],
            actions: ["read", "write", "idk"]
        }],
        authorized: true,
        name: "Wildcard for everything allows access to everything"
    },
    {
        has: [{
            namespace,
            resources: ["*"],
            actions: ["*"]
        }],
        requests: [{
            namespace,
            resources: ["someresource.*"],
            actions: ["*"]
        }],
        authorized: true,
        name: "Wildcard for wildcard is OK"
    },
    {
        has: [{
            namespace,
            resources: ["*"],
            actions: ["*"]
        }],
        requests: [{
            namespace: "",
            resources: ["someres"],
            actions: ["someac"]
        }],
        authorized: false,
        name: "Wildcard but different namespace is not OK"
    },
    {
        has: [{
            namespace,
            resources: ["someres*"],
            actions: ["someac*"]
        }],
        requests: [{
            namespace,
            resources: ["someres123123"],
            actions: ["someac.lalseu"]
        },{
            namespace,
            resources: ["someres123oasd"],
            actions: ["someacjabu"]
        }],
        authorized: true,
        name: "Wildcard for several scopes is OK"
    },
    {
        has: [{
            namespace,
            resources: ["someres*"],
            actions: ["someac*"]
        }],
        requests: [{
            namespace,
            resources: ["someres123123"],
            actions: ["someac.lalseu"]
        },{
            namespace,
            resources: ["somered"],
            actions: ["someacjabu"]
        }],
        authorized: false,
        name: "One out of two policies authorized is not OK"
    },
    {
        has: [{
            namespace,
            resources: ["res*"],
            actions: ["act*"]
        }],
        requests: [{
            namespace,
            resources: ["res1", "red", "res2"],
            actions: ["act1"]
        }],
        authorized: false,
        name: "No access to one of the resource is not OK"
    },
    {
        has: [{
            namespace,
            resources: ["res*"],
            actions: ["act*"]
        }],
        requests: [{
            namespace,
            resources: ["res1"],
            actions: ["act1", "acd", "act3"]
        }],
        authorized: false,
        name: "No access to one of the actions is not OK"
    }
] as Array<AuthCase> } 