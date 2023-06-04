export default {
  layout: {
    main: {
      modules: {
        namespace: {
          name: "Namespaces",
          description: "Manage multitenancy"
        },
        accessControl: {
          name: "Access Management",
          description: ""
        }
      }
    },
  },
  modules: {
    login: {
      header: 'Login to the OpenBP',
      hint: 'There is no public registration. In case of any problems you have to directly communicate with platform administrator.',

      usernameInput: 'Username',
      passwordInput: 'Password',

      loginButton: 'Login',

      loginOperationPendingNotify: 'Trying to login ...',
      failToLoginNotify: 'Failed to login. Check your credentials and try again. Error: {error}',
      successfullyLoggedInNotify: 'Successfully logged in.',
    },
    bootstrap: {
      header: 'System is not ready yet',
      subheader: 'OpenBP is not ready and needs additional inputs. Complete the bootstrap process to start.',

      steps: {
        status: {
          label: "Check status",

          getStatusOperationNotify: 'Tryying to get bootstrap status...',
          getStatusFailNotify: 'Failed to get bootstrap status. Error: {error}'
        },
        vault: {
          label: "Unseal the vault",

          passwordInput: "Password",

          unsealButton: "Unseal",
          unsealHint: "If you dont remember the password, the only possiblity to restore the access is to use root access to the physical HSM you are using. You have to manually login to the physical HSM. If you dont know root password to the HSM - there is no way to unseal the vault and the data is lost forever."
        },
        rootUser: {
          label: "Create root user",

          usernameInput: "Username",
          passwordInput: "Password",

          createButton: "Create",
          createHint: "This will create root user with full access to the entire system. You will not be able to create one more user, so make sure you remember the username and password.",

          blockedHeader: "Root user creation is blocked",
          blockedHint: "Root user creation was manually disabled. This is not standard behaviour. That means, if you are a system administrator at some point you disabled the ability to bootstrap root user. In this case you have to use other methods to do this. Bootstrap process using this GUI may not be completed.",
        
          createOperationNotify: 'Trying to create new root user ...',
          createSuccessNotify: 'Successfully created new root user',
          createFailNotify: 'Failed to create root user. Error: {error}'
        }
      }
    },
    namespace: {
      list: {
        table: {
          header: "Namespaces",
          nameColumn: "Name",
          fullNameColumn: "Full name",
          descriptionColumn: "Description",
          actionsColumn: "Actions",
          actionsMenu: {
            delete: "Delete"
          },
          search: "Search",
          noData: "There are no namespaces",
          failedToLoad: "Failed to load namespaces list: {error}",
          createButton: "Create",

          loadOperationNotify: 'Trying to load namespaces list ...',
          loadFailNotify: 'Failed to load namespaces list. Error: {error}'
        }
      },
      create: {
        header: "Create namespace",

        nameInput: "Name",
        fullNameInput: "Full name",
        descriptionInput: "Description",

        createButton: "Create",
        createHint: "After the namespace creation, service will raise event for entire system to inform all the services about namespace creation. You may have to wait several seconds before all the services will recognize the change.",

        createOperationNotify: 'Trying to create new namespace ...',
        createSuccessNotify: 'Successfully created new namespace',
        createFailNotify: 'Failed to create namespace. Error: {error}'
      },
      delete: {
        header: 'Delete namespace?',
        bodyText: 'You are about to delete "{namespaceName}" namespace. This action can not be undone - all the data related to the namespace will be lost. This action will not only remove the entry from the namespaces list, but delete the entire namespace database and inform all the services to clear the namespace information. This process is asynchronous, so you will have to wait for all the services will recognize the change (most probably up to several seconds).',
      
        deleteButton: 'Delete',

        deleteOperationNotify: 'Trying to delete namespace ...',
        deleteSuccessNotify: 'Successfully deleted namespace',
        deleteFailNotify: 'Failed to delete namespace. Error: {error}'
      }
    },
    accessControl: {
      iam: {
        identity: {
          list: {
            header: "Identities",
            uuidColumn: "UUID",
            nameColumn: "Name",
            managedColumn: "Managed by",
            actionsColumn: "Actions",
            actionsMenu: {
              delete: "Delete"
            },

            managedByServiceColumn: "Service",
            managedByIdentityColumn: "Identity",
            managedByNooneColumn: "Not managed",

            noData: "There are no identities in this namespace",
            failedToLoad: "Failed to load identities list: {error}",
            createButton: "Create",

            loadOperationNotify: 'Trying to load identities list ...',
            loadFailNotify: 'Failed to load identities list. Error: {error}'
          },
          create: {
            header: "Create identity",
    
            namespaceInput: "Namespace",
            nameInput: "Name",
            initiallyActiveInput: "Initially active",
            initiallyActiveInputCaption: "The identity will be active right from the creation and be able to access resources.",
    
            createButton: "Create",
            createHint: "System will create identity and mark it as managed by you (active user).",
    
            createOperationNotify: 'Trying to create new identity ...',
            createSuccessNotify: 'Successfully created new identity',
            createFailNotify: 'Failed to create identity. Error: {error}'
          },
          delete: {
            header: 'Delete identity?',
            bodyText: 'You are about to delete "{uuid}" identity from the "{namespace}" namespace. This action can not be undone - all the data related to the identity will be lost. This process is asynchronous, so you will have to wait for all the services will recognize the change (most probably up to several seconds).',
          
            deleteButton: 'Delete',
    
            deleteOperationNotify: 'Trying to delete identity ...',
            deleteSuccessNotify: 'Successfully deleted identity',
            deleteFailNotify: 'Failed to delete identity. Error: {error}'
          }
        },
        policy: {
          list: {
            header: "Policies",
            uuidColumn: "UUID",
            nameColumn: "Name",
            descriptionColumn: "Description",
            managedColumn: "Managed by",
            actionsColumn: "Actions",
            actionsMenu: {
              delete: "Delete"
            },

            managedByServiceColumn: "Service",
            managedByIdentityColumn: "Identity",
            managedBuiltInColumn: "Built In",
            managedByNooneColumn: "Not managed",

            noData: "There are no policies in this namespace",
            failedToLoad: "Failed to load policies list: {error}",
            createButton: "Create",

            loadOperationNotify: 'Trying to load policies list ...',
            loadFailNotify: 'Failed to load policies list. Error: {error}'
          },
          delete: {
            header: 'Delete policy?',
            bodyText: 'You are about to delete "{uuid}" policy from the "{namespace}" namespace. This action can not be undone - all the data related to the policy will be lost. This process is asynchronous, so you will have to wait for all the services will recognize the change (most probably up to several seconds).',
          
            deleteButton: 'Delete',
    
            deleteOperationNotify: 'Trying to delete policy ...',
            deleteSuccessNotify: 'Successfully deleted policy',
            deleteFailNotify: 'Failed to delete policy. Error: {error}'
          },
          create: {
            header: "Create policy",
    
            namespaceInput: "Namespace",
            nameInput: "Name",
            descriptionInput: "Description",
            namespaceIndependentInput: "Namespace independent",
            namespaceIndependentInputCaption: "Marking the policy as namespace independent will make this policy to work in all the namespaces. Selecting this option will required administrative access from you (current user).",
    
            resourcesList: {
              header: "Resources",
              caption: "Resources that can be accessed with this policy."
            },

            actionsList: {
              header: "Actions",
              caption: "Actions that can be performed using this policy on accessible resources."
            },

            createButton: "Create",
            createHint: "System will create policy and mark it as managed by you (active user).",
    
            createOperationNotify: 'Trying to create new policy ...',
            createSuccessNotify: 'Successfully created new policy',
            createFailNotify: 'Failed to create policy. Error: {error}'
          },
          view: {
            header: "Policy view",
            notSelected: "Policy to view is not selected.",
            loading: "Loading policy information ...",
            error: "Failed to load policy. Error: {error}",
            namespace: "Namespace",
            uuid: "UUID",
            name: "Name",
            description: "Description",
            namespaceIndependent: "Namespace independent",
            managedBy: "Managed by",
            created: "Created",
            updated: "Updated",
            version: "Version",
            resources: "Resources",
            resourcesCaption: "Resources that can be accessed with this policy.",
            actions: "Actions",
            actionsCaption: "Actions that can be performed using this policy on accessible resources.",
            loadOperationNotify: 'Trying to load policy information ...',
            loadFailNotify: 'Failed to load policy information. Error: {error}'
          }
        },
        role: {
          list: {
            header: "Roles",
            uuidColumn: "UUID",
            nameColumn:  "Name",
            descriptionColumn: "Description",
            managedColumn: "Managed by",
            actionsColumn: "Actions",
            actionsMenu: {
              delete: "Delete"
            },

            noData: "There are no roles in the namespace",
            failedToLoad: "Failed to load roles list: {error}",
            createButton: "Create",

            loadOperationNotify: "Trying to load roles list ...",
            loadFailNotify: "Failed to load roles list. Error: {error}"
          },
          view: {
            header: "Role view",
            notSelected: "Role to view is not selected.",
            loading: "Loading role information ...",
            error: "Failed to load role. Error: {error}",
            namespace: "Namespace",
            uuid: "UUID",
            name: "Name",
            description: "Description",
            created: "Created",
            updated: "Updated",
            version: "Version",
            loadOperationNotify: 'Trying to load role information ...',
            loadFailNotify: 'Failed to load role information. Error: {error}'
          }
        }
      }
    }
  }
};
