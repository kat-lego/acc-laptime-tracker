@description('Location for all resources')
param location string = resourceGroup().location

@secure()
@description('Cosmos DB connection string')
param cosmosDbConnectionString string

@description('Cosmos DB database name')
param cosmosDbDatabase string = 'sessions'

@description('Cosmos DB container name')
param cosmosDbContainer string = 'sessions'

@secure()
@description('tui ssh host key')
param tuiSshHostKeyPEM string

resource vnet 'Microsoft.Network/virtualNetworks@2024-05-01' = {
  name: 'vn-acc-laptime-tracker'
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    subnets: [
      {
        name: 'default'
        properties: {
          addressPrefix: '10.0.0.0/23'
        }
      }
    ]
  }
}

resource logAnalytics 'Microsoft.OperationalInsights/workspaces@2025-02-01' = {
  name: 'ws-acc-laptime-tracker'
  location: location
  properties: {
    sku: {
      name: 'PerGB2018'
    }
    retentionInDays: 30
  }
}

resource containerEnv 'Microsoft.App/managedEnvironments@2025-01-01' = {
  name: 'menv-acc-laptime-tracker'
  location: location
  properties: {
    vnetConfiguration: {
      infrastructureSubnetId: vnet.properties.subnets[0].id
    }
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: logAnalytics.properties.customerId
        sharedKey: logAnalytics.listKeys().primarySharedKey
      }
    }
  }
  dependsOn: [
    #disable-next-line no-unnecessary-dependson
    vnet
    #disable-next-line no-unnecessary-dependson
    logAnalytics
  ]
}

resource apiContainerApp 'Microsoft.App/containerApps@2025-01-01' = {
  name: 'capp-acc-laptime-tracker-api'
  location: location
  properties: {
    managedEnvironmentId: containerEnv.id
    configuration: {
      ingress: {
        external: true
        targetPort: 80
        transport: 'auto'
      }
      registries: []
    }
    template: {
      containers: [
        {
          name: 'api-container'
          image: '99percentmoses/acc-laptime-tracker-api:latest'
          resources: {
            #disable-next-line BCP036
            cpu: '0.5'
            memory: '1.0Gi'
          }
          env: [
            {
              name: 'PORT'
              value: '80'
            }
            {
              name: 'ACC_COSMOS_CONNECTION_STRING'
              value: cosmosDbConnectionString
            }
            {
              name: 'ACC_COSMOS_DATABASE'
              value: cosmosDbDatabase
            }
            {
              name: 'ACC_COSMOS_CONTAINER'
              value: cosmosDbContainer
            }
          ]
        }
      ]
    }
  }
  dependsOn: [
    #disable-next-line no-unnecessary-dependson
    containerEnv
  ]
}

resource tuiContainerApp 'Microsoft.App/containerApps@2025-01-01' = {
  name: 'capp-acc-laptime-tracker-tui'
  location: location
  properties: {
    managedEnvironmentId: containerEnv.id
    configuration: {
      ingress: {
        external: true
        targetPort: 8080
        exposedPort: 8080
        transport: 'tcp'
      }
      registries: []
    }
    template: {
      containers: [
        {
          name: 'api-container'
          image: '99percentmoses/acc-laptime-tracker-tui:latest'
          resources: {
            #disable-next-line BCP036
            cpu: '0.5'
            memory: '1.0Gi'
          }
          env: [
            {
              name: 'SSH_HOST_KEY_PEM'
              value: tuiSshHostKeyPEM
            }
          ]
        }
      ]
    }
  }
  dependsOn: [
    #disable-next-line no-unnecessary-dependson
    containerEnv
  ]
}
