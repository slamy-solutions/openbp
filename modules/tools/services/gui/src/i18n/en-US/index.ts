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
        },
        iot: {
          name: "IoT",
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
        actor: {
          user: {
            list: {
              header: "Users",
              uuidColumn: "UUID",
              loginColumn: "Login",
              fullNameColumn: "Full name",
              identityColumn: "Identity",
              actionsColumn: "Actions",
              actionsMenu: {
                delete: "Delete"
              },
  
              noData: "There are no users in this namespace",
              failedToLoad: "Failed to load users list: {error}",
              createButton: "Create",
  
              loadOperationNotify: 'Trying to load users list ...',
              loadFailNotify: 'Failed to load users list. Error: {error}'
            }
          }
        },
        auth: {
          password: {
            loadingError: "Failed to load password information. Error: {error}",
            header: {
              enabled: "Password enabled",
              disabled: "Password disabled"
            },
            caption: {
              enabled: "Identity has password setted up and can authenticate with it",
              disabled: "Identity password not setted up. Authentication with password is not possible"
            },
            newPasswordInput: "New password",
            disableButton: "Disable",
            setOrUpdateButton: "Set/Update password",

            loadingOperationNotify: "Loading identity password information...",
            loadingFailNotify: "Failed to load identity password information. Error: {error}",
            
            updateOperationNotify: "Loading identity password information...",
            updateFailNotify: "Failed to load identity password information. Error: {error}",
            updateSuccessNotify: "Successfully updated password for identity",
            
            disableOperationNotify: "Loading identity password information...",
            disableFailNotify: "Failed to load identity password information. Error: {error}",
            disableSuccessNotify: "Successfully disabled password for identity"
          },
          certificate: {
            register: {
              header: "Register public key and generate certificate",
              namespaceInput: "Namespace",
              identityInput: "Identity",
              descriptionInput: "Description",
              publicKeyInfo: "Select public key file in PEM format. Currently only RSA keys are supported. You can also generate key-pair and it will be created on the client and will never leave your machine.",
              fileInput: "Public Key",
              fileRequired: "Public key file is required",
              generateKeyPairButton: "Generate key pair",
              registerAndGenerateButton: "Register key and generate certificate",
              registerHint: "After registration, the public key will be saved in the system and identity will be able to authenticate using generated X509 certificate",

              registerOperationNotify: "Registering public key and generating certificate ...",
              registerSuccessNotify: "Successfully registered public key",
              registerFailNotify: "Failed to register public key. Error: {error}"
            }
          }
        },
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
          },
          view: {
            header: "Identity view",
            notSelected: "Identity to view is not selected.",
            loading: "Loading identity information ...",
            error: "Failed to load identity. Error: {error}",
            namespace: "Namespace",
            uuid: "UUID",
            name: "Name",
            disabled: "Disabled",
            managedBy: "Managed by",
            created: "Created",
            updated: "Updated",
            version: "Version",
            policies: "Policies",
            policiesCaption: "List of the policies assigned to the identity",
            roles: "Roles",
            rolesCaption: "List of the roles assigned to the identity",
            loadOperationNotify: 'Trying to load identity information ...',
            loadFailNotify: 'Failed to load identity information. Error: {error}',
            updateOperationNotify: "Updating identity information",
            updateSuccessfullNotify: "Successfully updated identity information",
            updateFailNotify: "Failed to update identity. Error: {error}",
            activeOperationNotify: "Updating identity state",
            activeSuccessfullNotify: "Successfully updated identity state",
            activeFailNotify: "Failed to update identity state. Error: {error}",

            addPolicyOperationNotify: "Adding policy to the identity ...",
            addPolicySuccessNotify: "Successfully added policy to the identity",
            addPolicyFailNotify: "Failed to add policy to the identity. Error: {error}",
            removePolicyOperationNotify: "Removing policy from the identity ...",
            removePolicySuccessNotify: "Successfully removed policy from the identity",
            removePolicyFailNotify: "Failed to remove policy from the identity. Error: {error}",
            addRoleOperationNotify: "Adding role to the identity ...",
            addRoleSuccessNotify: "Successfully added role to the identity",
            addRoleFailNotify: "Failed to add role to the identity. Error: {error}",
            removeRoleOperationNotify: "Removing role from the identity ...",
            removeRoleSuccessNotify: "Successfully removed role from the identity",
            removeRoleFailNotify: "Failed to remove role from the identity. Error: {error}",

            tabs: {
              privileges: "Privileges",
              password: "Password",
              oauth: "OAuth",
              tokens: "Tokens",
              "2fa": "2FA",
              certificates: "Certificates"
            },
            certificatesList: {
              uuidColumn: "UUID",
              descriptionColumn: "Description",
              disabledColumn: "Disabled",
              actionsColumn: "Actions",
              noData: "Identity hasn't registered keys",
              failedToLoad: "Failed to load certificates: Error: {error}",

              registerAndGenerateButton: "Register public key and generate certificate",

              loadOperationNotify: "Loding identity certificates ...",
              loadFailNotify: "Failed to load identity certificates. Error: {error}",
              deleteOperationNotify: "Deleting certificate ...",
              deleteSuccessNotify: "Successfully deleted certificate",
              deleteFailNotify: "Failed to delete certificate. Error: {error}",
              disableOperationNotify: "Disabling certificate ...",
              disableSuccessNotify: "Successfully disabled certificate",
              disableFailNotify: "Failed to disable certificate. Error: {error}",
              
              actionsMenu: {
                delete: "Delete",
                disable: "Disable"
              }
            }
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
          },
          select: {
            header: "Select policy",
            cancelButton: "Cancel",
            selectButton: "Select",
            noData: "There are no policies yet.",
            failedToLoad: "Failed to load policies list: {error}",

            loadOperationNotify: "Loading policies list ...",
            loadFailedNotify: "Failed to load policies list: {error}",
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
          delete: {
            header: 'Delete role?',
            bodyText: 'You are about to delete "{uuid}" role from the "{namespace}" namespace. This action can not be undone - all the data related to the role will be lost. This process is asynchronous, so you will have to wait for all the services will recognize the change (most probably up to several seconds).',
          
            deleteButton: 'Delete',
    
            deleteOperationNotify: 'Trying to delete role ...',
            deleteSuccessNotify: 'Successfully deleted role',
            deleteFailNotify: 'Failed to delete role. Error: {error}'
          },
          create: {
            header: "Create role",
    
            namespaceInput: "Namespace",
            nameInput: "Name",
            descriptionInput: "Description",

            createButton: "Create",
            createHint: "System will create role and mark it as managed by you (active user).",
    
            createOperationNotify: 'Trying to create new role ...',
            createSuccessNotify: 'Successfully created new role',
            createFailNotify: 'Failed to create role. Error: {error}'
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
            managedBy: "Managed by",
            created: "Created",
            updated: "Updated",
            version: "Version",
            policies: "Policies",
            updateButton: "Update",
            policiesCaption: "List of the policies assigned to the role",
            loadOperationNotify: 'Trying to load role information ...',
            loadFailNotify: 'Failed to load role information. Error: {error}',
            updateOperationNotify: "Updating role information",
            updateSuccessfullNotify: "Successfully updated role information",
            updateFailNotify: "Failed to update role. Error: {error}",
            addPolicyOperationNotify: "Adding policy to the role ...",
            addPolicySuccessNotify: "Successfully added policy to the role",
            addPolicyFailNotify: "Failed to add policy to the role. Error: {error}",
            removePolicyOperationNotify: "Removing policy from the role ...",
            removePolicySuccessNotify: "Successfully removed policy from the role",
            removePolicyFailNotify: "Failed to remove policy from the role. Error: {error}"
          },
          select: {
            header: "Select role",
            cancelButton: "Cancel",
            selectButton: "Select",
            noData: "There are no roles yet.",
            failedToLoad: "Failed to load roles list: {error}",
            uuidColumn: "UUID",
            nameColumn: "Name",
            descriptionColumn: "Description",
            managedColumn: "Managed by",

            loadOperationNotify: "Loading roles list ...",
            loadFailedNotify: "Failed to load roles list: {error}",
          }
        }
      }
    },
    iot: {
      device: {
        selectInput: {
          label: "Select device",
        },
        select: {
          header: "Select device",
          cancelButton: "Cancel",
          selectButton: "Select",
          noData: "There are no devices yet. Maybe try to change fleet or namespace?",
          failedToLoad: "Failed to load devices list: {error}",
          uuidColumn: "UUID",
          nameColumn: "Name",
          descriptionColumn: "Description",
        },
        create: {
          header: "Create device",
    
          namespaceInput: "Namespace",
          fleetInput: "Fleet",
          nameInput: "Name",
          descriptionInput: "Description",
  
          createButton: "Create",
          createHint: "System will create iot device. Also, if fleet is selected, device will be automatically added to the fleet.",
  
          createOperationNotify: 'Trying to create new device ...',
          addOperationNotify: 'Adding device to the fleet...',
          createSuccessNotify: 'Successfully created new device',
          createFailNotify: 'Failed to create device. Error: {error}'
        },
        delete: {
          header: 'Delete device?',
          bodyText: 'You are about to delete "{uuid}" device from the "{namespace}" namespace. This action can not be undone - all the data related to the device will be lost (including events, logs and basic statistics). This process is asynchronous, so you will have to wait for all the services will recognize the change (most probably up to several seconds).',
        
          deleteButton: 'Delete',
  
          deleteOperationNotify: 'Trying to delete device ...',
          deleteSuccessNotify: 'Successfully deleted device',
          deleteFailNotify: 'Failed to delete device. Error: {error}'
        },
      },
      fleet: {
        selectInput: {
          label: "Select fleet",
        },
        select: {
          header: "Select fleet",
          cancelButton: "Cancel",
          selectButton: "Select",
          noData: "There are no fleets in this namespace yet.",
          failedToLoad: "Failed to load fleets list: {error}",
          uuidColumn: "UUID",
          nameColumn: "Name",
          descriptionColumn: "Description",
        },
        create: {
          header: "Create fleet",
    
          namespaceInput: "Namespace",
          nameInput: "Name",
          descriptionInput: "Description",
  
          createButton: "Create",
          createHint: "System will create iot devices fleet. After creation, you will be able to add devices to the fleet",
  
          createOperationNotify: 'Trying to create new fleet ...',
          createSuccessNotify: 'Successfully created new fleet',
          createFailNotify: 'Failed to create fleet. Error: {error}'
        },
        delete: {
          header: 'Delete fleet?',
          bodyText: 'You are about to delete "{uuid}" fleet from the "{namespace}" namespace. This action can not be undone - all the data related to the fleet will be lost. Devices will be unbinded from the fleet. This process is asynchronous, so you will have to wait for all the services will recognize the change (most probably up to several seconds).',
        
          deleteButton: 'Delete',
  
          deleteOperationNotify: 'Trying to delete fleet ...',
          deleteSuccessNotify: 'Successfully deleted fleet',
          deleteFailNotify: 'Failed to delete fleet. Error: {error}'
        },
        list: {
          header: "Fleets",
          uuidColumn: "UUID",
          nameColumn: "Name",
          descriptionColumn: "Description",
          actionsColumn: "Actions",
          actionsMenu: {
            delete: "Delete"
          },

          noData: "There are no device fleets in this namespace",
          failedToLoad: "Failed to load fleets list: {error}",
          createButton: "Create",

          loadOperationNotify: 'Trying to load fleets list ...',
          loadFailNotify: 'Failed to load fleets list. Error: {error}'
        },
        view: {
          header: "Fleet view",
          notSelected: "Fleet to view is not selected.",
          loading: "Loading fleet information ...",
          error: "Failed to load fleet. Error: {error}",
          namespace: "Namespace",
          uuid: "UUID",
          name: "Name",
          description: "Description",
          created: "Created",
          updated: "Updated",
          version: "Version",
          devices: "Devices",
          updateButton: "Update",
          devicesCaption: "List of the devices assigned to the fleet",
          loadOperationNotify: 'Trying to load fleet information ...',
          loadFailNotify: 'Failed to load fleet information. Error: {error}',
          updateOperationNotify: "Updating fleet information",
          updateSuccessfullNotify: "Successfully updated fleet information",
          updateFailNotify: "Failed to update fleet. Error: {error}",
        },
        deviceList: {
          header: "Devices",
          uuidColumn: "UUID",
          nameColumn: "Name",
          descriptionColumn: "Description",
          actionsColumn: "Actions",
          identityColumn: "Identity",
          actionsMenu: {
            delete: "Delete",
            changeFleet: "Change fleet",
          },

          failedToLoad: "Failed to load devices list: {error}",
          noData: "There are no devices in this fleet",
          createButton: "Create",
        },
      },
      integration: {
        balena: {
          device: {
            list: {
              header: "Balena devices",
              uuidColumn: "UUID",
              serverColumn: "Server",
              balenaUUIDColumn: "Balena UUID",
              nameColumn: "Name",

              bindedDeviceNamespaceColumn: "Bind namespace",
              bindedDeviceColumn: "Bind device",
              statusColumn: "Status",
              isOnlineColumn: "Online",
              cpuUsageColumn: "CPU usage",
              cpuTempColumn: "CPU °C",
              ramUsageColumn: "RAM",
              lastConnectivityColumn: "Last connectivity",
              actionsColumn: "Actions",

              noData: "There are no balena devices yet.",
              failedToLoad: "Failed to load balena devices list: {error}",
              bindButton: "Disable",
              unbindButton: "Delete",

              actionsMenu: {
                bind: "Bind"
              },

              loadOperationNotify: 'Trying to load devices list ...',
              loadFailNotify: 'Failed to load devices list. Error: {error}'
            },
          },
          server: {
            list: {
              header: "Balena servers",
              uuidColumn: "UUID",
              nameColumn: "Name",
              descriptionColumn: "Description",
              enabledColumn: "Enabled",
              actionsColumn: "Actions",
              syncStatusColumn: "Sync status",
              syncDevicesColumn: "Synced devices",
              syncTimeColumn: "Sync duration",
              syncTimestampColumn: "Sync timestamp",

              actionsMenu: {
                delete: "Delete"
              },

              noData: "There are no balena servers yet.",
              failedToLoad: "Failed to load balena servers list: {error}",
              createButton: "Add balena server",
              enableButton: "Enable",
              disableButton: "Disable",
              deleteButton: "Delete",

              loadFailNotify: 'Failed to load servers list. Error: {error}',
              setEnabledOperationNotify: 'Trying to set server enabled state ...',
              setEnabledSuccessNotify: 'Successfully set server enabled state',
              setEnabledFailNotify: 'Failed to set server enabled state. Error: {error}',
            },
            add: {
              header: "Add Balena server",
        
              namespaceInput: "Namespace",
              nameInput: "Name",
              descriptionInput: "Description",
              urlInput: "API URL",
              apiTokenInput: "API token",
      
              validateConnectionButton: "Validate connection data",
              connectionDataNotValidated: "You have to validate connection data before adding the server. This will make request to the provided server with provided credentials and will check if server is accessible and credentials are valid.",
              connectionDataIsValid: "Connection data is valid. You can add the server now.",
              createButton: "Add server",
              createHint: "Provided information will be used to communicate with balena server in the background and get information about devices. The communication will not be enabled immediatelly, you will have to enable it manually after server is added. The information on how to connect to the server will be encrypted using system_vault hardware security module.",
      
              validateOperationNotify: 'Trying to validate connection data ...',
              validateFailNotify: 'Failed to validate connection data. Error: {error}',

              addOperationNotify: 'Trying to add server ...',
              addSuccessNotify: 'Successfully added server',
              addFailNotify: 'Failed to add server. Error: {error}'
            },
          }
        }
      }
    }
  }
};
